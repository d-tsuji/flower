package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

// ms_taskの実行順序から依存関係を決定し、kr_task_statusに実行待ちとして登録する
func InsertTaskDefinition(item *Item) (sql.Result, error) {

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
	res, err := stmt.Exec(uuid.New(), item.TaskId, time.Now())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return res, nil
}

type RestTask struct {
	Endpoint        string `db:"endpoint"`
	Method          string `db:"method"`
	ExtendParameter string `db:"extend_parameters"`
	// TODO: パラメータなど
}

// TaskIDに紐づくタスク一覧を取得する
func (task *Task) GetExecRestTaskDefinition() (*RestTask, error) {

	var endpoint string
	var method string
	var extendParameter string
	query := "select t.endpoint, t.method, t.extend_parameters from ms_task t where t.task_id = $1 and t.exec_order = $2"
	err := conn.QueryRow(query, task.TaskId, task.JobExecSeq).Scan(&endpoint, &method, &extendParameter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &RestTask{
		endpoint,
		method,
		extendParameter,
	}, nil
}

type Task struct {
	JobFlowId    string `db:"job_flow_id"`
	TaskId       string `db:"task_id"`
	JobExecSeq   int64  `db:"job_exec_seq"`
	ResponseBody string `db:"response_body"`
}

// 依存タスクが完了している実行待ちタスクを取得する
func SelectExecTarget(limit int) (*[]Task, error) {

	list := make([]Task, 0)
	query := `
select
	base.job_flow_id
,	base.task_id
,	base.job_exec_seq
,	coalesce(dep.response_body, '') response_body
from
	kr_task_status base
left join
	kr_task_status dep
	on	1=1
		and	base.job_depend_exec_seq = dep.job_exec_seq
		and	base.task_id = dep.task_id
		and	base.job_flow_id = dep.job_flow_id
where	1=1
	and	coalesce(dep.status, '3') = '3'
	and	base.status in ('0')
order by
	base.priority
limit
	$1
;
`

	stmt, err := conn.Query(query, limit)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	for stmt.Next() {
		var jobFlowId string
		var taskId string
		var jobExecSeq int64
		var responseBody string

		if err := stmt.Scan(&jobFlowId, &taskId, &jobExecSeq, &responseBody); err != nil {
			log.Fatal(err)
			return nil, err
		}
		list = append(list, Task{
			jobFlowId,
			taskId,
			jobExecSeq,
			responseBody,
		})
	}

	return &list, nil
}

func (task *Task) UpdateKrTaskStatus(fromStat StatusType, toStat StatusType) (sql.Result, error) {
	statement := `
UPDATE
	kr_task_status
SET
	status = $1
,	start_ts = $2
WHERE
	job_flow_id = $3
AND	task_id = $4
AND	job_exec_seq = $5
AND	status = $6
`
	stmt, err := conn.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	cnt, err := stmt.Exec(toStat, time.Now(), task.JobFlowId, task.TaskId, task.JobExecSeq, fromStat)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cnt, nil
}

func (task *Task) UpdateCompleteKrTaskStatus(fromStat StatusType, toStat StatusType, res []byte) (sql.Result, error) {
	statement := `
UPDATE
	kr_task_status
SET
	status = $1
,	finish_ts = $2
,	response_body = $7
WHERE
	job_flow_id = $3
AND	task_id = $4
AND	job_exec_seq = $5
AND	status = $6
`
	stmt, err := conn.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	cnt, err := stmt.Exec(toStat, time.Now(), task.JobFlowId, task.TaskId, task.JobExecSeq, fromStat, string(res))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cnt, nil
}

type Item struct {
	JobFlowId string // 実行時に発行される一意なフローID(初回はnil)
	TaskId    string // 実行するタスク群のID
	TaskType  string // タスク登録 or タスク監視実行
}
