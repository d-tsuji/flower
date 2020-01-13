package register

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/d-tsuji/flower/repository"

	"github.com/d-tsuji/flower/repository/testdb"
	"github.com/golang/glog"
)

var allTables = []string{"ms_task_definition"}
var db *sql.DB

func Test_Get_Empty(t *testing.T) {
	router := NewRouter(&repository.DB{db})
	ts := httptest.NewServer(router)
	defer ts.Close()

	// The destination of the request is the URL of the test server.
	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}
	if "" != string(data) {
		t.Fatalf("Data Error. %v", string(data))
	}
}

func TestServer_ServeHTTP_Normal(t *testing.T) {
	cleanTestDB(db, t)
	initTestDB(db, "../testdata/test_http_1.sql", t)
	router := NewRouter(&repository.DB{db})
	ts := httptest.NewServer(router)
	defer ts.Close()

	r, err := http.Post(fmt.Sprintf("%s/%s", ts.URL, "register/dummy"), "application/json", nil)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}
	if "\"{status: succeeded, taskId: dummy}\"" != string(data) {
		t.Fatalf("Data Error. %v", string(data))
	}
}

func TestServer_ServeHTTP_AbNormal(t *testing.T) {
	cleanTestDB(db, t)
	initTestDB(db, "../testdata/test_http_1.sql", t)
	router := NewRouter(&repository.DB{db})
	ts := httptest.NewServer(router)
	defer ts.Close()

	r, err := http.Post(fmt.Sprintf("%s/%s", ts.URL, "register/error"), "application/json", nil)
	if err != nil {
		t.Fatalf("error by http.Post(). %v", err)
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("error by ioutil.ReadAll(). %v", err)
	}
	if "\"{status: failed, taskId: error, description: no tasks registered}\"" != string(data) {
		t.Fatalf("error data. %v", string(data))
	}
}

// TestMain is test helper function.
func TestMain(m *testing.M) {
	flag.Parse()
	if !testdb.PGAvailable() {
		glog.Errorf("PG not available, skipping all PG storage tests")
		return
	}

	var done func(context.Context)
	db, done = openTestDBOrDie()

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
