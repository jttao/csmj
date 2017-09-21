package api

import (
	"net/http"

	userservice "game/user/service"

	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type CardsResponse struct {
	Cards int64 `json:"cards"`
}

func handleCards(rw http.ResponseWriter, req *http.Request) {

	us := userservice.UserServiceInContext(req.Context())
	ownerId := userservice.UserInContext(req.Context())

	cards, err := us.GetCardNum(ownerId)
	if err != nil {
		log.Error("get cards  error ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &CardsResponse{
		Cards: cards,
	}
	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
}
