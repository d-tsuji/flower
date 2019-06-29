package main

import (
	"time"

	"github.com/d-tsuji/flower/app"
	"github.com/d-tsuji/flower/repository"
)

func main() {

	taskChannel := make(chan repository.KrTaskStatus, 10)

	go app.Run(taskChannel)
	for {
		_ = app.WatchTask(taskChannel)
		time.Sleep(5 * time.Second)
	}

}
