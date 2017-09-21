package user

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type WxUserInfoForm struct {
	WxId string `form:"wxId"`
}

//禁
func handleInfoByWx(rw http.ResponseWriter, req *http.Request) {
	form := &WxUserInfoForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("请求用户信息,解析失败")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"form": form,
		}).Debug("请求用户信息")

	us := gmservice.UserServiceInContext(req.Context())
	wxId := form.WxId

	user, err := us.GetUserByWxId(wxId)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"weixin": wxId,
				"error":  err,
			}).Error("请求用户信息,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		result := &struct{}{}
		rr := gamepkghttputils.NewSuccessResult(result)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"weixin": wxId,
			}).Debug("请求用户信息,查无此人")
	} else {
		result := convertUserToResponse(user)
		rr := gamepkghttputils.NewSuccessResult(result)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"weixin": wxId,
			}).Debug("请求用户信息,成功")
	}
}
