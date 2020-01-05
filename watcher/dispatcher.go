package watcher

import (
	"context"
	"log"

	"github.com/d-tsuji/flower/repository"
)

var WorkerChannel = make(chan chan repository.ExecutableTask)

type Collector struct {
	Work chan repository.ExecutableTask
}

func StartDispatcher(ctx context.Context, workerCount int, dbClient *repository.DB) Collector {
	var i int
	var workers []Worker
	input := make(chan repository.ExecutableTask)
	collector := Collector{Work: input}

	for i < workerCount {
		i++
		log.Printf("[dispatcher] starting worker: %d\n", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan repository.ExecutableTask),
			WorkerChannel: WorkerChannel,
			DBClient:      dbClient,
		}
		worker.Start(ctx)
		workers = append(workers, worker)
	}

	// start collector
	go func() {
		for {
			select {
			case work := <-input:
				worker := <-WorkerChannel
				worker <- work
			}
		}
	}()

	return collector
}
