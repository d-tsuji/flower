package repository

import (
	"log"
	"testing"

	"github.com/d-tsuji/flower/queue"
)

func TestInsertTaskDifinision(t *testing.T) {
	item := &queue.Item{
		"",
		"sample.a.id",
		"Normal",
	}

	truncateTable()
	err := InsertTaskDifinision(item)
	if err != nil {
		t.Error(err)
	}
}

func truncateTable() {
	tableName := "kr_task_status"
	if _, err := Conn.Exec("TRUNCATE TABLE " + tableName); err != nil {
		log.Fatal(err)
	}
}

func TestSelectExecTarget(t *testing.T) {
	list, err := SelectExecTarget()
	if err != nil {
		t.Errorf("%s", err)
	}
	for _, v := range *list {
		t.Log(v)
	}
}
