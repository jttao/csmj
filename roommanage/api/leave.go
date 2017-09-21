package api

import (
	"fmt"
	"net/http"

	gamehttputils "game/pkg/httputils"
	roommanageservice "game/roommanage/service"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type LeaveRoomForm struct {
	RoomId   int64 `json:"roomId"`
	PlayerId int64 `json:"playerId"`
}

func handleLeave(rw http.ResponseWriter, req *http.Request) {

	form := &LeaveRoomForm{}
	if err := httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithField("req", fmt.Sprintf("%#v", req)).Warn("bind form failed")
		return
	}
	log.WithField("form", fmt.Sprintf("%#v", form)).Debug("请求离开房间")
	rms := roommanageservice.RoomManageInContext(req.Context())

	pId := form.PlayerId
	rId := form.RoomId
	//查找玩家是否在房间内
	p := rms.GetPlayerById(pId)
	if p == nil {
		rw.WriteHeader(http.StatusOK)
		log.WithFields(log.Fields{
			"playerId": pId,
			"roomId":   rId,
		}).Warn("请求离开房间,玩家已经离开了")
		return
	}

	r := rms.GetRoomById(rId)
	//房间不存在了
	if r == nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"playerId": pId,
			"roomId":   rId,
		}).Warn("请求离开房间,房间不存在")
		return
	}

	rms.LeaveRoom(rId, p.Id())

	result := &struct{}{}
	rr := gamehttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	us := userservice.UserServiceInContext(req.Context())

	err := us.Offline(pId)
	if err != nil {
		log.WithFields(log.Fields{
			"playerId": pId,
			"error":    err,
		}).Warn("离开房间,离线失败")
	}
}
