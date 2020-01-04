package runner

import (
	"fmt"
	"math/rand"
	"time"
)

// Executor is a struct for executing tasks.
type executor struct{}

// NewServer creates a new executor.
func NewExecutor() *executor {
	return &executor{}
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// Test1 is the sample task.
func (e *executor) Test1() error {
	fmt.Println("echo Test1")
	time.Sleep(time.Duration(random.Intn(10)) * time.Second)
	return nil
}

// Test2 is the sample task.
func (e *executor) Test2() error {
	fmt.Println("echo Test2")
	time.Sleep(time.Duration(random.Intn(10)) * time.Second)
	return nil
}

// Test3 is the sample task.
func (e *executor) Test3() error {
	fmt.Println("echo Test3")
	time.Sleep(time.Duration(random.Intn(10)) * time.Second)
	return nil
}
