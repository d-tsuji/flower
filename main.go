package main

import (
	"flag"
	"log"
	"time"

	"github.com/d-tsuji/flower/mock"

	"github.com/d-tsuji/flower/app"
	"github.com/d-tsuji/flower/repository"
)

var concurrency int
var pollingIntervalSecond time.Duration

func main() {

	flag.IntVar(&concurrency, "c", 3, "Concurrency")
	flag.DurationVar(&pollingIntervalSecond, "p", 20, "Polling interval second")

	// テスト用のHTTPサーバを起動し、リクエストに応じてタスクを登録
	go mock.RegisterTask()

	taskChannel := make(chan repository.KrTaskStatus, concurrency)

	for i := 0; i < concurrency; i++ {
		go app.Run(taskChannel)
	}

	for {
		err := app.WatchTaskLimit(taskChannel, concurrency)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(pollingIntervalSecond * time.Second)
	}
}
