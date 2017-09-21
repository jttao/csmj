package api

import (
	"net/http"

	pkghttputils "game/pkg/httputils"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleRefreshToken(rw http.ResponseWriter, req *http.Request) {
	log.Debug("刷新token")
	us := userservice.UserServiceInContext(req.Context())

	userId := userservice.UserInContext(req.Context())
	t, expiredTime, err := us.Login(userId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"userId": userId,
			"error":  err,
		}).Error("刷新token,失败")
		return
	}
	r := &LoginResponse{}
	r.ExpireTime = expiredTime
	r.Token = t
	rr := pkghttputils.NewSuccessResult(r)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.Debug("刷新token,成功")
}

func refreshTokenHandler() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		userservice.AuthHandler().ServeHTTP(rw, req, http.HandlerFunc(handleRefreshToken))
	})
}
