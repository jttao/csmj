package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	roomManagePath = "/roommanage"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(roomManagePath).Subrouter()
	sr.Path("/create").HandlerFunc(http.HandlerFunc(handleCreate))
	sr.Path("/join").HandlerFunc(http.HandlerFunc(handleJoin))
	sr.Path("/leave").HandlerFunc(http.HandlerFunc(handleLeave))
	sr.Path("/destroy").HandlerFunc(http.HandlerFunc(handleDestory))
	sr.Path("/list").HandlerFunc(http.HandlerFunc(handleList))
	sr.Path("/query").HandlerFunc(http.HandlerFunc(handleQuery))
	sr.Path("/auto").HandlerFunc(http.HandlerFunc(handleAuto))
}
