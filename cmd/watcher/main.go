package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/d-tsuji/lightenv"

	"github.com/d-tsuji/flower/repository"
	"github.com/d-tsuji/flower/watcher"
)

const (
	// WORKER_COUNT is the number of workers in the worker-pool model.
	// It is the same as the number of tasks that can be executed in parallel.
	WORKER_COUNT = 10

	// DB_CONCURRENCY is the number of tasks to fetch executable
	// tasks and put to workers. It is usually recommended to be
	// the WORKER_COUNT >= 2 * DB_CONCURRENCY.
	DB_CONCURRENCY = 5

	// POLLING_INTERVAL_SECOND is the interval for monitoring
	// tasks that can be executed from Database.
	POLLING_INTERVAL_SECOND = 5
)

type flowerEnv struct {
	DbUser string `name:"DB_USER" required:"true"`
	DbPass string `name:"DB_PASS" required:"true"`
	DbHost string `name:"DB_HOST" required:"true"`
	DbPort string `name:"DB_PORT" required:"true"`
	DbName string `name:"DB_NAME" required:"true"`
}

var f flowerEnv

func init() {
	if err := lightenv.Process(&f); err != nil {
		log.Fatalf(fmt.Sprintf("[watcher] environment initialize error: %v\n", err))
	}
}

func main() {
	dbClient, err := repository.New(&repository.Opt{
		DBName:   f.DbName,
		User:     f.DbUser,
		Password: f.DbPass,
		Host:     f.DbHost,
		Port:     f.DbPort,
	})
	defer dbClient.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("[watcher] postgres initialize error: %v\n", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start up worker pool
	collector := watcher.StartDispatcher(ctx, WORKER_COUNT, dbClient)
	w := watcher.NewWatcherTask(dbClient, make(chan repository.ExecutableTask))
	tic := time.NewTicker(POLLING_INTERVAL_SECOND * time.Second)
	go func() {
		for {
			select {
			case <-tic.C:
				log.Printf("[watcher] watching task...\n")
				if err := w.WatchTask(ctx, DB_CONCURRENCY); err != nil {
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
