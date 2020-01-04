package runner

import (
	"context"
	"log"
	"reflect"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/pkg/errors"
)

//　Runner is a struct for executing tasks
type runner struct {
	task db.ExecutableTask
	db   *db.DB
}

// NewServer creates a new Runner.
func NewRunner(task db.ExecutableTask, db *db.DB) *runner {
	return &runner{
		task: task,
		db:   db,
	}
}

// Run calls runTask. Update the status of the task according to the result.
func (r *runner) Run(ctx context.Context) error {
	// TODO: handle response params
	_, err := r.runTask(ctx)
	if err != nil {
		if _, err := r.db.UpdateExecutableTasksSuspended(ctx, r.task); err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(err)
	}

	if _, err = r.db.UpdateExecutableTasksFinished(ctx, r.task); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// RunTask method calls the program registered
// in the master of the task definition using reflect.
func (r *runner) runTask(ctx context.Context) ([]reflect.Value, error) {
	log.Printf("[runner] runTask executing...\n")

	programName, err := r.db.GetTaskProgramName(ctx, r.task)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	executor := NewExecutor()
	method := reflect.ValueOf(executor).MethodByName(programName)
	in := make([]reflect.Value, 0)
	values := method.Call(in)

	return values, nil
}
