// watcher is a package that monitors workflows and controls tasks
package watcher

import (
	"context"
	"time"

	"github.com/d-tsuji/flower/repository"
	"github.com/pkg/errors"
)

// watcherTask contains the channel of the task being watched.
type watcherTask struct {
	db         *repository.DB
	ExecTaskCh chan repository.ExecutableTask
}

// NewWatcherTask creates a new watcherTask.
func NewWatcherTask(db *repository.DB, execTaskCh chan repository.ExecutableTask) *watcherTask {
	return &watcherTask{
		db:         db,
		ExecTaskCh: execTaskCh,
	}
}

// WatchTask searches for tasks that are waiting to be executed and can be executed.
// If the target task exists, update the status of the task using an optimistic lock
// and assign the job to a worker existing in the worker pool.
func (w *watcherTask) WatchTask(ctx context.Context, concurrency int) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	waitingTasks, err := w.db.GetExecutableTask(ctx, concurrency)
	if err != nil {
		return errors.WithStack(err)
	}

	var runTasks []repository.ExecutableTask
	for _, wt := range waitingTasks {
		ok, err := w.db.UpdateExecutableTasksRunning(ctx, wt)
		if err != nil {
			return errors.WithStack(err)
		}

		// The record was already updated during execution, so it is not added to the execution target.
		if !ok {
			continue
		}
		runTasks = append(runTasks, wt)
	}

	// Put to worker as execution task.
	for _, rt := range runTasks {
		w.ExecTaskCh <- rt
	}

	return nil
}
