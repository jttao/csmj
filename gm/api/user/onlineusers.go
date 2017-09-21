package user

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type OnlineUsersResponse struct {
	Num int `json:"num"`
}

//获取总用户
func handleOnlineUsers(rw http.ResponseWriter, req *http.Request) {
	log.WithFields(
		log.Fields{
			"ip": req.RemoteAddr,
		}).Debug("请求获取在线用户数")
	us := gmservice.UserServiceInContext(req.Context())
	num, err := us.GetOnlinesUsers()
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("获取在线用户数,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &OnlineUsersResponse{}
	result.Num = num

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":  req.RemoteAddr,
			"num": num,
		}).Debug("请求获取在线用户数,成功")
}
