package repository

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/d-tsuji/flower/repository/testdb"
	"github.com/golang/glog"
)

var allTables = []string{"ms_task_definition", "kr_task_stat"}
var xdb *sql.DB

func TestDB_GetTaskProgramName(t *testing.T) {
	db := DB{xdb}
	cleanTestDB(xdb, t)
	initTestDB(xdb, "../testdata/test_ms.sql", t)
	e := &ExecutableTask{
		TaskFlowId:  "test_0",
		TaskExecSeq: 1,
		TaskId:      "test_0",
		TaskSeq:     1,
		Params:      nil,
	}
	res, err := db.GetTaskProgramName(context.TODO(), e)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if res != "1" {
		log.Fatalf("value is not equal: expected=%s, got=%s", "1", res)
	}
}

func TestDB_GetExecutableTask_Single(t *testing.T) {
	db := DB{xdb}
	tests := []struct {
		inputDataPath string
		expect        *ExecutableTask
	}{
		{
			inputDataPath: "../testdata/test_1.sql",
			expect: &ExecutableTask{
				TaskFlowId:  "test_flow_1",
				TaskExecSeq: 1,
				TaskId:      "test_1",
				TaskSeq:     1,
				Params:      map[string]string{},
			},
		},
		{
			inputDataPath: "../testdata/test_2.sql",
			expect: &ExecutableTask{
				TaskFlowId:  "test_flow_2",
				TaskExecSeq: 2,
				TaskId:      "test_2",
				TaskSeq:     2,
				Params:      map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		cleanTestDB(xdb, t)
		initTestDB(xdb, tt.inputDataPath, t)
		got, err := db.GetExecutableTask(context.TODO(), 1)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if !equalExecutableTask(tt.expect, &got[0]) {
			log.Fatalf("value is not equal: expected=%+v, got=%+v", tt.expect, got)
		}
	}
}

func TestDB_GetExecutableTask_Multi(t *testing.T) {
	db := DB{xdb}
	tests := []struct {
		inputDataPath string
		expect        []ExecutableTask
	}{
		{
			inputDataPath: "../testdata/test_4.sql",
			expect: []ExecutableTask{
				{
					TaskFlowId:  "test_flow_4_1",
					TaskExecSeq: 1,
					TaskId:      "test_4_1",
					TaskSeq:     1,
					Params:      map[string]string{},
				},
				{
					TaskFlowId:  "test_flow_4_2",
					TaskExecSeq: 1,
					TaskId:      "test_4_2",
					TaskSeq:     1,
					Params:      map[string]string{},
				},
			},
		},
		{
			inputDataPath: "../testdata/test_5.sql",
			expect: []ExecutableTask{
				{
					TaskFlowId:  "test_flow_5_1",
					TaskExecSeq: 1,
					TaskId:      "test_5",
					TaskSeq:     1,
					Params:      map[string]string{},
				},
			},
		},
		{
			inputDataPath: "../testdata/test_6.sql",
			expect: []ExecutableTask{
				{
					TaskFlowId:  "test_flow_6_2",
					TaskExecSeq: 2,
					TaskId:      "test_6",
					TaskSeq:     2,
					Params:      map[string]string{},
				},
			},
		},
	}

	for _, tt := range tests {
		cleanTestDB(xdb, t)
		initTestDB(xdb, tt.inputDataPath, t)
		got, err := db.GetExecutableTask(context.TODO(), 100)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if len(got) != len(tt.expect) {
			log.Fatalf("error executable count: expected=%+v, got=%+v", len(tt.expect), len(got))
		}
		for i, expect := range tt.expect {
			if !equalExecutableTask(&expect, &got[i]) {
				log.Fatalf("value is not equal: expected=%+v, got=%+v", expect, got[i])
			}
		}
	}
}

func TestDB_GetExecutableTask_Zero(t *testing.T) {
	db := DB{xdb}
	cleanTestDB(xdb, t)
	initTestDB(xdb, "../testdata/test_3.sql", t)
	got, err := db.GetExecutableTask(context.TODO(), 1)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if len(got) != 0 {
		log.Fatalf("error executable count: expected=%+v, got=%+v", 0, len(got))
	}
}

func TestDB_UpdateExecutableTasksRunning(t *testing.T) {
	db := DB{xdb}
	ctx := context.TODO()

	tests := []struct {
		inputDataPath       string
		inputExecutableTask *ExecutableTask
		expectExecStatus    int
		fn                  func(context.Context, *ExecutableTask) (bool, error)
	}{
		{
			inputDataPath: "../testdata/test_update.sql",
			inputExecutableTask: &ExecutableTask{
				TaskFlowId:  "test_flow_update_1",
				TaskExecSeq: 1,
				TaskId:      "test_update",
				TaskSeq:     1,
				Params:      map[string]string{},
			},
			expectExecStatus: ExecStatusRunning,
			fn:               db.UpdateExecutableTasksRunning,
		},
	}

	for _, tt := range tests {
		cleanTestDB(xdb, t)
		initTestDB(xdb, tt.inputDataPath, t)
		//_, err := db.UpdateExecutableTasksRunning(ctx, tt.inputExecutableTask)
		_, err := tt.fn(ctx, tt.inputExecutableTask)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		var gotExecStatus int
		err = db.QueryRowContext(ctx, "SELECT exec_status FROM kr_task_stat WHERE task_flow_id = $1 AND task_exec_seq = $2;",
			tt.inputExecutableTask.TaskFlowId, tt.inputExecutableTask.TaskExecSeq).Scan(&gotExecStatus)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if gotExecStatus != tt.expectExecStatus {
			log.Fatalf("error executable count: expected=%+v, got=%+v", tt.expectExecStatus, gotExecStatus)
		}
	}
}

func TestDB_UpdateExecutableTasksRunningMulti(t *testing.T) {
	db := DB{xdb}
	ctx := context.TODO()

	tests := []struct {
		inputDataPath       string
		inputExecutableTask *ExecutableTask
		expectExecStatus    int
		fn                  func(context.Context, *ExecutableTask) (bool, error)
	}{
		{
			inputDataPath: "../testdata/test_update.sql",
			inputExecutableTask: &ExecutableTask{
				TaskFlowId:  "test_flow_update_1",
				TaskExecSeq: 1,
				TaskId:      "test_update",
				TaskSeq:     1,
				Params:      map[string]string{},
			},
			expectExecStatus: ExecStatusRunning,
			fn:               db.UpdateExecutableTasksRunning,
		},
	}

	for _, tt := range tests {
		cleanTestDB(xdb, t)
		initTestDB(xdb, tt.inputDataPath, t)
		if _, err := tt.fn(ctx, tt.inputExecutableTask); err != nil {
			log.Fatalf("%+v", err)
		}

		ok, err := tt.fn(ctx, tt.inputExecutableTask)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if ok {
			log.Fatalf("error task (%v) updated multiple times", tt.inputExecutableTask)
		}
	}
}

// logic helper function
func equalExecutableTask(expect, got *ExecutableTask) bool {
	if expect.TaskFlowId != got.TaskFlowId {
		return false
	}
	if expect.TaskExecSeq != got.TaskExecSeq {
		return false
	}
	if expect.TaskId != got.TaskId {
		return false
	}
	if expect.TaskSeq != got.TaskSeq {
		return false
	}
	if !reflect.DeepEqual(expect.Params, got.Params) {
		return false
	}
	return true
}

// TestMain is test helper function.
func TestMain(m *testing.M) {
	flag.Parse()
	if !testdb.PGAvailable() {
		glog.Errorf("PG not available, skipping all PG storage tests")
		return
	}

	var done func(context.Context)
	xdb, done = openTestDBOrDie()

	status := m.Run()
	done(context.Background())
	os.Exit(status)
}

func openTestDBOrDie() (*sql.DB, func(context.Context)) {
	db, done, err := testdb.NewFlowerDB(context.TODO())
	if err != nil {
		panic(err)
	}
	return db, done
}

func cleanTestDB(db *sql.DB, t *testing.T) {
	t.Helper()
	for _, table := range allTables {
		if _, err := db.ExecContext(context.TODO(), fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			t.Fatal(fmt.Sprintf("Failed to delete rows in %s: %v", table, err))
		}
	}
}

func initTestDB(db *sql.DB, input string, t *testing.T) {
	t.Helper()
	sqlBytes, err := ioutil.ReadFile(input)
	if err != nil {
		t.Fatalf("error input file(%s) cannot read: %+v", input, err)
	}

	for _, stmt := range strings.Split(testdb.Sanitize(string(sqlBytes)), ";--end") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.ExecContext(context.TODO(), stmt); err != nil {
			t.Fatalf("error running statement %q: %v", stmt, err)
		}
	}
}
