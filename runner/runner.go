package runner

import (
	"context"
	"log"
	"reflect"

	"github.com/d-tsuji/flower-v2/db"
	internal "github.com/d-tsuji/flower-v2/internal"
	"github.com/pkg/errors"
)

type runner struct {
	task db.ExecutableTask
	db   *db.DB
}

func NewRunner(task db.ExecutableTask, db *db.DB) *runner {
	return &runner{
		task: task,
		db:   db,
	}
}

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

func (r *runner) runTask(ctx context.Context) ([]reflect.Value, error) {
	log.Printf("[runner] runTask executing...\n")

	programName, err := r.db.GetTaskProgramName(ctx, r.task)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	executor := internal.NewExecutor()
	method := reflect.ValueOf(executor).MethodByName(programName)
	in := make([]reflect.Value, 0)
	values := method.Call(in)

	return values, nil
}
