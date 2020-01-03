package watcher

import (
	"context"
	"fmt"
	"log"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/d-tsuji/flower-v2/runner"
)

type Worker struct {
	ID            int
	WorkerChannel chan chan db.ExecutableTask
	Channel       chan db.ExecutableTask
	DBClient      *db.DB
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
					fmt.Printf("runner.Run() is failed. err: %v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// end worker
func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.End <- true
}
