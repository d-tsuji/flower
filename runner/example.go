package runner

import (
	"fmt"
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// Test1 is the sample task.
func (e *executor) Test1() error {
	needTime := random.Intn(10)
	fmt.Printf("echo Test1 (%v second). param1: %s, param2: %s\n", needTime, e.params["hoge"], e.params["piyo"])
	time.Sleep(time.Duration(needTime) * time.Second)
	return nil
}

// Test2 is the sample task.
func (e *executor) Test2() error {
	needTime := random.Intn(10)
	fmt.Printf("echo Test2 (%v second).\n", needTime)
	time.Sleep(time.Duration(needTime) * time.Second)
	return nil
}

// Test3 is the sample task.
func (e *executor) Test3() error {
	needTime := random.Intn(10)
	fmt.Printf("echo Test3 (%v second).\n", needTime)
	time.Sleep(time.Duration(needTime) * time.Second)
	return nil
}
