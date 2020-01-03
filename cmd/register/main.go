package main

import (
	"log"
	"net/http"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/d-tsuji/flower-v2/register"
)

func main() {
	dbClient, err := db.New(&db.Opt{Password: "flower"})
	defer dbClient.Close()
	if err != nil {
		log.Printf("postgres initialize error: %v\n", err)
	}
	s := register.NewServer(dbClient)
	http.HandleFunc("/", s.ServeHTTP)
	address := "0.0.0.0:8000"
	log.Println("Starting server on address", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
