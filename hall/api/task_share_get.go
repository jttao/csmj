package api

import (
	"net/http"

	"game/hall/tasks"
	gamepkghttputils "game/pkg/httputils"
	
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleTaskShareGet(rw http.ResponseWriter, req *http.Request) {

	ss := tasks.TaskServiceInContext(req.Context())

	userId := userservice.UserInContext(req.Context())
	
	result,err := ss.GetUserTask(userId,1)
	
	if err != nil {
		log.Error("get user task share  error ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
}
