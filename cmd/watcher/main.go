package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/d-tsuji/flower-v2/watcher"
)

const (
	WORKER_COUNT            = 5
	POLLING_INTERVAL_SECOND = 5
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbClient, err := db.New(&db.Opt{Password: "flower"})
	defer dbClient.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("[watcher] postgres initialize error: %v\n", err))
	}
	// start up worker pool
	collector := watcher.StartDispatcher(ctx, WORKER_COUNT, dbClient)

	w := watcher.NewWatcherTask(dbClient, make(chan db.ExecutableTask))
	tic := time.NewTicker(POLLING_INTERVAL_SECOND * time.Second)
	go func() {
		for {
			select {
			case <-tic.C:
				log.Printf("[watcher] watching task...\n")
				if err := w.WatchTask(ctx); err != nil {
					log.Printf("[watcher] watcher task error: %+v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case t := <-w.ExecTaskCh:
			collector.Work <- t
		case <-ctx.Done():
			return
		}
	}
}
