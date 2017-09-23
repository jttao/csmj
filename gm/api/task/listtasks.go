package task

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"
	
	model "game/hall/model"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

const (
	defaultPageSize    = 15
	maxDefaultPageSize = 100
)

type TaskListForm struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type TaskListResponse struct {
	TotalPage int         `json:"totalPage"`
	Page      int         `json:"page"`
	PageSize  int         `json:"pageSize"`
	Data      interface{} `json:"data"`
}


//获取玩家列表
func handleListTasks(rw http.ResponseWriter, req *http.Request) {

	form := &TaskListForm{}
	
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("请求任务列表信息,解析失败")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	page := form.Page
	pageSize := form.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	if pageSize > maxDefaultPageSize {
		pageSize = maxDefaultPageSize
	}

	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"form": form,
		}).Debug("请求用户列表信息")

	us := gmservice.TaskServiceInContext(req.Context())

	totalPage, tasks, err := us.GetTasks(page, pageSize)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("请求用户列表信息,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := &TaskListResponse{}
	taskResults := make([]*TaskInfoResponse, 0, len(tasks))
	for _, user := range tasks {
		taskResult := convertTaskToResponse(user)
		taskResults = append(taskResults, taskResult)
	}

	result.Data = taskResults
	result.Page = page
	result.PageSize = pageSize
	result.TotalPage = totalPage

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip": req.RemoteAddr,
		}).Debug("请求任务列表,成功") 
	
}
