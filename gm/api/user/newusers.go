package user

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type NewUsersResponse struct {
	Num int `json:"num"`
}

//获取新创建用户
func handleNewUsers(rw http.ResponseWriter, req *http.Request) {
	log.WithFields(
		log.Fields{
			"ip": req.RemoteAddr,
		}).Debug("请求获取新创建用户数")
	us := gmservice.UserServiceInContext(req.Context())
	num, err := us.GetTodayNewUsers()
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("获取新创建用户数,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &NewUsersResponse{}
	result.Num = num

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":  req.RemoteAddr,
			"num": num,
		}).Debug("请求获取新创建用户数,成功")
}
