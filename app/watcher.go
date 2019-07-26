package app

import (
	"log"

	"go.uber.org/zap"

	"github.com/d-tsuji/flower/repository"
)

func WatchTaskLimit(ch chan<- repository.Task, concurrency int) error {

	logger, _ := zap.NewDevelopment()

	logger.Info("Task Watching.")
	list, err := repository.SelectExecTarget(concurrency)
	if err != nil {
		log.Fatal("%s", err)
		return err
	}

	for _, v := range *list {
		log.Printf("Executable task found. Put channel. %v", v)
		// 管理テーブルの更新(実行待ち->実行可能)
		v.UpdateKrTaskStatus(repository.WaitExecute, repository.Executable)
		ch <- v
	}
	return nil
}

func WatchTask(ch chan<- repository.Task) error {
	if err := WatchTaskLimit(ch, 1000); err != nil {
		return err
	}
	return nil
}
