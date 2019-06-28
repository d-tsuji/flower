package queue

import (
	"fmt"
	"testing"
)

func TestEnqueue(t *testing.T) {

	var item = Item{
		"",
		"sample.task",
		"Register",
		"Normal",
		1,
	}

	q, err := NewDQueue()
	if err != nil {
		t.Errorf("%s", err)
	}
	defer func() {
		for q.dque.Size() > 0 {
			_, err := q.Pop()
			if err != nil {
				t.Errorf("%s", err)
			}
		}
	}()

	err = q.Push(item)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(q.dque.Size())
}

func TestDequeue(t *testing.T) {

	q, err := NewDQueue()
	if err != nil {
		t.Errorf("%s", err)
	}
	defer func() {
		for q.dque.Size() > 0 {
			_, err := q.Pop()
			if err != nil {
				t.Errorf("%s", err)
			}
		}
	}()

	err = q.Push(Item{
		"",
		"sample.task",
		"Register",
		"Normal",
		1,
	})

	message, err := q.Pop()
	if err != nil {
		t.Errorf("%s", err)
	}
	if message == nil {
		t.Log("No message.")
	}

	err = q.Push(Item{
		"",
		"sample.task",
		"Register",
		"Normal",
		1,
	})
	if err != nil {
		t.Errorf("%s", err)
	}

	t.Log(message)
	t.Log(q.dque.Size())

	err = q.Push(Item{
		"",
		"sample.task",
		"Register",
		"Normal",
		1,
	})
	if err != nil {
		t.Errorf("%s", err)
	}
	t.Log(message)
	t.Log(q.dque.Size())
}
