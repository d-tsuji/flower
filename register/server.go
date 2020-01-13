// Package register is a package related to workflow registration
package register

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/d-tsuji/flower/repository"
)

// Router contains settings for connecting to DB.
type Router struct {
	db *repository.DB
}

// NewRouter creates a new Server.
func NewRouter(db *repository.DB) *Router {
	return &Router{
		db: db,
	}
}

// ServeHTTP handles path routing.
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
func (s *Router) register(w http.ResponseWriter, r *http.Request) {
	taskId := strings.TrimPrefix(r.URL.Path, "/register/")
	if strings.HasSuffix(taskId, "/") {
		taskId = taskId[:len(taskId)-1]
	}

	ctx := context.Background()
	ok, err := s.db.InsertExecutableTasks(ctx, taskId)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(fmt.Sprintf("\"{status: failed, taskId: %s, description: %s}\"", taskId, "no tasks registered")))
		log.Printf("[register] no tasks registered. taskId: %s\n", taskId)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(fmt.Sprintf("\"{status: succeeded, taskId: %s}\"", taskId)))
	log.Printf("[register] task registered. taskId: %s\n", taskId)
}
