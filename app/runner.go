package app

import (
	"log"

	"github.com/d-tsuji/flower/repository"
)

func Run(ch chan repository.KrTaskStatus) {
	for {
		v := <-ch
		log.Printf("Task starting... : %v", v)

		// task_id, job_exec_seq から実行するrestタスクを取得
		rest, err := v.GetExecRestTaskDefinition()
		if err != nil {
			log.Fatal("%s", err)
		}

		// 管理テーブルの更新(実行可能->実行中)
		sqlResult, err := v.UpdateKrTaskStatus(&repository.Status{repository.Executable}, &repository.Status{repository.Running})
		cnt, err := sqlResult.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		if cnt == 0 {
			log.Printf("This task still started by other process. %v", v)
			continue
		}

		// rest apiを発行
		res, err := RestCall(rest)
		log.Println(string(res))
		if err != nil {
			log.Fatal("%s", err)
			// 管理テーブルの更新(実行中->異常終了)
			v.UpdateKrTaskStatus(&repository.Status{repository.Running}, &repository.Status{repository.Error})
		}

		// 管理テーブルの更新(実行中->正常終了)
		v.UpdateCompleteKrTaskStatus(&repository.Status{repository.Running}, &repository.Status{repository.Completed}, res)
		log.Printf("Task finished : %v", v)

		// 非同期で後続タスクを呼び出す。
		// 同期にするとチャネルのバッファがいっぱいのときに呼び出し元のrunner#Runが終了できなくなり、デッドロックになる
		go WatchTask(ch)
	}
}
