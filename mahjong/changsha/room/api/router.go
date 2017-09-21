package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	roomPath = "/room"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(roomPath).Subrouter()
	sr.Path("/create").HandlerFunc(http.HandlerFunc(handleCreate))
	sr.Path("/join").HandlerFunc(http.HandlerFunc(handleJoin))
	sr.Path("/auto").HandlerFunc(http.HandlerFunc(handleAuto))
}
