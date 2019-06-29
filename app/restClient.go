package app

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/d-tsuji/flower/repository"
)

func RestCall(r *repository.RestTask) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest(r.Method, r.Endpoint, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal(res)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return body, nil

}
