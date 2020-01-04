package runner

import (
	"fmt"
	"math/rand"
	"time"
)

type executor struct{}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewExecutor() *executor {
	return &executor{}
}

func (e *executor) Test1() error {
	fmt.Println("echo Test1")
	time.Sleep(time.Duration(random.Intn(10)) * time.Second)
	return nil
}

func (e *executor) Test2() error {
	fmt.Println("echo Test2")
	time.Sleep(time.Duration(random.Intn(10)) * time.Second)
	return nil
}

func (e *executor) Test3() error {
	fmt.Println("echo Test3")
	time.Sleep(time.Duration(random.Intn(10)) * time.Second)
	return nil
}
