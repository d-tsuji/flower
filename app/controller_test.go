package app

import (
	"testing"

	"github.com/d-tsuji/flower/queue"
)

func TestQueueConsume(t *testing.T) {

	q, err := queue.NewDQueue()
	if err != nil {
		t.Errorf("%s", err)
	}
	defer func() {
		for q.Dque.Size() > 0 {
			_, err := q.Pop()
			if err != nil {
				t.Errorf("%s", err)
			}
		}
	}()

	err = q.Push(&queue.Item{
		"sample.task",
		"Register",
		"Normal",
	})

	item, err := ConsumeQueue(q)
	if err != nil {
		t.Errorf("%s", err)
	}
	t.Log(item)
}
