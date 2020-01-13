package component

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mux *http.ServeMux

func TestComponent_HTTPPostRequest_Normal(t *testing.T) {
	ts := httptest.NewServer(mux)
	defer ts.Close()

	c := NewComponent(map[string]string{
		"URL":  ts.URL + "/test1",
		"BODY": "{id: test1}",
	})
	err := c.HTTPPostRequest()
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestComponent_HTTPPostRequest_AbNormal(t *testing.T) {
	ts := httptest.NewServer(mux)
	defer ts.Close()

	c := NewComponent(map[string]string{
		"URL":  ts.URL + "/xxx",
		"BODY": "{id: test1}",
	})
	err := c.HTTPPostRequest()
	if err == nil {
		t.Fatalf("error wrong url request, but no error")
	}
}

func TestMain(m *testing.M) {
	mux = http.NewServeMux()
	mux.HandleFunc(
		"/test1", func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)
			defer func() {
				if err := r.Body.Close(); err != nil {
					fmt.Printf("%+v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}()
			if err != nil {
				fmt.Printf("%+v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
		},
	)
	status := m.Run()
	os.Exit(status)
}
