package roomcard

import (
	"net/http"

	gmservice "game/gm/service"
	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type FreeForm struct {
	Flag bool `form:"flag"`
}

//是否免房卡
func handleFree(rw http.ResponseWriter, req *http.Request) {
	form := &FreeForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("强求免房卡")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	flag := form.Flag

	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"flag": flag,
		}).Debug("请求免房卡")
	rcs := gmservice.RoomCardServiceInContext(req.Context())
	err = rcs.Free(flag)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("请求免房卡,失败")
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
		}).Debug("请求免房卡,成功")
}

type IfFreeResponse struct {
	Flag bool `json:"flag"`
}

func handleIfFree(rw http.ResponseWriter, req *http.Request) {

	log.WithFields(
		log.Fields{
			"ip": req.RemoteAddr,
		}).Debug("请求免房卡")
	rcs := gmservice.RoomCardServiceInContext(req.Context())
	flag, err := rcs.IfFree()
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Error("请求免房卡,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &IfFreeResponse{}
	result.Flag = flag

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"flag": flag,
		}).Debug("请求免房卡,成功")
}
