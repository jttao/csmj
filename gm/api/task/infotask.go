package task

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	model "game/hall/model"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type TaskInfoForm struct {
	TaskId int64 `form:"taskId"`
}

type TaskInfoResponse struct {
	Id            int32  `json:"id"`
	Reward          int32 `json:"reward"`
	TargetNum       int32  `json:"targetNum"`
	Content         string    `json:"content"` 
	CreateTime    int64  `json:"createTime"`
}


func convertTaskToResponse(u *model.TaskModel) *TaskInfoResponse {
	uir := &TaskInfoResponse{}
	uir.Id = u.Id
	uir.Reward = u.Reward
	uir.TargetNum = u.TargetNum
	uir.Content = u.Content 
	uir.CreateTime = u.CreateTime
	return uir
}

func handleInfoTask(rw http.ResponseWriter, req *http.Request) {
	form := &TaskInfoForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("请求任务信息,解析失败")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"form": form,
		}).Debug("请求任务信息")

	us := gmservice.TaskServiceInContext(req.Context())
	taskId := form.TaskId

	task, err := us.GetTaskById(taskId)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"taskId": taskId,
				"error":  err,
			}).Error("请求任务信息,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if task == nil {
		result := &struct{}{}
		rr := gamepkghttputils.NewSuccessResult(result)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"taskId": taskId,
			}).Debug("请求任务信息,查无此人")
	} else {
		result := convertTaskToResponse(user)
		rr := gamepkghttputils.NewSuccessResult(result)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"taskId": taskId,
			}).Debug("请求任务信息,成功")
	}
}
