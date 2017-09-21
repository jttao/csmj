package user

import "github.com/gorilla/mux"
import "net/http"

const (
	userPath = "/user"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(userPath).Subrouter()
	sr.Path("/new").HandlerFunc(http.HandlerFunc(handleNewUsers))
	sr.Path("/total").HandlerFunc(http.HandlerFunc(handleTotalUsers))
	sr.Path("/online").HandlerFunc(http.HandlerFunc(handleOnlineUsers))
	sr.Path("/forbid").HandlerFunc(http.HandlerFunc(handleForbid))
	sr.Path("/info").HandlerFunc(http.HandlerFunc(handleInfo))
	sr.Path("/list").HandlerFunc(http.HandlerFunc(handleList))
	sr.Path("/info_by_wx").HandlerFunc(http.HandlerFunc(handleInfoByWx))

}
