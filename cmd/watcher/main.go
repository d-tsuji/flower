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
		log.Fatal(fmt.Sprintf("postgres initialize error: %v\n", err))
	}
	collector := watcher.StartDispatcher(ctx, WORKER_COUNT, dbClient) // start up worker pool

	w := watcher.NewWatcherTask(dbClient, make(chan db.ExecutableTask))
	tic := time.NewTicker(POLLING_INTERVAL_SECOND * time.Second)
	go func() {
		for {
			select {
			case <-tic.C:
				fmt.Println("watching task...")
				if err := w.WatchTask(ctx); err != nil {
					fmt.Printf("watcher task error: %+v\n", err)
				}
			}
		}
	}()

	for t := range w.ExecTaskCh {
		collector.Work <- t
	}
}
