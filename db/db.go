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

func (db *DB) getRegisterTask(taskId string) ([]task, error) {
	var tasks []task

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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

type ExecutableTask struct {
	TaskFlowId  string `db:"task_flow_id"`
	TaskExecSeq int    `db:"task_exec_seq"`
	TaskId      string `db:"task_id"`
	TaskSeq     int    `db:"task_seq"`
}

func (db *DB) GetExecutableTask(ctx context.Context) ([]ExecutableTask, error) {
	var tasks []ExecutableTask

	rows, err := db.QueryContext(ctx, `
		SELECT
			base.task_flow_id
		,	base.task_exec_seq
		,	base.task_id
		,	base.task_seq
		FROM
			kr_task_stat base
		LEFT JOIN
			kr_task_stat dep
			ON	1=1
				AND	base.depends_task_exec_seq = dep.task_exec_seq
				AND	base.task_id = dep.task_id
				AND	base.task_flow_id = dep.task_flow_id
		WHERE	1=1
			AND	COALESCE(dep.exec_status, '3') = '3' --依存するタスクが完了しているタスク
			AND	base.exec_status IN ('0')            --実行待ちのタスク
		ORDER BY
			base.task_priority
		`)
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

func (db *DB) InsertExecutableTasks(taskId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	tasks, err := db.getRegisterTask(taskId)
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

func (db *DB) UpdateExecutableTasksRunning(e ExecutableTask) (bool, error) {
	return db.updateExecutableTasks(e, ExecStatusWait, ExecStatusRunning)
}

func (db *DB) UpdateExecutableTasksFinished(e ExecutableTask) (bool, error) {
	return db.updateExecutableTasks(e, ExecStatusRunning, ExecStatusFinish)
}

func (db *DB) UpdateExecutableTasksSuspended(e ExecutableTask) (bool, error) {
	return db.updateExecutableTasks(e, ExecStatusRunning, ExecStatusSuspend)
}

// UpdateExecutableTasks updates task status with an optimistic lock.
// It is not an error if the record has already been updated,
// but returns false as the first argument of the return value.
func (db *DB) updateExecutableTasks(e ExecutableTask, beforeTaskStatus, afterTaskStatus int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, errors.New(fmt.Sprintf("begin transaction error: %v", err))
	}
	stmtInsertExecutableTask, err := tx.PrepareContext(ctx, `
	UPDATE kr_task_stat SET exec_status = $1 WHERE task_flow_id = $2 and task_exec_seq = $3 and exec_status = $4;`)
	if err != nil {
		return false, errors.New(fmt.Sprintf("query prepare error 'UPDATE kr_task_stat': %v", err))
	}
	result, err := stmtInsertExecutableTask.ExecContext(ctx, afterTaskStatus, e.TaskFlowId, e.TaskExecSeq, beforeTaskStatus)
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
	if err := tx.Commit(); err != nil {
		return false, errors.New(fmt.Sprintf("transaction commit error: %v", err))
	}
	return true, nil
}

func (db *DB) GetTaskProgramName(task ExecutableTask) (string, error) {
	var programName string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
