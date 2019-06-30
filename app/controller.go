package app

import (
	"log"

	"github.com/d-tsuji/flower/queue"
	"github.com/d-tsuji/flower/repository"
)

func ConsumeQueue(q *queue.DQueue) (*queue.Item, error) {

	item, err := q.Pop()
	if err != nil {
		log.Fatal(err.Error())
		return &queue.Item{}, err
	}

	log.Println(item)

	return item, nil

}

func InvokeTaskRegister(q *queue.DQueue) {

	item, err := ConsumeQueue(q)
	if err != nil {
		log.Fatal(err.Error())
	}

	if item.TaskType == "Register" {
		_, err := repository.InsertTaskDifinision(item)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func RunRestTask() {

}
