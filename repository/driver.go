package repository

import (
	"database/sql"
)

var (
	conn *sql.DB
)

func init() {
	var err error
	conn, err = sql.Open("postgres", "user=postgres dbname=dev password=postgres host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}
}
