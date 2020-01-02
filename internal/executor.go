package internal

import "fmt"

type Executor struct{}

func (e *Executor) Test1() error {
	fmt.Println("echo Test1")
	return nil
}

func (e *Executor) Test2() error {
	fmt.Println("echo Test2")
	return nil
}

func (e *Executor) Test3() error {
	fmt.Println("echo Test3")
	return nil
}
