package app

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/d-tsuji/flower/repository"
)

func RestCall(r *repository.RestTask) ([]byte, error) {
	logger, _ := zap.NewDevelopment()

	client := &http.Client{}
	request, err := http.NewRequest(r.Method, r.Endpoint, strings.NewReader(r.ExtendParameter))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	res, err := client.Do(request)
	if err != nil {
		logger.Warn("An unexpected error has occurred", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Warn("An unexpected error has occurred", zap.Int("Status Code", res.StatusCode))
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Warn("An unexpected error has occurred", zap.Error(err))
		return nil, err
	}
	return body, nil

}
