package watcher

import (
	"context"
	"log"

	"github.com/d-tsuji/flower/repository"
	"github.com/d-tsuji/flower/runner"
)

type Worker struct {
	ID            int
	WorkerChannel chan chan repository.ExecutableTask
	Channel       chan repository.ExecutableTask
	DBClient      *repository.DB
}

// start worker
func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			// when the worker is available place channel in queue
			w.WorkerChannel <- w.Channel
			select {
			case job := <-w.Channel:
				r := runner.NewRunner(job, w.DBClient)
				if err := r.Run(ctx); err != nil {
					log.Printf("[worker] runner.Run() is failed. err: %v\n", err)
				}
			case <-ctx.Done():
				log.Printf("[worker] worker [%d] is stopping\n", w.ID)
				return
			}
		}
	}()
}
