package api

import (
	"net/http"

	"game/hall/tasks"
	gamepkghttputils "game/pkg/httputils"
	
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleUserTaskList(rw http.ResponseWriter, req *http.Request) {

	ns := tasks.TaskServiceInContext(req.Context())

	userId := userservice.UserInContext(req.Context())

	result,err := ns.GetUserTasks(userId)
	
	if err != nil {
		log.Error("get user tasks  error ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
}
