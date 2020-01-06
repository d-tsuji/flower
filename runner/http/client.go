package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type client struct {
	url        string
	httpClient *http.Client
}

func New(url string) *client {
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	return &client{
		url:        url,
		httpClient: http.DefaultClient,
	}
}

func (c *client) Post(payload interface{}) error {
	var payloadBytes []byte
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return errors.WithStack(err)
		}
		payloadBytes = b
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%+v: %+v", resp.Status, string(b)))
	}

	log.Printf("POST response: %v", string(b))
	err = resp.Body.Close()
	return err
}

func (c *client) Get(values url.Values) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if values != nil {
		q := req.URL.Query()
		for k, v := range values {
			for _, vs := range v {
				q.Add(k, vs)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("HTTP Status: %v %v", resp.StatusCode, resp.Status))
	}

	return body, nil
}
