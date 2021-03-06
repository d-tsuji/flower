package watcher

import (
	"context"
	"log"

	"github.com/d-tsuji/flower/repository"
	"github.com/d-tsuji/flower/runner"
)

// Worker play the role of workers in the Worker-Pools model.
// The structure that executes the task is runner.
type Worker struct {
	id            int
	WorkerChannel chan chan repository.ExecutableTask
	Channel       chan repository.ExecutableTask
	dbClient      *repository.DB
}

// Start makes one worker start.
func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			// when the worker is available place channel in queue
			w.WorkerChannel <- w.Channel
			select {
			case job := <-w.Channel:
				r := runner.NewRunner(&job, w.dbClient)
				if err := r.Run(ctx); err != nil {
					log.Printf("[worker] runner.Run() is failed. err: %v\n", err)
				}
			case <-ctx.Done():
				log.Printf("[worker] worker [%d] is stopping\n", w.id)
				return
			}
		}
	}()
}
