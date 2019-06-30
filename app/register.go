package app

import (
	"log"

	"github.com/d-tsuji/flower/queue"
	"github.com/d-tsuji/flower/repository"
)

func InvokeExecTaskRegistry(q *queue.DQueue) {

	item, err := ConsumeQueue(q)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = repository.InsertTaskDifinision(item)
	if err != nil {
		log.Fatal(err.Error())
	}
}
