package roomcard

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type TodayUseResponse struct {
	Num int64 `json:"num"`
}

//获取当日使用数量
func handleTodayUse(rw http.ResponseWriter, req *http.Request) {
	log.WithFields(
		log.Fields{
			"ip": req.RemoteAddr,
		}).Debug("请求当日使用数量")
	rcs := gmservice.RoomCardServiceInContext(req.Context())
	num, err := rcs.GetTodayUse()
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("请求当日使用数量,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &TodayUseResponse{}
	result.Num = num

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":  req.RemoteAddr,
			"num": num,
		}).Debug("请求当日使用数量,成功")
}
