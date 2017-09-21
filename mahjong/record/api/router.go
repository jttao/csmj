package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	recordPath = "/record"
)

func Router(r *mux.Router) {
	recordList := recordPath + "/list"
	r.Path(recordList).HandlerFunc(http.HandlerFunc(handleRecordList))
	recordGet := recordPath + "/get"
	r.Path(recordGet).HandlerFunc(http.HandlerFunc(handleRecordGet))
	roundGet := recordPath + "/round/get"
	r.Path(roundGet).HandlerFunc(http.HandlerFunc(handleRoundGet))
}
