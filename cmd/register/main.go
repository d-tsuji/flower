package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/d-tsuji/flower/register"
	"github.com/d-tsuji/flower/repository"
)

func main() {
	dbuser := flag.String("dbuser", "", "postgres user")
	dbpass := flag.String("dbpass", "", "postgres user password")
	dbhost := flag.String("dbhost", "", "postgres host")
	dbport := flag.String("dbport", "", "postgres port")
	dbname := flag.String("dbname", "", "postgres database name")
	webhost := flag.String("webhost", "localhost", "web server host")
	webport := flag.String("webport", "8000", "web server port")
	flag.Parse()

	dbClient, err := repository.New(&repository.Opt{
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
	s := register.NewRouter(dbClient)
	http.HandleFunc("/", s.ServeHTTP)
	address := fmt.Sprintf("%s:%s", *webhost, *webport)
	log.Printf("[register] starting server on address: %s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
