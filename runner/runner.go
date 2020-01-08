package runner

import (
	"context"
	"log"
	"reflect"

	"github.com/d-tsuji/flower/component"

	"github.com/d-tsuji/flower/repository"
	"github.com/pkg/errors"
)

// Runner is a struct for executing tasks.
type runner struct {
	task repository.ExecutableTask
	db   *repository.DB
}

// NewServer creates a new Runner.
func NewRunner(task repository.ExecutableTask, db *repository.DB) *runner {
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
		if err := r.db.UpdateExecutableTasksSuspended(ctx, r.task); err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(err)
	}

	if err = r.db.UpdateExecutableTasksFinished(ctx, r.task); err != nil {
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

	executor := component.NewExecutor(r.task.Params)
	method := reflect.ValueOf(executor).MethodByName(programName)
	values := method.Call(nil)

	return values, nil
}
