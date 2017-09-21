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

type joinRoomForm struct {
	RoomId   int64 `json:"roomId"`
	PlayerId int64 `json:"playerId"`
	Location string `form:"location"`
}

type joinRoomResponse struct {
	RoomId int64  `json:"roomId"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

func handleJoin(rw http.ResponseWriter, req *http.Request) {

	form := &joinRoomForm{}
	if err := httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithField("req", fmt.Sprintf("%#v", req)).Warn("bind form failed")
		return
	}
	
	log.WithField("form", fmt.Sprintf("%#v", form)).Debug("请求加入房间")
	rms := roommanageservice.RoomManageInContext(req.Context())

	pId := form.PlayerId
	rId := form.RoomId
	location := form.Location
	
	//查找玩家是否在房间内
	p := rms.GetPlayerById(pId)
	if p != nil {
		rr := gamehttputils.NewFailedResult(PlayerAlreadyInRoomErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithField("playerId", pId).Warn("玩家已经在房间内")
		return
	}

	us := userservice.UserServiceInContext(req.Context())

	//判断是否被封号
	forbid, err := us.IsForbid(pId)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"playerId": pId,
			"roomId":   rId,
			"error":    err,
		}).Error("加入房间,判断是否被封号失败")
		return
	}
	if forbid {
		rr := gamehttputils.NewFailedResult(AccountIsLockErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		return
	}

	//判断是否房间存在
	r := rms.GetRoomById(rId)
	if r == nil {
		rr := gamehttputils.NewFailedResult(RoomNoExistErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithField("roomId", rId).Warn("房间不存在")
		return
	}
	
	//人数已久满了
	if r.Full() {
		rr := gamehttputils.NewFailedResult(RoomFullErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithField("roomId", rId).Warn("房间人数已满")
		return 
	}
	
	ip := req.RemoteAddr
	
	//判断是否同IP禁止
	f := rms.IfForbidIp(rId,ip)   
	if f { 
		rr := gamehttputils.NewFailedResult(JoinRoomIpErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithField("roomId", rId).Warn("同IP禁止加入房间")
		return
	}
	
	//跟新地址
	err = us.ChangeUserLocation(pId, location)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"playerId":    form.PlayerId, 
			"location" : form.Location,
			"error":      err,
		}).Error("加入房间,更新位置失败")
	}  
	
	r = rms.JoinRoom(rId, pId , ip)

	if r == nil {
		//加入房间失败
		rr := gamehttputils.NewFailedResult(JoinRoomErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithField("playerId", pId).Warn("加入房间失败")
		return
	}

	// //创建玩家
	// p = roommanageservice.NewPlayer(pId, int64(0))

	// //判断是否房间存在
	// r := rms.GetRoomById(rId)
	// if r == nil {
	// 	rr := gamehttputils.NewFailedResult(RoomNoExistErrorCode)
	// 	httputils.WriteJSON(rw, http.StatusOK, rr)
	// 	log.WithField("roomId", rId).Warn("房间不存在")
	// 	return
	// }

	// flag := r.JoinPlayer(p)
	// if !flag {
	// 	//加入房间失败
	// 	rr := gamehttputils.NewFailedResult(JoinRoomErrorCode)
	// 	httputils.WriteJSON(rw, http.StatusOK, rr)
	// 	log.WithField("playerId", pId).Warn("加入房间失败")
	// 	return
	// }
	// rms.AddPlayer(p)
	result := &joinRoomResponse{}
	sc := rms.GetServerByServerId(r.ServerId())
	result.Host = sc.Host
	result.Port = sc.Port
	rr := gamehttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)

	err = us.Online(pId)
	if err != nil {
		log.WithFields(log.Fields{
			"playerId": pId,
			"error":    err,
		}).Warn("加入房间,在线失败")
	}
}
