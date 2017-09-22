package api

import (
	"fmt"
	"net/http"

	gamepkghttputils "game/pkg/httputils"
	roommanageclient "game/roommanage/client"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type LoginResponse struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Ip     string `json:"ip"`
	RoomId int64  `json:"roomId"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Sex    int    `json:"sex"`
}

func handleLogin(rw http.ResponseWriter, req *http.Request) {
	
	playerId := userservice.UserInContext(req.Context())
	log.WithFields(
		log.Fields{
			"玩家id": playerId,
		}).Debug("大厅请求登陆")
	rmc := roommanageclient.RoomManageClientInContext(req.Context())
	rId, _, host, port, _, _, _,_,_,_,_,err := rmc.Query(playerId)
	if err != nil {
		log.WithFields(
			log.Fields{
				"玩家id":  playerId,
				"error": err,
			}).Error("大厅请求登陆失败1")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	us := userservice.UserServiceInContext(req.Context())
	user, err := us.GetUserById(playerId)
	if err != nil {
		log.WithFields(
			log.Fields{
				"玩家id":  playerId,
				"error": err,
			}).Error("大厅请求登陆失败2")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.WithFields(
			log.Fields{
				"玩家id":  playerId,
				"error": err,
			}).Warn("大厅请求登陆,用户不存在")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	name := fmt.Sprintf("guest_%d", playerId)
	if user.Name != "" {
		name = user.Name
	}
	lastLoginIp := req.RemoteAddr

	err = us.UpdateUser(playerId, name, lastLoginIp)
	if err != nil {
		log.WithFields(
			log.Fields{
				"玩家id":  playerId,
				"error": err,
			}).Error("大厅请求登陆,更新基本信息")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &LoginResponse{
		Id:     playerId,
		Ip:     lastLoginIp,
		Name:   name,
		RoomId: rId,
		Host:   host,
		Port:   port,
		Image:  user.Image,
		Sex:    user.Sex,
	}
	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"玩家id": playerId,
		}).Debug("大厅请求登陆成功")
}
