package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	changsharoomservice "game/mahjong/changsha/room/service"
	userservice "game/user/service"
	"github.com/xozrc/pkg/httputils"
)

type JoinForm struct {
	RoomId int64 `form:"roomId"`
	Location string `form:"location"`
}

func handleJoin(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求加入长沙房间")
	form := &JoinForm{}

	if err := httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Warn("请求加入长沙房间 解析错误")
		return
	}

	csrs := changsharoomservice.ChangShaRoomServiceInContext(req.Context())

	playerId := userservice.UserInContext(req.Context())
	
	result, err := csrs.JoinRoom(playerId, form.RoomId , form.Location )
	
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Error("请求加入长沙房间 失败")
		return
	}

	httputils.WriteJSON(rw, http.StatusOK, result)
	log.Debug("请求加入长沙房间 成功")
}
