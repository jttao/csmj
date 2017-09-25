package api

import (
	"net/http"

	usermodel "game/user/model"
	userservice "game/user/service"
	
	gamepkghttputils "game/pkg/httputils"

	taskservice "game/hall/tasks"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type TaskFinishForm struct { 
	PlayerId     	int64	`form:"playerId"`  
	TaskShareId     int32	`form:"taskShareId"` 
	State     		bool    `form:"state"` 
}

func handleTaskFinish(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求每日任务完成")
	taskForm := &TaskFinishForm{}
	if err := httputils.Bind(req, taskForm); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	
	playerId := taskForm.PlayerId
	taskShareId := taskForm.TaskShareId
	state := taskForm.State 

	ss := taskservice.TaskServiceInContext(req.Context()) 
	
	log.WithFields(
		log.Fields{
			"userId": playerId,
		}).Debug("请求每日任务完成")
	
	reward,ut,err := ss.FinishUserTask(playerId,taskShareId,state)
	
	if err != nil { 
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(
			log.Fields{
				"userId": playerId,
				"error":  err,
			}).Error("请求每日任务完成,错误")
		return
	} 
	
	if reward {

		us := userservice.UserServiceInContext(req.Context())  

		reason := usermodel.ReasonType(0)
		if taskShareId==1 {
			reason = usermodel.ReasonTypeTask1
		}
		if taskShareId==2 {
			reason = usermodel.ReasonTypeTask2
		}
		if taskShareId==3 {
			reason = usermodel.ReasonTypeTask3
		}
		
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
	
	log.Debug("请求每日任务完成,成功")
	rr := gamepkghttputils.NewSuccessResult(ut)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	
}
