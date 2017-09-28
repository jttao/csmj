package task

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"
	
	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type SaveTaskForm struct {
	TaskId int32 `form:"taskId"` 
	Reward int32 `form:"reward"` 
	TargetNum int32 `form:"targetNum"` 
	Content string `form:"content"` 
}

func handleSaveTask(rw http.ResponseWriter, req *http.Request) {
	form := &SaveTaskForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("修改任务,解析失败")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	
	ts := gmservice.TaskServiceInContext(req.Context())
	taskId := form.TaskId
	reward := form.Reward
	targetNum := form.TargetNum
	content := form.Content
	
	//判断任务是否存在
	err = ts.UpdateTask(taskId,reward,targetNum,content)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"taskId": taskId,
				"error":  err,
			}).Error("修改任务,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	result := &struct{}{}

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":     req.RemoteAddr,
			"taskId": taskId,
			"reward":   reward,
			"targetNum":   targetNum, 
		}).Debug("修改任务,成功")

}
