package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func init() {

	connStr := "postgres://dev:dev@0.0.0.0/dev?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err.Error)
	}
	Conn = conn
}
