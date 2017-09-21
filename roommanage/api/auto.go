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

func handleAuto(rw http.ResponseWriter, req *http.Request) {
	var err error
	form := &CreateRoomForm{}
	if err = httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Warn("create room request bind failed")
		return
	}

	log.WithField("form", fmt.Sprintf("%#v", form)).Debug("请求创建房间")
	roomType := roommanageservice.RoomType(form.RoomType)
	if !roomType.Valid() {
		log.WithField("roomType", form.RoomType).Warn("房间类型不存在")
		return
	}

	ownerId := form.OwnerId
	maxPlayers := form.MaxPlayers
	round := form.Round
	cost := form.Cost
	roomConfig := form.RoomConfig
	forbidIp := form.ForbidIp
	location := form.Location
	ip := req.RemoteAddr
	openRoomType := roommanageservice.OpenRoomType(form.OpenRoomType)
	
	rms := roommanageservice.RoomManageInContext(req.Context())
	p := rms.GetPlayerById(ownerId)

	//玩家在房间内了
	if p != nil {
		rr := gamehttputils.NewFailedResult(PlayerAlreadyInRoomErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithField("ownerId", ownerId).Warn("玩家已经在房间内")
		return
	}

	check, err := rms.IfCheck()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"error":      err,
		}).Error("自动加入房间失败")
		return
	}
	//非审核不可以自动加入
	if !check {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
		}).Warn("自动加入房间,但是房间不是审核状态")
		return
	}
	
	cost = 0
	r, err := rms.AutoRoom(roomType, ownerId, maxPlayers, round, cost, roomConfig , forbidIp,ip,openRoomType)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"error":      err,
		}).Error("自动加入房间失败")
		return
	}

	if r == nil {
		rr := gamehttputils.NewFailedResult(ServerIsFullErrorCode)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"error":      err,
		}).Warn("创建房间失败,服务器已满")
		httputils.WriteJSON(rw, http.StatusOK, rr)
		return
	}


	us := userservice.UserServiceInContext(req.Context())

	err = us.ChangeUserLocation(ownerId,location)
	if err != nil {
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"error":      err,
		}).Warn("创建房间,更新位置失败")
	}
	
	sc := rms.GetServerByServerId(r.ServerId())
	result := &createRoomResponse{}
	result.RoomId = r.Id()
	result.Host = sc.Host
	result.Port = sc.Port
	rr := gamehttputils.NewSuccessResult(result)

	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"room":       fmt.Sprintf("%#v", r),
		}).Debug("请求自动加入房间成功")
	
	
	err = us.Online(ownerId)
	if err != nil {
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"error":      err,
		}).Warn("创建房间,在线失败")
	}
}
