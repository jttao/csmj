package roomcard

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	roomcardPath = "/roomcard"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(roomcardPath).Subrouter()
	sr.Path("/change").HandlerFunc(http.HandlerFunc(handleChange))
	sr.Path("/today_use").HandlerFunc(http.HandlerFunc(handleTodayUse))
	sr.Path("/free").HandlerFunc(http.HandlerFunc(handleFree))
	sr.Path("/if_free").HandlerFunc(http.HandlerFunc(handleIfFree))
	sr.Path("/check").HandlerFunc(http.HandlerFunc(handleCheck))
	sr.Path("/if_check").HandlerFunc(http.HandlerFunc(handleIfCheck))

}
