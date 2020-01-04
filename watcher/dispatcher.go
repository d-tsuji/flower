package watcher

import (
	"context"
	"log"

	"github.com/d-tsuji/flower-v2/db"
)

var WorkerChannel = make(chan chan db.ExecutableTask)

type Collector struct {
	Work chan db.ExecutableTask
}

func StartDispatcher(ctx context.Context, workerCount int, dbClient *db.DB) Collector {
	var i int
	var workers []Worker
	input := make(chan db.ExecutableTask)
	collector := Collector{Work: input}

	for i < workerCount {
		i++
		log.Printf("[dispatcher] starting worker: %d\n", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan db.ExecutableTask),
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
