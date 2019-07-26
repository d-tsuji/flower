package repository

import (
	"database/sql"
	"log"
	"time"
)

type Task struct {
	JobFlowId    string `db:"job_flow_id"`
	TaskId       string `db:"task_id"`
	JobExecSeq   int64  `db:"job_exec_seq"`
	ResponseBody string `db:"response_body"`
}

type RestTask struct {
	Endpoint        string `db:"endpoint"`
	Method          string `db:"method"`
	ExtendParameter string `db:"extend_parameters"`
	// TODO: パラメータなど
}

// TaskIDに紐づくRestタスク一覧を取得する
// 呼び出す際に用いるパラメータは依存している直前のタスクが完了したときのレスポンスボディのパラメータ
func (task *Task) GetExecRestTaskDefinition() (*RestTask, error) {

	var endpoint string
	var method string
	var extendParameter string
	query := `
SELECT
	t.endpoint
,	t.method
,	t.extend_parameters
FROM
	ms_task t
WHERE
	t.task_id = $1
AND	t.exec_order = $2
;
`
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
