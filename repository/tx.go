// I implemented transaction management with reference to
// https://github.com/google/trillian/blob/master/storage/postgres/admin_storage.go.

package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"

	"github.com/pkg/errors"
)

type adminTX struct {
	tx *sql.Tx

	// mu guards *direct* reads/writes on closed, which happen only on
	// Commit/Rollback/IsClosed/Close methods.
	// We don't check closed on *all* methods (apart from the ones above),
	// as we trust tx to keep tabs on its state (and consequently fail to do
	// queries after closed).
	mu     sync.RWMutex
	closed bool
}

func (t *adminTX) Commit() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	return t.tx.Commit()
}

func (t *adminTX) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	return t.tx.Rollback()
}

func (t *adminTX) IsClosed() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.closed
}

func (t *adminTX) Close() error {
	// Acquire and release read lock manually, without defer, as if the txn
	// is not closed Rollback() will attempt to acquire the rw lock.
	t.mu.RLock()
	closed := t.closed
	t.mu.RUnlock()
	if !closed {
		err := t.Rollback()
		if err != nil {
			log.Printf("Rollback error on Close(): %v", err)
		}
		return err
	}
	return nil
}

// AdminTXFunc is the signature for functions passed to ReadWriteTransaction.
type AdminTXFunc func(context.Context, *adminTX) error

func (t *adminTX) updateExecutableTasks(ctx context.Context, e ExecutableTask, beforeTaskStatus, afterTaskStatus int, ok *bool) error {
	// Acquire a lock for updating a record.
	stmtUpdateExecutableTaskLock, err := t.tx.PrepareContext(ctx, selectUpdateExecutableTaskLock)
	if err != nil {
		*ok = false
		return errors.New(fmt.Sprintf("query prepare error 'SELECT kr_task_stat ... FOR UPDATE': %v", err))
	}
	defer stmtUpdateExecutableTaskLock.Close()
	result, err := stmtUpdateExecutableTaskLock.ExecContext(ctx, e.TaskFlowId, e.TaskExecSeq, beforeTaskStatus)
	if err != nil {
		*ok = false
		return errors.New(fmt.Sprintf("query error: %v", err))
	}
	num, err := result.RowsAffected()
	if err != nil {
		*ok = false
		return errors.WithStack(err)
	}
	if num == 0 {
		*ok = false
		return nil
	}

	// Updates the record that acquired the lock.
	stmtUpdateExecutableTask, err := t.tx.PrepareContext(ctx, updateExecutableTask)
	if err != nil {
		*ok = false
		return errors.New(fmt.Sprintf("query prepare error 'UPDATE kr_task_stat': %v", err))
	}
	defer stmtUpdateExecutableTask.Close()
	_, err = stmtUpdateExecutableTask.ExecContext(ctx, afterTaskStatus, e.TaskFlowId, e.TaskExecSeq, beforeTaskStatus)
	if err != nil {
		*ok = false
		return errors.New(fmt.Sprintf("query error: %v", err))
	}
	*ok = true
	return nil
}

func (t *adminTX) getRegisterTask(ctx context.Context, taskId string) ([]task, error) {
	var tasks []task
	rows, err := t.tx.QueryContext(ctx, selectRegisterTask, taskId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("query error: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		task, err := readTask(rows)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		tasks = append(tasks, *task)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %v", err))
	}

	return tasks, nil
}

func (t *adminTX) insertExecutableTasks(ctx context.Context, tasks []task) error {
	stmt, err := t.tx.PrepareContext(ctx, insertExecutableTasks)
	if err != nil {
		return errors.New(fmt.Sprintf("query prepare error 'INSERT INTO kr_task_stat': %v", err))
	}
	defer stmt.Close()

	taskExecSec, dependsTaskExecSec := 1, -1
	taskFlowId, err := uuid.NewUUID()
	if err != nil {
		return errors.New(fmt.Sprintf("generate uuid error: %v", err))
	}
	for _, task := range tasks {
		params, err := json.Marshal(task.Params)
		if err != nil {
			return errors.New(fmt.Sprintf("json marshal error: %v", err))
		}
		_, err = stmt.ExecContext(ctx,
			taskFlowId,
			taskExecSec,
			dependsTaskExecSec,
			string(params),
			task.TaskId,
			task.TaskSeq,
			0,
			task.TaskPriority,
		)
		if err != nil {
			return errors.New(fmt.Sprintf("query error: %v", err))
		}
		dependsTaskExecSec = taskExecSec
		taskExecSec++
	}
	return nil
}
