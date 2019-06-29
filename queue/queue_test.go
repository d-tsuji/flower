package queue

import (
	"fmt"
	"testing"
)

func TestEnqueue(t *testing.T) {

	var item = Item{
		"sample.task",
		"Register",
		"Normal",
	}

	q, err := NewDQueue()
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

	err = q.Push(&item)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(q.Dque.Size())
}

func TestDequeue(t *testing.T) {

	q, err := NewDQueue()
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

	err = q.Push(&Item{
		"sample.task",
		"Register",
		"Normal",
	})

	message, err := q.Pop()
	if err != nil {
		t.Errorf("%s", err)
	}
	if message == nil {
		t.Log("No message.")
	}

	err = q.Push(&Item{
		"sample.task",
		"Register",
		"Normal",
	})
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(message)
	t.Log(q.Dque.Size())

	err = q.Push(&Item{
		"sample.task",
		"Register",
		"Normal",
	})
	if err != nil {
		t.Errorf("%s", err)
	}
	t.Log(message)
	t.Log(q.Dque.Size())
}
