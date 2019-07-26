package app

import (
	"go.uber.org/zap"

	"github.com/d-tsuji/flower/repository"
)

func Run(ch chan repository.Task) {
	logger, _ := zap.NewDevelopment()

	for {
		v := <-ch
		logger.Info("Task starting.")

		// task_id, job_exec_seq から実行するrestタスクを取得
		rest, err := v.GetExecRestTaskDefinition()
		if err != nil {
			logger.Warn("None get tasks", zap.Error(err))
			continue
		}

		// 管理テーブルの更新(実行可能->実行中)
		sqlResult, err := v.UpdateKrTaskStatus(repository.Executable, repository.Running)
		cnt, err := sqlResult.RowsAffected()
		if err != nil {
			logger.Warn("An unexpected error has occurred", zap.Error(err))
			continue
		}
		if cnt == 0 {
			logger.Warn("This task still started by other process. " + v.String())
			//log.Printf("This task still started by other process. %v", v)
			continue
		}

		// rest apiを発行
		res, err := RestCall(rest)
		logger.Info(string(res))
		if err != nil {
			// 管理テーブルの更新(実行中->異常終了)
			v.UpdateKrTaskStatus(repository.Running, repository.Error)
			logger.Warn("An unexpected error has occurred", zap.Error(err))
			continue
		}

		// 管理テーブルの更新(実行中->正常終了)
		v.UpdateCompleteKrTaskStatus(repository.Running, repository.Completed, res)
		logger.Info("Task finished.")
		//logger.Info("Task finished : %v", v)

		// 非同期で後続タスクを呼び出す。
		// 同期にするとチャネルのバッファがいっぱいのときに呼び出し元のrunner#Runが終了できなくなり、デッドロックになる
		go WatchTask(ch)
	}
}
