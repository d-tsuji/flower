package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/d-tsuji/flower-v2/watcher"
)

const (
	WORKER_COUNT            = 5
	DB_CONCURRENCY          = WORKER_COUNT
	POLLING_INTERVAL_SECOND = 5
)

func main() {
	dbuser := flag.String("dbuser", "", "postgres user")
	dbpass := flag.String("dbpass", "", "postgres user password")
	dbhost := flag.String("dbhost", "", "postgres host")
	dbport := flag.String("dbport", "", "postgres port")
	dbname := flag.String("dbname", "", "postgres database name")
	flag.Parse()

	dbClient, err := db.New(&db.Opt{
		DBName:   *dbname,
		User:     *dbuser,
		Password: *dbpass,
		Host:     *dbhost,
		Port:     *dbport,
	})
	defer dbClient.Close()
	if err != nil {
		log.Fatal(fmt.Sprintf("[watcher] postgres initialize error: %v\n", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start up worker pool
	collector := watcher.StartDispatcher(ctx, WORKER_COUNT, dbClient)
	w := watcher.NewWatcherTask(dbClient, make(chan db.ExecutableTask))
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
