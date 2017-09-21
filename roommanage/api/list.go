package api

import (
	"fmt"
	"net/http"

	gamehttputils "game/pkg/httputils"
	roommanageservice "game/roommanage/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type roomListForm struct {
	ServerId int `form:"serverId"`
}

type roomListResponse struct {
	Rooms []*roomResponse `json:"rooms"`
}

type roomResponse struct {
	Id int64 `json:"id"`
}

func toRoomResponse(r roommanageservice.Room) *roomResponse {
	rp := &roomResponse{}
	rp.Id = r.Id()
	return rp
}

func toRoomListResponse(rs []roommanageservice.Room) *roomListResponse {
	rlp := &roomListResponse{}
	rlp.Rooms = make([]*roomResponse, 0, len(rs))
	for _, r := range rs {
		rp := toRoomResponse(r)
		rlp.Rooms = append(rlp.Rooms, rp)
	}
	return rlp
}

func handleList(rw http.ResponseWriter, req *http.Request) {
	form := &roomListForm{}
	if err := httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Warn("create room request bind failed")
		return
	}
	log.WithField("form", fmt.Sprintf("%#v", form)).Debug("请求房间列表")

	//serverId := form.ServerId

	rms := roommanageservice.RoomManageInContext(req.Context())
	rs := rms.Rooms()

	result := toRoomListResponse(rs)

	gamehttputils.NewSuccessResult(result)

	httputils.WriteJSON(rw, http.StatusOK, result)
	log.Debug("请求房间列表成功")
}
