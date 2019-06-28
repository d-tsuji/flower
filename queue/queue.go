package queue

import (
	"os"

	"github.com/joncrlsn/dque"
	"github.com/pkg/errors"
)

// Item is what we'll be storing in the queue.  It can be any struct
// as long as the fields you want stored are public.
type Item struct {
	JobId          string
	TaskId         string
	TaskType       string
	DistributeMode string
	ExecNumber     int64
}

// ItemBuilder creates a new item and returns a pointer to it.
// This is used when we load a segment of the queue from disk.
func ItemBuilder() interface{} {
	return &Item{}
}

// DQueue implements a Queue using a DQueue as the underlying provider
type DQueue struct {
	dque *dque.DQue
}

// NewDQueue constrcuts a new Queue with an underlying DQueue provider
func NewDQueue() (*DQueue, error) {
	d, err := dque.NewOrOpen("test-dqueue", os.TempDir(), 50, ItemBuilder)

	if err != nil {
		return nil, errors.Wrap(err, "failed to construct new DQue")
	}

	q := &DQueue{
		dque: d,
	}

	return q, nil
}

// Push item to the end of the queue
func (q *DQueue) Push(o interface{}) error {
	return q.dque.Enqueue(o)
}

// Pop item from top of the queue
func (q *DQueue) Pop() (interface{}, error) {
	return q.dque.Dequeue()
}
