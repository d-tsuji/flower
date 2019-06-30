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

func RegisterTask() {
	http.HandleFunc("/register", handler)
	http.ListenAndServe(":8021", nil)
}
