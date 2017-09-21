package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	recordmodel "game/mahjong/record/model"
	recordservice "game/mahjong/record/service"
	gamepkghttputils "game/pkg/httputils"
	"github.com/xozrc/pkg/httputils"
)

type RoundBasic struct {
	Id         int64  `json:"id"`
	RoomType   int64  `json:"roomType"`
	RoomId     int64  `json:"roomId"`
	Round      int32  `json:"round"`
	Settle     string `json:"settle"`
	CreateTime int64  `gorm:"column:createTime"`
}

func ConvertRoundModelToBasic(rm *recordmodel.RoundModel) *RoundBasic {
	r := &RoundBasic{}
	r.Id = rm.Id
	r.RoomType = rm.RoomType
	r.RoomId = rm.RoomId
	r.Round = rm.Round
	r.Settle = rm.Settle
	r.CreateTime = rm.CreateTime
	return r
}

type RecordGetForm struct {
	RoomId int64 `form:"roomId"`
}

func handleRecordGet(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求录像")

	form := &RecordGetForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("请求录像,参数错误")
		return
	}

	roomId := form.RoomId

	rs := recordservice.RecordServiceInContext(req.Context())
	rml, err := rs.GetRoundList(roomId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"roomId": roomId,
			"error":  err.Error(),
		}).Error("请求录像,错误")
		return
	}
	rl := make([]*RoundBasic, 0, len(rml))
	for _, rm := range rml {
		r := ConvertRoundModelToBasic(rm)
		rl = append(rl, r)
	}
	restR := gamepkghttputils.NewSuccessResult(rl)
	httputils.WriteJSON(rw, http.StatusOK, restR)
	log.Debug("请求录像成功")
}
