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

type CreateRoomForm struct {
	RoomType   int    `form:"roomType"`
	OwnerId    int64  `form:"ownerId"`
	MaxPlayers int    `form:"maxPlayers"`
	Round      int    `form:"round"`
	Cost       int    `form:"cost"`
	RoomConfig string `form:"roomConfig"`
	Location   string `form:"location"`
	ForbidIp   int    `form:"forbidIp"`
	OpenRoomType int  `form:"openRoomType"`
}

type createRoomResponse struct {
	RoomId int64  `json:"roomId"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	OpenRoomType int `json:"openRoomType"`
}

func handleCreate(rw http.ResponseWriter, req *http.Request) {
	var err error
	form := &CreateRoomForm{}
	if err = httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Warn("创建房间解析失败")
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
	location := form.Location
	forbidIp := form.ForbidIp
	ip := req.RemoteAddr
	openRoomType := roommanageservice.OpenRoomType(form.OpenRoomType)

	rms := roommanageservice.RoomManageInContext(req.Context())
	p := rms.GetPlayerById(ownerId)

	
	//非代理开房，玩家在房间内了
	if p != nil {
		rr := gamehttputils.NewFailedResult(PlayerAlreadyInRoomErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithField("ownerId", ownerId).Warn("玩家已经在房间内")
		return
	}

	us := userservice.UserServiceInContext(req.Context())

	forbid, err := us.IsForbid(ownerId)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"location" : form.Location,
			"ip": ip,
			"error":      err,
		}).Error("创建房间,判断是否被封号")
		return
	}

	if forbid {
		rr := gamehttputils.NewFailedResult(AccountIsLockErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		return
	}

	free, err := rms.IfFree()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"location" : form.Location,
			"ip": ip,
			"error":      err,
		}).Error("创建房间,判断是否免房卡")
		return
	}
	
	if !rms.Debug() {
		if !free {
			//判断是否有钱
			cardNum, err := us.GetCardNum(ownerId)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				log.WithFields(log.Fields{
					"ownerId":    form.OwnerId,
					"maxPlayers": form.MaxPlayers,
					"roomType":   form.RoomType,
					"round":      form.Round,
					"cost":       form.Cost,
					"roomConfig": form.RoomConfig,
					"location" : form.Location,
					"ip": ip,
					"error":      err,
				}).Error("创建房间,获取房卡数失败")
				return
			}

			//不够钱
			if cardNum < int64(cost) {
				rr := gamehttputils.NewFailedResult(RoomCardNoEnoughErrorCode)
				httputils.WriteJSON(rw, http.StatusOK, rr)
				return
			}
		}
		
		//跟新地址
		err = us.ChangeUserLocation(ownerId, location)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"ownerId":    form.OwnerId,
				"maxPlayers": form.MaxPlayers,
				"roomType":   form.RoomType,
				"round":      form.Round,
				"cost":       form.Cost,
				"roomConfig": form.RoomConfig,
				"location" : form.Location,
				"ip": ip,
				"error":      err,
			}).Error("创建房间,更细位置失败")
		}  
		
	}

	if free {
		cost = 0
	}
	
	r, err := rms.CreateRoom(roomType, ownerId, maxPlayers, round, cost, roomConfig,forbidIp,ip,openRoomType)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"location" : form.Location,
			"ip": ip,
			"error":      err,
		}).Error("创建房间失败")
		return
	}
	
	//服务器已满
	if r == nil {
		rr := gamehttputils.NewFailedResult(ServerIsFullErrorCode)
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"location" : form.Location,
			"ip": ip,
			"error":      err,
		}).Warn("创建房间失败,服务器已满")
		httputils.WriteJSON(rw, http.StatusOK, rr)
		return
	}
	
	log.WithFields(log.Fields{
		"cost":    cost, 
	}).Debug("创建房间开始扣款")
	
	if !rms.Debug() {
		//扣钱 
		err = us.ChangeCardNum(ownerId, int64(-cost), usermodel.ReasonTypeCost)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"ownerId":    form.OwnerId,
				"maxPlayers": form.MaxPlayers,
				"roomType":   form.RoomType,
				"round":      form.Round,
				"cost":       form.Cost,
				"roomConfig": form.RoomConfig,
				"location" : form.Location,
				"ip": ip,
				"error":      err,
			}).Error("创建房间,扣款失败")
		}  
	}
	
	sc := rms.GetServerByServerId(r.ServerId())
	result := &createRoomResponse{}
	result.RoomId = r.Id()
	result.Host = sc.Host
	result.Port = sc.Port
	result.OpenRoomType = int(r.GetOpenRoomType())
	rr := gamehttputils.NewSuccessResult(result)
	
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"room":       fmt.Sprintf("%#v", r),
		}).Debug("请求创建房间成功")
	err = us.Online(ownerId)
	if err != nil {
		log.WithFields(log.Fields{
			"ownerId":    form.OwnerId,
			"maxPlayers": form.MaxPlayers,
			"roomType":   form.RoomType,
			"round":      form.Round,
			"cost":       form.Cost,
			"roomConfig": form.RoomConfig,
			"location" :  form.Location,
			"ip": ip,
			"error":      err,
		}).Warn("创建房间,在线失败")
	}
	
}
