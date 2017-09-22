package api

import (
	"fmt"
	"net/http"

	gamehttputils "game/pkg/httputils"
	roommanageservice "game/roommanage/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type QueryForm struct {
	PlayerId int64 `form:"playerId" json:"playerId"`
}

type QueryResponse struct {
	RoomId     int64  `json:"roomId"`
	OwnerId    int64  `json:"ownerId"`
	Round      int    `json:"round"`
	MaxPlayers int    `json:"maxPlayers"`
	RoomConfig string `json:"roomConfig"`
	ServerId   string `json:"serverId"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	ForbidIp   int    `json:"forbidIp"` 
	OpenRoomType int  `json:"openRoomType"`
	CreateTime 	int64 `json:"createTime"`
	ForbidJoinTime 	int64 `json:"forbidJoinTime"`
	LastGameTime int64 `json:"lastGameTime"`
}

func handleQuery(rw http.ResponseWriter, req *http.Request) {
	form := &QueryForm{}
	if err := httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Warn("create room request bind failed")
		return
	}
	log.WithField("form", fmt.Sprintf("%#v", form)).Debug("查询玩家")
	pId := form.PlayerId
	rms := roommanageservice.RoomManageInContext(req.Context())
	p := rms.GetPlayerById(form.PlayerId)
	result := QueryResponse{
		RoomId: 0,
	}
	if p == nil {
		log.WithField("playerId", pId).Warn("玩家不在房间内")
	} else { 
		r := rms.GetRoomById(p.RoomId())
		result.RoomId = p.RoomId()
		result.OwnerId = r.OwnerId()
		result.RoomConfig = r.RoomConfig()
		result.Round = r.Round()
		result.MaxPlayers = r.MaxPlayers()
		result.ForbidIp = r.ForbidIp() 
		result.OpenRoomType = int(r.GetOpenRoomType()) 
		result.CreateTime = r.CreateTime()
		result.ForbidJoinTime = r.ForbidJoinTime()
		result.LastGameTime = r.LastGameTime()  
		sc := rms.GetServerByServerId(r.ServerId())
		result.Host = sc.Host
		result.Port = sc.Port
	}
	rr := gamehttputils.NewSuccessResult(result)

	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"result": fmt.Sprintf("%#v", result),
		}).Debug("查询玩家房间")
}
