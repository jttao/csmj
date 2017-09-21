package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	hallPath = "/hall"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(hallPath).Subrouter()
	sr.Path("/login").HandlerFunc(http.HandlerFunc(handleLogin))
	sr.Path("/news").HandlerFunc(http.HandlerFunc(handleNews))
	sr.Path("/notices").HandlerFunc(http.HandlerFunc(handleNotices))
	sr.Path("/cards").HandlerFunc(http.HandlerFunc(handleCards))
	
	sr.Path("/task_get").HandlerFunc(http.HandlerFunc(handleTaskGet))
	sr.Path("/task_reward").HandlerFunc(http.HandlerFunc(handleTaskReward))
	sr.Path("/task_finish").HandlerFunc(http.HandlerFunc(handleTaskFinish))
	
}
