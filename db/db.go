package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	// database/sql driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	// DefaultDBName is the default name of postgres database.
	DefaultDBName = "flower"
	// DefaultDBUserName is the default postgres user name.
	DefaultDBUserName = "flower"
	// DefaultDBHostname is the default postgres host name.
	DefaultDBHostname = "localhost"
	// DefaultDBPort is the default postgres port.
	DefaultDBPort = "5432"
	// ConnectTimeout is the default timeout of the connection to the postgres server.
	ConnectTimeout = 5

	// ExecStatusWait is the status waiting to be executed.
	ExecStatusWait = 0
	// ExecStatusRunning is the running status.
	ExecStatusRunning = 1
	// ExecStatusSuspend is the suspended status.
	ExecStatusSuspend = 2
	// ExecStatusFinish is the completed status.
	ExecStatusFinish = 3
	// ExecStatusIgnore is the status to be ignored.
	ExecStatusIgnore = 9
)

// DB represents a Database handler.
type DB struct {
	*sql.DB
}

// Opt are options for database connection.
// https://godoc.org/github.com/lib/pq
type Opt struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

// New creates the DB object.
func New(opt *Opt) (*DB, error) {
	var user, dbname, host, port, sslmode string
	if user = opt.User; user == "" {
		user = DefaultDBUserName
	}
	if dbname = opt.DBName; dbname == "" {
		dbname = DefaultDBName
	}
	if host = opt.Host; host == "" {
		host = DefaultDBHostname
	}
	if port = opt.Port; port == "" {
		port = DefaultDBPort
	}
	if sslmode = opt.SSLMode; sslmode == "" {
		sslmode = "disable"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s connect_timeout=%d",
		user, opt.Password, host, port, dbname, sslmode, ConnectTimeout,
	))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("postgres open error: %v", err))
	}
	if err = db.Ping(); err != nil {
		return nil, errors.New(fmt.Sprintf("postgres ping error: %v", err))
	}
	return &DB{db}, nil
}

// Task is one record of ms_task.
type task struct {
	TaskId       string `db:"task_id"`
	TaskSeq      int    `db:"task_seq"`
	Program      string `db:"program"`
	TaskPriority int    `db:"task_priority"`
}

// GetRegisterTask gets a series of tasks registered in the master from taskId.
func (db *DB) getRegisterTask(ctx context.Context, taskId string) ([]task, error) {
	var tasks []task

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	rows, err := db.QueryContext(ctx, `SELECT task_id, task_seq, program, task_priority FROM ms_task_definition WHERE task_id = $1`, taskId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("query error: %v", err))
	}

	for rows.Next() {
		var (
			taskId       string
			taskSeq      int
			program      string
			taskPriority int
		)
		if err := rows.Scan(&taskId, &taskSeq, &program, &taskPriority); err != nil {
			return nil, errors.New(fmt.Sprintf("rows scan error: %v", err))
		}
		tasks = append(tasks, task{
			TaskId:       taskId,
			TaskSeq:      taskSeq,
			Program:      program,
			TaskPriority: taskPriority,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %v", err))
	}

	return tasks, nil
}

// ExecutableTask is the struct of executable tasks.
type ExecutableTask struct {
	TaskFlowId  string `db:"task_flow_id"`
	TaskExecSeq int    `db:"task_exec_seq"`
	TaskId      string `db:"task_id"`
	TaskSeq     int    `db:"task_seq"`
}

// GetExecutableTask is the main method.
// From the series of tasks waiting to be registered for each task flow Id,
// resolve the dependency of the execution order.
// Get the executable task and register the task to be executed in the channel to run task.
func (db *DB) GetExecutableTask(ctx context.Context, concurrency int) ([]ExecutableTask, error) {
	var tasks []ExecutableTask

	rows, err := db.QueryContext(ctx, `
		SELECT
			base.task_flow_id
		,	base.task_exec_seq
		,	base.task_id
		,	base.task_seq
		FROM
			(
				SELECT
					base.task_flow_id
				,	base.task_exec_seq
				,	base.task_id
				,	base.task_seq
				,	base.exec_status
				,	ROW_NUMBER() OVER (ORDER BY base.exec_status DESC, base.task_exec_seq, base.task_priority) rowno
				FROM
					kr_task_stat base
				LEFT JOIN
					kr_task_stat dep
					ON	1=1
						AND	base.depends_task_exec_seq = dep.task_exec_seq
						AND	base.task_id = dep.task_id
						AND	base.task_flow_id = dep.task_flow_id
				-- Task is running or dependent task has completed and is waiting to run
				WHERE	base.exec_status = 1 OR (COALESCE(dep.exec_status, 3) = 3 AND base.exec_status = 0)
			)base
		-- Control within the number of concurrent executions
		WHERE rowno <= $1 AND base.exec_status = 0;`, concurrency)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("query error: %v", err))
	}

	for rows.Next() {
		var (
			taskFlowId  string
			taskExecSeq int
			taskId      string
			taskSeq     int
		)
		if err := rows.Scan(&taskFlowId, &taskExecSeq, &taskId, &taskSeq); err != nil {
			return nil, errors.New(fmt.Sprintf("rows scan error: %v", err))
		}
		tasks = append(tasks, ExecutableTask{
			TaskFlowId:  taskFlowId,
			TaskExecSeq: taskExecSeq,
			TaskId:      taskId,
			TaskSeq:     taskSeq,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %v", err))
	}

	return tasks, nil
}

// InsertExecutableTasks registers the task waiting to be executed from the called taskId.
func (db *DB) InsertExecutableTasks(ctx context.Context, taskId string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	tasks, err := db.getRegisterTask(ctx, taskId)
	if err != nil {
		return errors.WithStack(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("begin transaction error: %v", err))
	}
	stmtInsertExecutableTask, err := tx.PrepareContext(ctx, `
	INSERT INTO kr_task_stat(
		task_flow_id
	,	task_exec_seq
	,	depends_task_exec_seq
	,	parameters
	,	task_id
	,	task_seq
	,	exec_status
	,	task_priority)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`)
	if err != nil {
		return errors.New(fmt.Sprintf("query prepare error 'INSERT INTO kr_task_stat': %v", err))
	}

	taskExecSec, dependsTaskExecSec := 1, 0
	taskFlowId, err := uuid.NewUUID()
	if err != nil {
		return errors.New(fmt.Sprintf("generate uuid error: %v", err))
	}
	for _, task := range tasks {
		_, err := stmtInsertExecutableTask.ExecContext(ctx,
			taskFlowId,
			taskExecSec,
			dependsTaskExecSec,
			"{}",
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

	if err := tx.Commit(); err != nil {
		return errors.New(fmt.Sprintf("transaction commit error: %v", err))
	}
	return nil
}

// UpdateExecutableTasksRunning updates task status to running.
func (db *DB) UpdateExecutableTasksRunning(ctx context.Context, e ExecutableTask) (bool, error) {
	return db.updateExecutableTasks(ctx, e, ExecStatusWait, ExecStatusRunning)
}

// UpdateExecutableTasksRunning updates the status of tasks to finished.
func (db *DB) UpdateExecutableTasksFinished(ctx context.Context, e ExecutableTask) (bool, error) {
	return db.updateExecutableTasks(ctx, e, ExecStatusRunning, ExecStatusFinish)
}

// UpdateExecutableTasksRunning updates the status of a task to suspended.
func (db *DB) UpdateExecutableTasksSuspended(ctx context.Context, e ExecutableTask) (bool, error) {
	return db.updateExecutableTasks(ctx, e, ExecStatusRunning, ExecStatusSuspend)
}

// UpdateExecutableTasks updates task status with an optimistic lock.
// It is not an error if the record has already been updated,
// but returns false as the first argument of the return value.
func (db *DB) updateExecutableTasks(ctx context.Context, e ExecutableTask, beforeTaskStatus, afterTaskStatus int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, errors.New(fmt.Sprintf("begin transaction error: %v", err))
	}

	// Acquire a lock for updating a record.
	stmtUpdateExecutableTaskLock, err := tx.PrepareContext(ctx, `
	SELECT * FROM kr_task_stat WHERE task_flow_id = $1 and task_exec_seq = $2 and exec_status = $3 FOR UPDATE;`)
	if err != nil {
		return false, errors.New(fmt.Sprintf("query prepare error 'SELECT kr_task_stat ... FOR UPDATE': %v", err))
	}
	result, err := stmtUpdateExecutableTaskLock.ExecContext(ctx, e.TaskFlowId, e.TaskExecSeq, beforeTaskStatus)
	if err != nil {
		return false, errors.New(fmt.Sprintf("query error: %v", err))
	}
	num, err := result.RowsAffected()
	if err != nil {
		return false, errors.WithStack(err)
	}
	if num == 0 {
		return false, nil
	}

	// Updates the record that acquired the lock.
	stmtUpdateExecutableTask, err := tx.PrepareContext(ctx, `
	UPDATE kr_task_stat SET exec_status = $1 WHERE task_flow_id = $2 and task_exec_seq = $3 and exec_status = $4;`)
	if err != nil {
		return false, errors.New(fmt.Sprintf("query prepare error 'UPDATE kr_task_stat': %v", err))
	}
	result, err = stmtUpdateExecutableTask.ExecContext(ctx, afterTaskStatus, e.TaskFlowId, e.TaskExecSeq, beforeTaskStatus)
	if err != nil {
		return false, errors.New(fmt.Sprintf("query error: %v", err))
	}
	num, err = result.RowsAffected()
	if err != nil {
		return false, errors.WithStack(err)
	}
	if num == 0 {
		return false, nil
	}

	if err := tx.Commit(); err != nil {
		return false, errors.New(fmt.Sprintf("transaction commit error: %v", err))
	}
	return true, nil
}

// GetTaskProgramName gets the name of the program to be executed from taskId and taskSeq.
func (db *DB) GetTaskProgramName(ctx context.Context, task ExecutableTask) (string, error) {
	var programName string

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	rows, err := db.QueryContext(ctx, `SELECT program FROM ms_task_definition WHERE task_id = $1 AND task_seq = $2`, task.TaskId, task.TaskSeq)
	if err != nil {
		return "", errors.New(fmt.Sprintf("query error: %v", err))
	}

	for rows.Next() {
		if err := rows.Scan(&programName); err != nil {
			return "", errors.New(fmt.Sprintf("rows scan error: %v", err))
		}
	}
	if err := rows.Err(); err != nil {
		return "", errors.New(fmt.Sprintf("rows error: %v", err))
	}

	return programName, nil
}
