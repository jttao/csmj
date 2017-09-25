package api

import (
	"net/http"

	"game/hall/tasks"
	gamepkghttputils "game/pkg/httputils"
	
	usermodel "game/user/model"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleTaskShareReward(rw http.ResponseWriter, req *http.Request) {

	ss := tasks.TaskServiceInContext(req.Context())

	playerId := userservice.UserInContext(req.Context())
	
	reward,ut,err := ss.RewardUserTask(playerId,1)
	
	if err != nil {
		log.Error("get user task share  error ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} 
	
	if reward { 
		us := userservice.UserServiceInContext(req.Context())   
		reason := usermodel.ReasonTypeTask1 
		err := us.ChangeCardNum( playerId , int64(ut.Reward), reason ) 	
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"userId": playerId,
				"Reward":   ut.Reward,
				"error":   err,
			}).Error("请求每日任务完成,发送奖励失败")
			return
		}
	}

	rr := gamepkghttputils.NewSuccessResult(ut)

	httputils.WriteJSON(rw, http.StatusOK, rr)
}
