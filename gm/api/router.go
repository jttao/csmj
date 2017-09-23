package api

import "github.com/gorilla/mux"
import "game/gm/api/roomcard"
import "game/gm/api/user"
import "game/gm/api/task"

const (
	gmPath = "/gm"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(gmPath).Subrouter()
	roomcard.Router(sr)
	user.Router(sr)
	task.Router(sr)
}
