package mock

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/d-tsuji/flower/repository"
	"go.uber.org/zap"
)

// リクエストに応じてタスクをDBに登録する
func register(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()

	qs := r.URL.Query()

	item := &repository.Job{
		JobFlowId: "",
		TaskId:    qs.Get("taskId"),
		TaskType:  "Normal",
	}
	res, err := repository.InsertTaskDefinition(item)
	if err != nil {
		logger.Error("Error to register task to database.", zap.Error(err))
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
	logger, _ := zap.NewDevelopment()
	logger.Info("Request received!")
	time.Sleep(2 * time.Second)
	//fmt.Fprintln(w, "Hello world!")
}

// テスト用
func heavy(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()
	time.Sleep(10 * time.Second)
	logger.Info("Heavy Process start.")
	//fmt.Fprintln(w, "Heavy Process finish.")
}

// テスト用
func param(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()
	time.Sleep(2 * time.Second)
	method := r.Method
	//logger.Info("[Method] %v\n", method)
	//logger.Info("[Method] %v\n", zap. )
	//for k, v := range r.Header {
	//	logger.Info("[Header] %v: %s\n", k, zap.String(strings.Join(v, ",")))
	//}

	if method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error("Error request", zap.Error(err))
		}

		decoded, error := url.QueryUnescape(string(body))
		if error != nil {
			logger.Error("Error decode body", zap.Error(error))
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
