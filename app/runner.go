package app

import (
	"log"

	"github.com/d-tsuji/flower/repository"
)

func Run(ch chan repository.KrTaskStatus) {
	for {
		// 実行可能なタスクが存在する場合
		v := <-ch
		log.Printf("Task starting... : %v", v)

		// task_id, job_exec_seq から実行するrestタスクを取得
		rest, err := repository.GetExecRestTaskDefinision(&v)
		if err != nil {
			log.Fatal("%s", err)
		}

		// 管理テーブルの更新(実行可能->実行中)
		sqlResult, err := repository.UpdateKrTaskStatus(&repository.Status{"1"}, &repository.Status{"2"}, &v)
		cnt, err := sqlResult.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		if cnt == 0 {
			log.Printf("This task still started by other process. v")
			continue
		}

		// rest apiを発行
		res, err := RestCall(rest)
		log.Println(string(res))
		if err != nil {
			log.Fatal("%s", err)
			// 管理テーブルの更新(実行中->異常終了)
			repository.UpdateKrTaskStatus(&repository.Status{"2"}, &repository.Status{"9"}, &v)
		}

		// 管理テーブルの更新(実行中->正常終了)
		repository.UpdateKrTaskStatus(&repository.Status{"2"}, &repository.Status{"3"}, &v)
		WatchTask(ch)
	}
}
