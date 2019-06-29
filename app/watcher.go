package app

import (
	"log"

	"github.com/d-tsuji/flower/repository"
)

func WatchTask(ch chan<- repository.KrTaskStatus) error {

	log.Println("Task Watching...")
	list, err := repository.SelectExecTarget()
	if err != nil {
		log.Fatal("%s", err)
		return err
	}

	for _, v := range *list {
		log.Printf("Executable task found. ", v)
		ch <- v
	}
	return nil
}
