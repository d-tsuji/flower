package runner

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/d-tsuji/flower/runner/http"

	"github.com/pkg/errors"
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

func (e *executor) TestHTTPPostRequest() error {
	url, ok := e.params["URL"]
	if !ok {
		return errors.New("executor param does not contain URL err.")
	}
	body, _ := e.params["BODY"]
	var dummy map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dummy); err != nil {
		return errors.New("params BODY is not json format.")
	}

	client := http.New(url)
	if err := client.Post(body); err != nil {
		return errors.New(fmt.Sprintf("post request error. url: %s.", url))
	}
	return nil
}
