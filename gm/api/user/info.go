package user

import (
	gmservice "game/gm/service"
	"net/http"

	gamepkghttputils "game/pkg/httputils"

	usermodel "game/user/model"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type UserInfoForm struct {
	UserId int64 `form:"userId"`
}

type UserInfoResponse struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	CardNum       int64  `json:"cardNum"`
	State         int    `json:"state"`
	LastLoginTime int64  `json:"lastLoginTime"`
	LastLoginIp   string `json:"lastLoginIp"`
	Forbid        int    `json:"forbid"`
	CreateTime    int64  `json:"createTime"`
}

func convertUserToResponse(u *usermodel.User) *UserInfoResponse {
	uir := &UserInfoResponse{}
	uir.CardNum = u.CardNum
	uir.CreateTime = u.CreateTime
	uir.Id = u.Id
	uir.Name = u.Name
	uir.State = u.State
	uir.LastLoginIp = u.LastLoginIp
	uir.LastLoginTime = u.LastLoginTime
	uir.Forbid = u.Forbid
	return uir
}

//禁
func handleInfo(rw http.ResponseWriter, req *http.Request) {
	form := &UserInfoForm{}
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
	userId := form.UserId

	user, err := us.GetUserById(userId)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"userId": userId,
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
				"userId": userId,
			}).Debug("请求用户信息,查无此人")
	} else {
		result := convertUserToResponse(user)
		rr := gamepkghttputils.NewSuccessResult(result)
		httputils.WriteJSON(rw, http.StatusOK, rr)
		log.WithFields(
			log.Fields{
				"ip":     req.RemoteAddr,
				"userId": userId,
			}).Debug("请求用户信息,成功")
	}
}
