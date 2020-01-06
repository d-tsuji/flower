package repository

import (
	"context"
	"database/sql"
	"fmt"

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

	selectRegisterTask = `
		SELECT
			task_id, task_seq, program, task_priority,
			param1_key, param1_value, param2_key, param2_value, param3_key, param3_value,
			param4_key, param4_value, param5_key, param5_value
		FROM ms_task_definition WHERE task_id = $1;`
	selectTaskProgramName = `SELECT program FROM ms_task_definition WHERE task_id = $1 AND task_seq = $2;`
	selectExecutableTask  = `
		SELECT
			base.task_flow_id
		,	base.task_exec_seq
		,	base.task_id
		,	base.task_seq
		,	base.parameters
		FROM
			(
				SELECT
					base.task_flow_id
				,	base.task_exec_seq
				,	base.task_id
				,	base.task_seq
				,	base.exec_status
				,	base.parameters
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
		WHERE rowno <= $1 AND base.exec_status = 0;`
	insertExecutableTasks = `
	INSERT INTO kr_task_stat(
		task_flow_id
	,	task_exec_seq
	,	depends_task_exec_seq
	,	parameters
	,	task_id
	,	task_seq
	,	exec_status
	,	task_priority)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	selectUpdateExecutableTaskLock = `SELECT * FROM kr_task_stat WHERE task_flow_id = $1 and task_exec_seq = $2 and exec_status = $3 FOR UPDATE;`
	updateExecutableTask           = `UPDATE kr_task_stat SET exec_status = $1 WHERE task_flow_id = $2 and task_exec_seq = $3 and exec_status = $4;`
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

// Task is one record of ms_task.
type task struct {
	TaskId       string `db:"task_id"`
	TaskSeq      int    `db:"task_seq"`
	Program      string `db:"program"`
	TaskPriority int    `db:"task_priority"`
	Params       map[string]string
}

// ExecutableTask is the struct of executable tasks.
type ExecutableTask struct {
	TaskFlowId  string            `db:"task_flow_id"`
	TaskExecSeq int               `db:"task_exec_seq"`
	TaskId      string            `db:"task_id"`
	TaskSeq     int               `db:"task_seq"`
	Params      map[string]string `db:"parameters"`
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

func (db *DB) beginInternal(ctx context.Context) (*adminTX, error) {
	tx, err := db.BeginTx(ctx, nil /* opts */)
	if err != nil {
		return nil, err
	}
	return &adminTX{tx: tx}, nil
}

// ReadWriteTransaction creates a transaction, and runs f with it.
// Some storage implementations may retry aborted transactions, so
// f MUST be idempotent.
func (db *DB) ReadWriteTransaction(ctx context.Context, f AdminTXFunc) error {
	tx, err := db.beginInternal(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	if err := f(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

// GetExecutableTask is the main method.
// From the series of tasks waiting to be registered for each task flow Id,
// resolve the dependency of the execution order.
// Get the executable task and register the task to be executed in the channel to run task.
func (db *DB) GetExecutableTask(ctx context.Context, concurrency int) ([]ExecutableTask, error) {
	var tasks []ExecutableTask

	rows, err := db.QueryContext(ctx, selectExecutableTask, concurrency)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("query error: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		task, err := readExecutableTask(rows)
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

// InsertExecutableTasks registers the task waiting to be executed from the called taskId.
func (db *DB) InsertExecutableTasks(ctx context.Context, taskId string) error {
	err := db.ReadWriteTransaction(ctx, func(ctx context.Context, t *adminTX) error {
		tasks, err := t.getRegisterTask(ctx, taskId)
		if err != nil {
			return errors.WithStack(err)
		}

		err = t.insertExecutableTasks(ctx, tasks)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	return err
}

// UpdateExecutableTasksRunning updates task status to running.
func (db *DB) UpdateExecutableTasksRunning(ctx context.Context, e ExecutableTask) (bool, error) {
	var ok bool
	err := db.ReadWriteTransaction(ctx, func(ctx context.Context, t *adminTX) error {
		return t.updateExecutableTasks(ctx, e, ExecStatusWait, ExecStatusRunning, &ok)
	})
	return ok, err
}

// UpdateExecutableTasksRunning updates the status of tasks to finished.
func (db *DB) UpdateExecutableTasksFinished(ctx context.Context, e ExecutableTask) (bool, error) {
	var ok bool
	err := db.ReadWriteTransaction(ctx, func(ctx context.Context, t *adminTX) error {
		return t.updateExecutableTasks(ctx, e, ExecStatusRunning, ExecStatusFinish, &ok)
	})
	return ok, err
}

// UpdateExecutableTasksRunning updates the status of a task to suspended.
func (db *DB) UpdateExecutableTasksSuspended(ctx context.Context, e ExecutableTask) (bool, error) {
	var ok bool
	err := db.ReadWriteTransaction(ctx, func(ctx context.Context, t *adminTX) error {
		err := t.updateExecutableTasks(ctx, e, ExecStatusRunning, ExecStatusSuspend, &ok)
		return err
	})
	return ok, err
}

// GetTaskProgramName gets the name of the program to be executed from taskId and taskSeq.
func (db *DB) GetTaskProgramName(ctx context.Context, task ExecutableTask) (string, error) {
	var programName string
	rows, err := db.QueryContext(ctx, selectTaskProgramName, task.TaskId, task.TaskSeq)
	if err != nil {
		return "", errors.New(fmt.Sprintf("query error: %v", err))
	}
	defer rows.Close()

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
