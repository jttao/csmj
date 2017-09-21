package api

import (
	"net/http"

	pkghttputils "game/pkg/httputils"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleLogout(rw http.ResponseWriter, req *http.Request) {
	log.Debug("登出")
	us := userservice.UserServiceInContext(req.Context())

	userId := userservice.UserInContext(req.Context())
	err := us.Logout(userId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"用户Id":  userId,
			"error": err,
		}).Error("登出,失败")
		return
	}
	rr := pkghttputils.NewSuccessResult(nil)
	httputils.WriteJSON(rw, http.StatusOK, rr)
	log.Debug("登出,成功")
}

func logoutHandler() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		userservice.AuthHandler().ServeHTTP(rw, req, http.HandlerFunc(handleLogout))
	})
}
