package main

import (
	"log"
	"time"

	"github.com/d-tsuji/flower/mock"

	"github.com/d-tsuji/flower/app"
	"github.com/d-tsuji/flower/repository"
)

var Concurrency = 3
var pollingIntervalSecond time.Duration = 10

func main() {

	// テスト用のHTTPサーバを起動し、リクエストに応じてタスクを登録
	go mock.RegisterTask()

	taskChannel := make(chan repository.KrTaskStatus, Concurrency)

	for i := 0; i < Concurrency; i++ {
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
