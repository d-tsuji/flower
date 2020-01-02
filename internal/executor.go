package internal

import "fmt"

type executor struct{}

func NewExecutor() *executor {
	return &executor{}
}

func (e *executor) Test1() error {
	fmt.Println("echo Test1")
	return nil
}

func (e *executor) Test2() error {
	fmt.Println("echo Test2")
	return nil
}

func (e *executor) Test3() error {
	fmt.Println("echo Test3")
	return nil
}
