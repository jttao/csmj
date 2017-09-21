package api

import (
	"net/http"

	userservice "game/user/service"
	
	gamepkghttputils "game/pkg/httputils"

	taskservice "game/hall/tasks"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleTaskFinish(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求微信分享完成")

	ss := taskservice.TaskServiceInContext(req.Context())
	playerId := userservice.UserInContext(req.Context())
	
	log.WithFields(
		log.Fields{
			"userId": playerId,
		}).Debug("请求微信分享完成")

	us, err := ss.FinishUserTask(playerId,1)
			
	if err != nil { 
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(
			log.Fields{
				"userId": playerId,
				"error":  err,
			}).Error("请求微信分享完成,错误")
		return
	} 
	log.Debug("请求微信分享完成,成功")
	rr := gamepkghttputils.NewSuccessResult(us)
	httputils.WriteJSON(rw, http.StatusOK, rr)
}
