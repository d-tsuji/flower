package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/d-tsuji/flower-v2/register"
)

func main() {
	dbuser := flag.String("dbuser", "", "postgres user")
	dbpass := flag.String("dbpass", "", "postgres user password")
	dbhost := flag.String("dbhost", "", "postgres host")
	dbport := flag.String("dbport", "", "postgres port")
	dbname := flag.String("dbname", "", "postgres database name")
	flag.Parse()

	dbClient, err := db.New(&db.Opt{
		DBName:   *dbname,
		User:     *dbuser,
		Password: *dbpass,
		Host:     *dbhost,
		Port:     *dbport,
	})
	defer dbClient.Close()
	if err != nil {
		log.Printf("[register] postgres initialize error: %v\n", err)
	}
	s := register.NewServer(dbClient)
	http.HandleFunc("/", s.ServeHTTP)
	address := "0.0.0.0:8000"
	log.Printf("[register] Starting server on address: %s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
