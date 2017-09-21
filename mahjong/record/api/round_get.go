package api

import (
	"net/http"

	recordmodel "game/mahjong/record/model"
	recordservice "game/mahjong/record/service"
	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type Round struct {
	Id         int64  `json:"id"`
	RoomType   int64  `json:"roomType"`
	RoomId     int64  `json:"roomId"`
	Round      int32  `json:"round"`
	TotalRound int32  `json:"totalRound"`
	Logs       string `json:"logs"`
	Settle     string `json:"settle"`
	Config     string `json:"config"`
	CreateTime int64  `gorm:"column:createTime"`
}

func ConvertRoundModel(rm *recordmodel.RoundModel) *Round {
	r := &Round{}
	r.Id = rm.Id
	r.RoomType = rm.RoomType
	r.RoomId = rm.RoomId
	r.Round = rm.Round
	r.TotalRound = rm.TotoalRound
	r.Settle = rm.Settle
	r.Logs = rm.Logs
	r.Config = rm.Config
	r.CreateTime = rm.CreateTime
	return r
}

type RoundGetForm struct {
	RoundId int64 `form:"roundId"`
}

func handleRoundGet(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求单局录像")

	form := &RoundGetForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("请求单局录像,参数错误")
		return
	}

	roundId := form.RoundId

	rs := recordservice.RecordServiceInContext(req.Context())
	rm, err := rs.GetRound(roundId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"roundId": roundId,
			"error":   err.Error(),
		}).Error("请求单局录像,错误")
		return
	}
	if rm == nil {
		errR := gamepkghttputils.NewFailedResult(RoundNoFoundErrorCode)
		httputils.WriteJSON(rw, http.StatusOK, errR)
		log.WithFields(log.Fields{
			"roundId": roundId,
		}).Warn("请求单局,没找到")
		return
	}
	r := ConvertRoundModel(rm)
	restR := gamepkghttputils.NewSuccessResult(r)
	httputils.WriteJSON(rw, http.StatusOK, restR)
	log.Debug("请求单局录像成功")
}
