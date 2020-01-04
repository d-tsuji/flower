package register

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/d-tsuji/flower-v2/db"
)

// Payload is included in the request body from the client
type Payload struct {
	TaskId     string `json:"taskId"`
	Parameters string `json:"parameters"`
}

// Server contains settings for connecting to DB.
type Server struct {
	db *db.DB
}

// NewServer creates a new Server.
func NewServer(db *db.DB) *Server {
	return &Server{
		db: db,
	}
}

// ServeHTTP handles path routing.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//　Validate request
	if r.Method != "POST" {
		http.Error(w, "Use the appropriate Method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Use the appropriate Content-Type", http.StatusBadRequest)
		return
	}

	if r.URL.Path == "/register" {
		s.register(w, r)
		return
	}
}

// Register registers a waiting task from taskId.
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
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

	var payload Payload
	err = json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	if err := s.db.InsertExecutableTasks(ctx, payload.TaskId); err != nil {
		fmt.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal("{status: succeeded}")
	if err != nil {
		fmt.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
