package main

import (
	"log"
	"net/http"

	"github.com/d-tsuji/flower-v2/db"
	"github.com/d-tsuji/flower-v2/register"
)

func main() {
	db, err := db.New(&db.Opt{Password: "flower"})
	if err != nil {
		log.Printf("postgres initialize error: %v\n", err)
	}
	s := register.NewServer(db)
	http.HandleFunc("/", s.ServeHTTP)
	address := ":8000"
	log.Println("Starting server on address", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
