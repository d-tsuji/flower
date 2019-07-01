package queue

import (
	"os"

	"github.com/joncrlsn/dque"
	"github.com/pkg/errors"
)

// Item is what we'll be storing in the queue.  It can be any struct
// as long as the fields you want stored are public.
type Item struct {
	JobFlowId string // 実行時に発行される一意なフローID(初回はnil)
	TaskId    string // 実行するタスク群のID
	TaskType  string // タスク登録 or タスク監視実行
}

// ItemBuilder creates a new item and returns a pointer to it.
// This is used when we load a segment of the queue from disk.
func ItemBuilder() interface{} {
	return &Item{}
}

// DQueue implements a Queue using a DQueue as the underlying provider
type DQueue struct {
	Dque *dque.DQue
}

// NewDQueue constrcuts a new Queue with an underlying DQueue provider
func NewDQueue() (*DQueue, error) {
	d, err := dque.NewOrOpen("test-dqueue", os.TempDir(), 50, ItemBuilder)

	if err != nil {
		return nil, errors.Wrap(err, "failed to construct new DQue")
	}

	q := &DQueue{
		Dque: d,
	}

	return q, nil
}

// Push item to the end of the queue
func (q *DQueue) Push(o *Item) error {
	return q.Dque.Enqueue(o)
}

// Pop item from top of the queue
func (q *DQueue) Pop() (*Item, error) {
	itf, err := q.Dque.Dequeue()
	if err != nil {
		return nil, err
	}

	item, ok := itf.(*Item)
	if !ok {
		return nil, errors.New("Dequeued object is not an Item pointer")
	}
	return item, nil
}