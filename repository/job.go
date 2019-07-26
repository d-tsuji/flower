package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	JobFlowId string // 実行時に発行される一意なフローID(初回はnil)
	TaskId    string // 実行するタスク群のID
	TaskType  string // タスク登録 or タスク監視実行
}

// ms_taskの実行順序から依存関係を決定し、kr_task_statusに実行待ちとして登録する
func InsertTaskDefinition(job *Job) (sql.Result, error) {

	statement := `
INSERT INTO
	kr_task_status(
		job_flow_id
	,	task_id
	,	job_exec_seq
	,	job_depend_exec_seq
	,	wait_mode
	,	status
	,	response_body
	,	priority
	,	create_ts
	,	start_ts
	, 	finish_ts
)
SELECT
	$1
	,	task_id
	,	exec_order
	,	depend_exec_order
	,	wait_mode
	,	status
	,	'' response_body
	,	0 priority
	,	$3 create_ts
	,	null start_ts
	,	null finish_ts
FROM  
	(
		SELECT
			task_id
		,	exec_order
		,	lag(exec_order, 1) over(partition by task_id order by exec_order) depend_exec_order
		,	wait_mode
		,	'0' status
		FROM
			ms_task t
		WHERE
			t.task_id = $2
	) res
;
`

	stmt, err := conn.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	res, err := stmt.Exec(uuid.New(), job.TaskId, time.Now())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return res, nil
}
