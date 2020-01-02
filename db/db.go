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
	TaskId  string `db:"task_id"`
	TaskSeq int    `db:"task_seq"`
	Program string `db:"program"`
}

func (db *DB) getRegisterTask(taskId string) ([]task, error) {
	var tasks []task

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	rows, err := db.QueryContext(ctx, `SELECT task_id, task_seq, program FROM ms_task_definition WHERE task_id = $1`, taskId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("no task get error: %v", err))
	}

	for rows.Next() {
		var (
			taskId  string
			taskSeq int
			program string
		)
		if err := rows.Scan(&taskId, &taskSeq, &program); err != nil {
			return nil, errors.New(fmt.Sprintf("rows scan error: %v", err))
		}
		tasks = append(tasks, task{
			TaskId:  taskId,
			TaskSeq: taskSeq,
			Program: program,
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
	INSERT INTO kr_task_stat (task_flow_id, task_exec_seq, depends_task_exec_seq, parameters, task_id, task_seq, exec_status) VALUES ($1, $2, $3, $4, $5, $6, $7);`)
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
