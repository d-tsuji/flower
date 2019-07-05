package repository

import (
	"log"
	"testing"
)

func TestInsertTaskDifinision(t *testing.T) {
	truncateTable()

	item := &Item{
		"",
		"sample.a.id",
		"Normal",
	}
	_, err := InsertTaskDefinition(item)
	if err != nil {
		t.Error(err)
	}
	item = &Item{
		"",
		"sample.b.id",
		"Normal",
	}
	_, err = InsertTaskDefinition(item)
	if err != nil {
		t.Error(err)
	}
}

func truncateTable() {
	tableName := "kr_task_status"
	if _, err := conn.Exec("TRUNCATE TABLE " + tableName); err != nil {
		log.Fatal(err)
	}
}

func TestSelectExecTarget(t *testing.T) {
	list, err := SelectExecTarget(100)
	if err != nil {
		t.Errorf("%s", err)
	}
	for _, v := range *list {
		t.Log(v)
	}
}
