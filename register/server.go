// Package register is a package related to workflow registration
package register

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/d-tsuji/flower/repository"
)

// Server contains settings for connecting to DB.
type Server struct {
	db *repository.DB
}

// NewServer creates a new Server.
func NewServer(db *repository.DB) *Server {
	return &Server{
		db: db,
	}
}

// ServeHTTP handles path routing.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Validate request
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/register") {
		s.register(w, r)
		return
	}
}

// Register registers a waiting task from taskId.
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	taskId := strings.TrimPrefix(r.URL.Path, "/register/")
	if strings.HasSuffix(taskId, "/") {
		taskId = taskId[:len(taskId)-1]
	}

	ctx := context.Background()
	if err := s.db.InsertExecutableTasks(ctx, taskId); err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(fmt.Sprintf("{status: succeeded, taskId: %s}", taskId))
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
	log.Printf("[register] task registered. taskId: %s\n", taskId)
}
