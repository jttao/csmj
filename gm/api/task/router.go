package task

import "github.com/gorilla/mux"
import "net/http"

const (
	taskPath = "/task"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(taskPath).Subrouter()
	sr.Path("/list").HandlerFunc(http.HandlerFunc(handleListTasks))
	sr.Path("/info").HandlerFunc(http.HandlerFunc(handleInfoTask))
	sr.Path("/save").HandlerFunc(http.HandlerFunc(handleSaveTask)) 
}
