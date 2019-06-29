package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/d-tsuji/flower/queue"
)

type TaskDifinition struct {
}

// ms_taskの実行順序から依存関係を決定し、kr_task_statusに実行待ちとして登録する
func InsertTaskDifinision(item *queue.Item) error {

	statement :=
		"INSERT INTO kr_task_status (job_flow_id, task_id, job_exec_seq, job_depend_exec_seq, wait_mode, status, response_body, priority, create_ts, update_ts)" +
			"SELECT $1, task_id, exec_order, depend_exec_order, wait_mode, status, '', 0 ,$3, $4 " +
			"FROM " +
			"(SELECT task_id, exec_order, lag(exec_order, 1) over(partition by task_id order by exec_order) depend_exec_order, wait_mode, '0' status " +
			"FROM ms_task t WHERE t.task_id = $2) res"
	stmt, err := Conn.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer stmt.Close()
	if _, err := stmt.Exec(uuid.New(), item.TaskId, time.Now(), time.Now()); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

type RestTask struct {
	Endpoint string `db:"endpoint"`
	Method   string `db:"method"`
	// TODO: パラメータなど
}

// TaskIDに紐づくタスク一覧を取得する
func GetExecRestTaskDefinision(item *KrTaskStatus) (*RestTask, error) {

	var endpoint string
	var method string
	query := "select t.endpoint, t.method from ms_task t where t.task_id = $1 and t.exec_order = $2"
	err := Conn.QueryRow(query, item.TaskId, item.JobExecSeq).Scan(&endpoint, &method)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &RestTask{
		endpoint,
		method,
	}, nil
}

type KrTaskStatus struct {
	JobFlowId    string `db:"job_flow_id"`
	TaskId       string `db:"task_id"`
	JobExecSeq   int64  `db:"job_exec_seq"`
	ResponseBody string `db:"response_body"`
}

func SelectExecTarget() (*[]KrTaskStatus, error) {

	list := make([]KrTaskStatus, 0)
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
		;
	`

	stmt, err := Conn.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	for stmt.Next() {
		var job_flow_id string
		var task_id string
		var job_exec_seq int64
		var response_body string

		if err := stmt.Scan(&job_flow_id, &task_id, &job_exec_seq, &response_body); err != nil {
			log.Fatal(err)
			return nil, err
		}
		list = append(list, KrTaskStatus{
			job_flow_id,
			task_id,
			job_exec_seq,
			response_body,
		})
	}

	return &list, nil
}

type Status struct {
	S string
}

func UpdateKrTaskStatus(fromStat *Status, toStat *Status, task *KrTaskStatus) (sql.Result, error) {
	statement := "UPDATE kr_task_status SET status = $1, update_ts = $2 WHERE job_flow_id = $3 AND task_id = $4 AND job_exec_seq = $5 AND status = $6"
	stmt, err := Conn.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer stmt.Close()
	cnt, err := stmt.Exec(toStat.S, time.Now(), task.JobFlowId, task.TaskId, task.JobExecSeq, fromStat.S)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cnt, nil
}
