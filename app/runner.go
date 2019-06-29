package app

import (
	"log"

	"github.com/d-tsuji/flower/repository"
)

func Run(ch chan repository.KrTaskStatus) {
	for {
		// 実行可能なタスクが存在する場合(チャネルで呼び出し可能となる)
		v := <-ch
		log.Println(v)

		// task_id, job_exec_seq から実行するrestタスクを特定してロックを取得
		rest, err := repository.GetExecRestTaskDefinision(&v)
		if err != nil {
			log.Fatal("%s", err)
		}

		// 管理テーブルの更新(実行待ち->実行中)
		repository.UpdateKrTaskStatus(&repository.Status{"1"}, &v)

		// rest apiを発行
		res, err := RestCall(rest)
		log.Println(string(res))
		if err != nil {
			log.Fatal("%s", err)
			// 管理テーブルの更新(実行中->異常終了)
			repository.UpdateKrTaskStatus(&repository.Status{"9"}, &v)
		}

		// 管理テーブルの更新(実行中->正常終了)
		repository.UpdateKrTaskStatus(&repository.Status{"3"}, &v)
		WatchTask(ch)
	}
}
