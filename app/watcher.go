package app

import (
	"go.uber.org/zap"

	"github.com/d-tsuji/flower/repository"
)

func WatchTaskLimit(ch chan<- repository.Task, limit int) error {
	logger, _ := zap.NewDevelopment()

	logger.Info("Task Watching.")
	list, err := repository.SelectExecTarget(limit)
	if err != nil {
		logger.Error("Error getting target tasks", zap.Error(err))
		return err
	}

	for _, v := range *list {
		logger.Info("Executable task found. Put channel. " + v.String())
		// 管理テーブルの更新(実行待ち->実行可能)
		v.UpdateKrTaskStatus(repository.WaitExecute, repository.Executable)
		ch <- v
	}
	return nil
}

func WatchTask(ch chan<- repository.Task) error {
	logger, _ := zap.NewDevelopment()

	if err := WatchTaskLimit(ch, 1000); err != nil {
		logger.Error("Error getting target tasks", zap.Error(err))
		return err
	}
	return nil
}
