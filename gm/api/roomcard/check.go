package roomcard

import (
	"net/http"

	gmservice "game/gm/service"
	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type CheckForm struct {
	Flag bool `form:"flag"`
}

//是否审核
func handleCheck(rw http.ResponseWriter, req *http.Request) {
	form := &CheckForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("请求审核")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	flag := form.Flag

	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"flag": flag,
		}).Debug("请求审核")
	rcs := gmservice.RoomCardServiceInContext(req.Context())
	err = rcs.Check(flag)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("请求审核,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &struct{}{}

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"flag": flag,
		}).Debug("请求审核,成功")
}

type IfCheckResponse struct {
	Flag bool `json:"flag"`
}

func handleIfCheck(rw http.ResponseWriter, req *http.Request) {

	log.WithFields(
		log.Fields{
			"ip": req.RemoteAddr,
		}).Debug("请求审核")
	rcs := gmservice.RoomCardServiceInContext(req.Context())
	flag, err := rcs.IfCheck()
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("请求审核,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &IfCheckResponse{}
	result.Flag = flag

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"flag": flag,
		}).Debug("请求审核,成功")
}
