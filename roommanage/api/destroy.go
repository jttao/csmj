package api

import (
	"fmt"
	"net/http"

	gamehttputils "game/pkg/httputils"
	roommanageservice "game/roommanage/service"
	usermodel "game/user/model"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type DestroyRoomForm struct {
	RoomId int64 `form:"roomId"`
	Refund bool  `form:"refund"`
}

func handleDestory(rw http.ResponseWriter, req *http.Request) {
	form := &DestroyRoomForm{}
	if err := httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Warn("create room request bind failed")
		return
	}
	log.WithField("form", fmt.Sprintf("%#v", form)).Debug("请求摧毁房间")
	rms := roommanageservice.RoomManageInContext(req.Context())
	roomId := form.RoomId
	refund := form.Refund
	r := rms.GetRoomById(roomId)
	if r == nil {
		log.WithField("roomId", roomId).Warn("房间不存在")
		rw.WriteHeader(http.StatusOK)
		return
	}
	us := userservice.UserServiceInContext(req.Context())
	if !rms.Debug() {
		//扣钱
		if refund && r.Cost() != 0 {
			err := us.ChangeCardNum(r.OwnerId(), int64(r.Cost()), usermodel.ReasonTypeRefund)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				log.WithFields(log.Fields{
					"ownerId": r.OwnerId(),
					"cost":    r.Cost,
					"error":   err,
				}).Error("摧毁房间,返款失败")
			}
		}
	}

	for _, p := range r.Players() {
		err := us.Offline(p.Id())
		if err != nil {
			log.WithFields(log.Fields{
				"playerId": p.Id(),
				"error":    err,
			}).Warn("离开房间,离线失败")
		}
	}

	rms.DestroyRoom(roomId)
	result := &struct{}{}
	rr := gamehttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"roomId": roomId,
		}).Debug("请求摧毁房间成功")
}
