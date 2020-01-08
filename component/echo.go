package component

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func (e *executor) EchoRandomTimeSleep() error {
	randTime := random.Intn(5) + 1

	fmt.Printf("[executor] start EchoRandomTimeSleep. (%v second sleep)\n", randTime)
	time.Sleep(time.Duration(randTime) * time.Second)
	fmt.Printf("[executor] finish EchoRandomTimeSleep\n")

	return nil
}

func (e *executor) EchoParamTimeSleep() error {
	sleepTimeStr, ok := e.params["SLEEP_TIME_SECOND"]
	if !ok {
		return errors.New("EchoParamTimeSleep() required SLEEP_TIME_SECOND parameter")
	}
	sleepTime, err := strconv.Atoi(sleepTimeStr)
	if err != nil {
		return errors.New(fmt.Sprintf("string %s connot convert int", sleepTimeStr))
	}

	fmt.Printf("[executor] start EchoParamTimeSleep. (%v second sleep)\n", sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	fmt.Printf("[executor] finish EchoParamTimeSleep\n")

	return nil
}