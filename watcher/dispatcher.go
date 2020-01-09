package watcher

import (
	"context"
	"log"

	"github.com/d-tsuji/flower/repository"
)

var workerChannel = make(chan chan repository.ExecutableTask)

type collector struct {
	Work chan repository.ExecutableTask
}

// StartDispatcher starts workerCount number of workers. In addition,
// wait for repository.ExecutableTask to be input to channel,and dispatch the task to worker.
func StartDispatcher(ctx context.Context, workerCount int, dbClient *repository.DB) collector {
	var workers []Worker
	input := make(chan repository.ExecutableTask)
	collector := collector{Work: input}

	for i := 0; i < workerCount; i++ {
		log.Printf("[dispatcher] starting worker: %d\n", i)
		worker := Worker{
			id:            i,
			Channel:       make(chan repository.ExecutableTask),
			WorkerChannel: workerChannel,
			dbClient:      dbClient,
		}
		worker.Start(ctx)
		workers = append(workers, worker)
	}

	// start collector
	go func() {
		for {
			select {
			case work := <-input:
				worker := <-workerChannel
				worker <- work
			}
		}
	}()

	return collector
}
