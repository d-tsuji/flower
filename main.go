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

	flag.IntVar(&concurrency, "c", 3, "Concurrency (Goroutine count)")
	flag.DurationVar(&pollingIntervalSecond, "p", 20, "Polling interval")

	// テスト用のHTTPサーバを起動し、リクエストに応じてタスクを登録
	go mock.RegisterTask()

	taskChannel := make(chan repository.KrTaskStatus, concurrency)

	for i := 0; i < concurrency; i++ {
		go app.Run(taskChannel)
	}
	for {
		log.Println("main() watching...")
		err := app.WatchTask(taskChannel)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(pollingIntervalSecond * time.Second)
	}
}
