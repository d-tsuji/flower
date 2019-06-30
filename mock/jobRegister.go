package mock

import (
	"fmt"
	"log"
	"net/http"

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
	err := repository.InsertTaskDifinision(item)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Complete Task registration. Taskid -> %s.", qs.Get("taskId"))
}

func StartServer() {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(":8021", nil)
}
