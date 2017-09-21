package api

import (
	"net/http"

	userservice "game/user/service"
	
	gamepkghttputils "game/pkg/httputils"

	taskservice "game/hall/tasks"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleTaskGet(rw http.ResponseWriter, req *http.Request) {
	
	log.Debug("请求获取任务列表信息")
	
	ss := taskservice.TaskServiceInContext(req.Context())
	playerId := userservice.UserInContext(req.Context())

	log.WithFields(
		log.Fields{
			"userId": playerId,
		}).Debug("请求获取微信分享信息")
	us, err := ss.GetUserTasks(playerId) 
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(
			log.Fields{
				"error": err,
			}).Error("请求获取任务列表,错误")
		return 
	}	
	rr := gamepkghttputils.NewSuccessResult(us)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.Debug("请求获取任务列表信息,成功")
}
