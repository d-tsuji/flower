package repository

import (
	"log"

	_ "github.com/lib/pq"
)

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
