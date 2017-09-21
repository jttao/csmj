package roomcard

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type ChangeCardForm struct {
	UserId int64 `form:"userId"`
	Change int64 `form:"change"`
}

//改变卡牌数
func handleChange(rw http.ResponseWriter, req *http.Request) {
	form := &ChangeCardForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("请求改变卡牌数,解析失败")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := form.UserId
	change := form.Change
	log.WithFields(
		log.Fields{
			"ip":     req.RemoteAddr,
			"userId": userId,
			"change": change,
		}).Debug("请求改变卡牌数")
	rcs := gmservice.RoomCardServiceInContext(req.Context())
	err = rcs.ChangeCardNum(userId, change)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"userId": userId,
				"change": change,
				"error":  err,
			}).Error("请求改变卡牌数,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &struct{}{}

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":     req.RemoteAddr,
			"userId": userId,
			"change": change,
		}).Debug("请求改变卡牌数,成功")
}
