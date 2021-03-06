package component

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func (c *component) EchoRandomTimeSleep() error {
	randTime := random.Intn(5) + 1

	fmt.Printf("[component] start EchoRandomTimeSleep. (%v second sleep)\n", randTime)
	time.Sleep(time.Duration(randTime) * time.Second)
	fmt.Printf("[component] finish EchoRandomTimeSleep\n")

	return nil
}

func (c *component) EchoParamTimeSleep() error {
	sleepTimeStr, ok := c.params["SLEEP_TIME_SECOND"]
	if !ok {
		return errors.New("EchoParamTimeSleep() required SLEEP_TIME_SECOND parameter")
	}
	sleepTime, err := strconv.Atoi(sleepTimeStr)
	if err != nil {
		return errors.New(fmt.Sprintf("string %s connot convert int", sleepTimeStr))
	}

	fmt.Printf("[component] start EchoParamTimeSleep. (%v second sleep)\n", sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	fmt.Printf("[component] finish EchoParamTimeSleep\n")

	return nil
}
