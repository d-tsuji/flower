package watcher

import (
	"context"
	"time"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/pkg/errors"
)

type watcherTask struct {
	pollingIntervalSecond time.Duration
	concurrency           int
	db                    *db.DB
	ExecTaskCh            chan db.ExecutableTask
}

const (
	defaultPollingInterval = 10
	defaultConcurrency     = 1
)

func NewWatcherTask(db *db.DB, execTaskCh chan db.ExecutableTask) *watcherTask {
	return &watcherTask{
		pollingIntervalSecond: defaultPollingInterval,
		concurrency:           defaultConcurrency,
		db:                    db,
		ExecTaskCh:            execTaskCh,
	}
}

// WatchTask searches for tasks that are waiting to be executed and can be executed.
// If the target task exists, update the status of the task using an optimistic lock
// and assign the job to a worker existing in the worker pool.
func (w *watcherTask) WatchTask(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	waitingTasks, err := w.db.GetExecutableTask(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	var runTasks []db.ExecutableTask
	for _, wt := range waitingTasks {
		ok, err := w.db.UpdateExecutableTasksRunning(wt)
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
