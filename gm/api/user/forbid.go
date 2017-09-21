package user

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type ForbidUserForm struct {
	UserId int64 `form:"userId"`
	Flag   bool  `form:"flag"`
}

//禁
func handleForbid(rw http.ResponseWriter, req *http.Request) {
	form := &ForbidUserForm{}
	err := httputils.Bind(req, form)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":    req.RemoteAddr,
				"error": err,
			}).Warn("禁止用户,解析失败")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.WithFields(
		log.Fields{
			"ip":   req.RemoteAddr,
			"form": form,
		}).Debug("请求禁止用户")

	us := gmservice.UserServiceInContext(req.Context())
	userId := form.UserId

	//判断用户是否存在
	user, err := us.GetUserById(userId)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"userId": userId,
				"error":  err,
			}).Error("禁止用户,失败")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"userId": userId,
			}).Warn("禁止用户,用户不存在")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	flag := form.Flag
	if flag {
		err := us.ForbidUser(userId)
		if err != nil {
			log.WithFields(
				log.Fields{
					"ip":     req.RemoteAddr,
					"userId": userId,
					"flag":   flag,
					"error":  err,
				}).Error("禁止用户,失败")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := us.UnforbidUser(userId)
		if err != nil {
			log.WithFields(
				log.Fields{
					"ip":     req.RemoteAddr,
					"userId": userId,
					"flag":   flag,
					"error":  err,
				}).Error("禁止用户,失败")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	result := &struct{}{}

	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.WithFields(
		log.Fields{
			"ip":     req.RemoteAddr,
			"userId": userId,
			"flag":   flag,
		}).Debug("禁止用户,成功")
}
