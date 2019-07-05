package app

import (
	"log"

	"github.com/d-tsuji/flower/repository"
)

func WatchTaskLimit(ch chan<- repository.KrTaskStatus, concurrency int) error {

	log.Println("Task Watching...")
	list, err := repository.SelectExecTarget(concurrency)
	if err != nil {
		log.Fatal("%s", err)
		return err
	}

	for _, v := range *list {
		log.Printf("Executable task found. Put channel. %v", v)
		// 管理テーブルの更新(実行待ち->実行可能)
		v.UpdateKrTaskStatus(&repository.Status{repository.WaitExecute}, &repository.Status{repository.Executable})
		ch <- v
	}
	return nil
}

func WatchTask(ch chan<- repository.KrTaskStatus) error {
	if err := WatchTaskLimit(ch, 1000); err != nil {
		return err
	}
	return nil
}
