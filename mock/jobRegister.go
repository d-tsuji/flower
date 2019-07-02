package mock

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/d-tsuji/flower/queue"
	"github.com/d-tsuji/flower/repository"
)

func handler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	item := &queue.Item{
		"",
		qs.Get("taskId"),
		"Normal",
	}
	res, err := repository.InsertTaskDifinision(item)
	if err != nil {
		log.Fatal(err)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if cnt == 0 {
		log.Printf("taskId :[%s] is not registered in ms_task. Please Check configuration.", item.TaskId)
		fmt.Fprintf(w, "taskId :[%s] is not registered in ms_task. Please Check configuration.", item.TaskId)
		return
	}
	fmt.Fprintf(w, "Complete registered task. TaskId -> %s.", qs.Get("taskId"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received!")
	time.Sleep(2 * time.Second)
	fmt.Fprintln(w, "Hello world!")
}

func heavy(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	log.Println("Heavy Process start.")
	fmt.Fprintln(w, "Heavy Process finish.")
}

func param(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	method := r.Method
	log.Printf("[Method] %v\n", method)
	for k, v := range r.Header {
		log.Printf("[Header] %v: %s\n", k, strings.Join(v, ","))
	}

	if method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		decoded, error := url.QueryUnescape(string(body))
		if error != nil {
			log.Fatal(error)
		}
		fmt.Fprintln(w, decoded)
	}
}

func RegisterTask() {
	http.HandleFunc("/register", handler)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/heavy", heavy)
	http.HandleFunc("/param", param)
	http.ListenAndServe(":8021", nil)
}
