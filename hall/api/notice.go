package api

import (
	"net/http"

	notice "game/hall/notice"
	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleNotices(rw http.ResponseWriter, req *http.Request) {

	ns := notice.NoticeServiceInContext(req.Context())
	notices, err := ns.GetNotices()
	if err != nil {
		log.Error("get notices  error ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rr := gamepkghttputils.NewSuccessResult(notices)
	httputils.WriteJSON(rw, http.StatusOK, rr)
}
