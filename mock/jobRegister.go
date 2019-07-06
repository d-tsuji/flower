package mock

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/d-tsuji/flower/repository"
)

// リクエストに応じてタスクをDBに登録する
func register(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	item := &repository.Item{
		"",
		qs.Get("taskId"),
		"Normal",
	}
	res, err := repository.InsertTaskDefinition(item)
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
	fmt.Fprintf(w, "%v : Complete registered task. TaskId -> %s.", time.Now(), qs.Get("taskId"))
}

// テスト用
func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received!")
	time.Sleep(2 * time.Second)
	fmt.Fprintln(w, "Hello world!")
}

// テスト用
func heavy(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	log.Println("Heavy Process start.")
	fmt.Fprintln(w, "Heavy Process finish.")
}

// テスト用
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
	http.HandleFunc("/register", register)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/heavy", heavy)
	http.HandleFunc("/param", param)
	http.ListenAndServe(":8021", nil)
}
