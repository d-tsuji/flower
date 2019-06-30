package main

import (
	"log"
	"time"

	"github.com/d-tsuji/flower/mock"

	"github.com/d-tsuji/flower/app"
	"github.com/d-tsuji/flower/repository"
)

var concurrency = 3
var pollingIntervalSecond time.Duration = 5

func main() {

	// テスト用のHTTPサーバを起動し、リクエストに応じてタスクを登録
	go mock.StartServer()

	taskChannel := make(chan repository.KrTaskStatus, concurrency)

	for i := 0; i < 10; i++ {
		go app.Run(taskChannel)
	}
	for {
		err := app.WatchTask(taskChannel)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(pollingIntervalSecond * time.Second)
	}
}
