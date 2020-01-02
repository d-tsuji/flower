package watcher

import (
	"fmt"

	"github.com/d-tsuji/flower-v2/db"
)

var WorkerChannel = make(chan chan db.ExecutableTask)

type Collector struct {
	Work chan db.ExecutableTask
	End  chan bool
}

func StartDispatcher(workerCount int, dbClient *db.DB) Collector {
	var i int
	var workers []Worker
	input := make(chan db.ExecutableTask) // channel to recieve work
	end := make(chan bool)                // channel to spin down workers
	collector := Collector{Work: input, End: end}

	for i < workerCount {
		i++
		fmt.Println("starting worker: ", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan db.ExecutableTask),
			WorkerChannel: WorkerChannel,
			End:           make(chan bool),
			DBClient:      dbClient,
		}
		worker.Start()
		workers = append(workers, worker) // store worker
	}

	// start collector
	go func() {
		for {
			select {
			case <-end:
				for _, w := range workers {
					w.Stop() // stop worker
				}
				return
			case work := <-input:
				worker := <-WorkerChannel // wait for available channel
				worker <- work            // dispatch work to worker
			}
		}
	}()

	return collector
}
