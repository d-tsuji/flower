package app

import (
	"log"

	"github.com/d-tsuji/flower/repository"
)

func WatchTask(ch chan<- repository.KrTaskStatus) error {

	log.Println("Task Watching...")
	list, err := repository.SelectExecTarget()
	if err != nil {
		log.Fatal("%s", err)
		return err
	}

	for _, v := range *list {
		log.Printf("Executable task found. Put channel. %v", v)
		// 管理テーブルの更新(実行待ち->実行可能)
		repository.UpdateKrTaskStatus(&repository.Status{"0"}, &repository.Status{"1"}, &v)
		ch <- v
	}
	return nil
}
