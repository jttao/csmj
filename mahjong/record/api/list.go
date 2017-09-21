package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	recordmodel "game/mahjong/record/model"
	recordservice "game/mahjong/record/service"
	gamepkghttputils "game/pkg/httputils"
	userservice "game/user/service"
	"github.com/xozrc/pkg/httputils"
)

type RoomRecord struct {
	Id         int64  `json:"id"`
	RoomType   int64  `json:"roomType"`
	RoomId     int64  `json:"roomId"`
	OwnerId    int64  `json:"ownerId"`
	Settle     string `json:"settle"`
	CreateTime int64  `gorm:"column:createTime"`
}

func Convert(rrm *recordmodel.RoomRecordModel) *RoomRecord {
	rr := &RoomRecord{}
	rr.Id = rrm.Id
	rr.RoomType = rrm.RoomType
	rr.RoomId = rrm.RoomId
	rr.OwnerId = rrm.OwnerId
	rr.Settle = rrm.Settle
	rr.CreateTime = rrm.CreateTime
	return rr
}

func handleRecordList(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求录像列表")

	playerId := userservice.UserInContext(req.Context())

	rs := recordservice.RecordServiceInContext(req.Context())
	rrml, err := rs.GetRecordList(playerId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"playerId": playerId,
			"error":    err.Error(),
		}).Error("请求录像列表,错误")
		return
	}
	rrl := make([]*RoomRecord, 0, len(rrml))
	for _, rrm := range rrml {
		rr := Convert(rrm)
		rrl = append(rrl, rr)
	}
	restR := gamepkghttputils.NewSuccessResult(rrl)
	httputils.WriteJSON(rw, http.StatusOK, restR)
	log.Debug("请求录像列表成功")
}
