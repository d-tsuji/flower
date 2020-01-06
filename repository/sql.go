package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Row defines a common interface between sql.Row and sql.Rows(!)
type Row interface {
	Scan(dest ...interface{}) error
}

func readTask(row Row) (*task, error) {
	var (
		taskId       string
		taskSeq      int
		program      string
		taskPriority int
		param1Key    sql.NullString
		param1Value  sql.NullString
		param2Key    sql.NullString
		param2Value  sql.NullString
		param3Key    sql.NullString
		param3Value  sql.NullString
		param4Key    sql.NullString
		param4Value  sql.NullString
		param5Key    sql.NullString
		param5Value  sql.NullString
	)
	err := row.Scan(
		&taskId,
		&taskSeq,
		&program,
		&taskPriority,
		&param1Key,
		&param1Value,
		&param2Key,
		&param2Value,
		&param3Key,
		&param3Value,
		&param4Key,
		&param4Value,
		&param5Key,
		&param5Value,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("rows scan error: %v", err))
	}
	pm := make(map[string]string)
	addMapIfNullString(param1Key, param1Value, &pm)
	addMapIfNullString(param2Key, param2Value, &pm)
	addMapIfNullString(param3Key, param3Value, &pm)
	addMapIfNullString(param4Key, param4Value, &pm)
	addMapIfNullString(param5Key, param5Value, &pm)

	return &task{
		TaskId:       taskId,
		TaskSeq:      taskSeq,
		Program:      program,
		TaskPriority: taskPriority,
		Params:       pm,
	}, nil
}

func readExecutableTask(row Row) (*ExecutableTask, error) {
	var (
		taskFlowId  string
		taskExecSeq int
		taskId      string
		taskSeq     int
		data        string
	)
	if err := row.Scan(&taskFlowId, &taskExecSeq, &taskId, &taskSeq, &data); err != nil {
		return nil, errors.New(fmt.Sprintf("rows scan error: %v", err))
	}
	params := make(map[string]string)
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		return nil, errors.New(fmt.Sprintf("json unmarshal error: %v", err))
	}
	return &ExecutableTask{
		TaskFlowId:  taskFlowId,
		TaskExecSeq: taskExecSeq,
		TaskId:      taskId,
		TaskSeq:     taskSeq,
		Params:      params,
	}, nil
}

func addMapIfNullString(key sql.NullString, value sql.NullString, m *map[string]string) {
	if key.Valid {
		(*m)[key.String] = value.String
	}
}
