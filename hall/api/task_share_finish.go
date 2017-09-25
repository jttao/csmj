package api

import (
	"net/http"

	usermodel "game/user/model"
	userservice "game/user/service"
	
	gamepkghttputils "game/pkg/httputils"

	"game/hall/tasks"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleTaskShareFinish(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求每日分享任务完成")
	
	ss := tasks.TaskServiceInContext(req.Context())
	playerId := userservice.UserInContext(req.Context())
	
	log.WithFields(
		log.Fields{
			"userId": playerId,
		}).Debug("请求每日分享任务完成")
	
	reward,ut,err := ss.FinishUserTask(playerId,1,true)
	
	if err != nil { 
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(
			log.Fields{
				"userId": playerId,
				"error":  err,
			}).Error("请求每日分享任务完成,错误")
		return
	} 
	
	if reward {
		us := userservice.UserServiceInContext(req.Context())  
		err := us.ChangeCardNum( playerId , int64(ut.Reward), usermodel.ReasonTypeTask1) 	
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"userId": playerId,
				"Reward":   ut.Reward,
				"error":   err,
			}).Error("请求每日分享任务完成,发送奖励失败")
			return
		}
	}
	
	log.Debug("请求每日分享任务完成,成功")
	rr := gamepkghttputils.NewSuccessResult(ut)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	
}
