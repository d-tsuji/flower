package testdb

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/golang/glog"
)

var (
	flowerSQL = "../assets/schema/01_createTables.sql"
	opts      = flag.String("pg_opts", "sslmode=disable", "Database options to be included when connecting to the db")
	dbName    = flag.String("db_name", "test", "The database name to be used when checking for pg connectivity")
)

// PGAvailable indicates whether a default PG database is available.
func PGAvailable() bool {
	db, err := sql.Open("postgres", getConnStr(*dbName))
	if err != nil {
		log.Printf("sql.Open(): %v", err)
		return false
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Printf("db.Ping(): %v", err)
		return false
	}
	return true
}

func getConnStr(name string) string {
	return fmt.Sprintf("database=%s %s", name, *opts)
}

// NewFlowerDB creates an empty database with the Flower schema. The database name is randomly
// generated.
// NewFlowerDB is equivalent to Default().NewFlowerDB(ctx).
func NewFlowerDB(ctx context.Context) (*sql.DB, func(context.Context), error) {
	db, done, err := newEmptyDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	sqlBytes, err := ioutil.ReadFile(flowerSQL)
	if err != nil {
		return nil, nil, err
	}

	for _, stmt := range strings.Split(Sanitize(string(sqlBytes)), ";--end") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return nil, nil, fmt.Errorf("error running statement %q: %v", stmt, err)
		}
	}
	return db, done, nil
}

// newEmptyDB creates a new, empty database.
// The returned clean up function should be called once the caller no longer
// needs the test DB.
func newEmptyDB(ctx context.Context) (*sql.DB, func(context.Context), error) {
	db, err := sql.Open("postgres", getConnStr(*dbName))
	if err != nil {
		return nil, nil, err
	}

	// Create a randomly-named database and then connect using the new name.
	name := fmt.Sprintf("trl_%v", time.Now().UnixNano())
	stmt := fmt.Sprintf("CREATE DATABASE %v", name)
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return nil, nil, fmt.Errorf("error running statement %q: %v", stmt, err)
	}
	db.Close()
	db, err = sql.Open("postgres", getConnStr(name))
	if err != nil {
		return nil, nil, err
	}

	done := func(ctx context.Context) {
		defer db.Close()
		if _, err := db.ExecContext(ctx, "DROP DATABASE %v", name); err != nil {
			glog.Warningf("Failed to drop test database %q: %v", name, err)
		}
	}

	return db, done, db.Ping()
}

// sanitize tries to remove empty lines and comments from a sql script
// to prevent them from being executed.
func Sanitize(script string) string {
	buf := &bytes.Buffer{}
	for _, line := range strings.Split(string(script), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' || strings.Index(line, "--") == 0 {
			continue // skip empty lines and comments
		}
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf.String()
}
