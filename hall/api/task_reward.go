package api

import (
	"net/http"

	userservice "game/user/service"
	
	gamepkghttputils "game/pkg/httputils"

	taskservice "game/hall/tasks"
	
	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleTaskReward(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求微信分享领取奖励")

	ss := taskservice.TaskServiceInContext(req.Context())
	playerId := userservice.UserInContext(req.Context())

	log.WithFields(
		log.Fields{
			"userId": playerId,
		}).Debug("请求微信分享领取奖励")
			
	us, err := ss.GetReward(playerId,1)
	
	if err != nil {
		log.WithFields(
			log.Fields{
				"error": err,
			}).Error("请求微信分享领取奖励 错误")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	rr := gamepkghttputils.NewSuccessResult(us)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.Debug("请求微信分享领取奖励,成功")
}
