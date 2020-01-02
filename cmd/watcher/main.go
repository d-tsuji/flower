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
	if err != nil {
		log.Fatal(fmt.Sprintf("postgres initialize error: %v\n", err))
	}
	collector := watcher.StartDispatcher(WORKER_COUNT, dbClient) // start up worker pool

	w := watcher.NewWatcherTask(dbClient, make(chan db.ExecutableTask))
	go func() {
		for {
			fmt.Println("Start watcher query.")
			if err := w.WatchTask(ctx); err != nil {
				fmt.Printf("watcher task error: %+v\n", err)
			}
			time.Sleep(POLLING_INTERVAL_SECOND * time.Second)
		}
	}()

	for t := range w.ExecTaskCh {
		collector.Work <- t
	}
}
