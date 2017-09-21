package api

import (
	"log"
	"net/http"
	"strings"

	userservice "game/user/service"

	"github.com/xozrc/pkg/httputils"
)

type WxLoginForm struct {
	WxId      string `form:"wxId"`
	WxName    string `form:"wxName"`
	WxSex     int    `form:"wxSex"`
	WxImage   string `form:"wxImage"`
	DeviceMac string `form:"deviceMac"`
}

func handleWxLogin(rw http.ResponseWriter, req *http.Request) {
	loginForm := &WxLoginForm{}
	if err := httputils.Bind(req, loginForm); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	deviceMac := loginForm.DeviceMac
	wxId := loginForm.WxId
	wxName := loginForm.WxName
	deviceMac = strings.TrimSpace(deviceMac)
	//TODO 验证参数
	if len(deviceMac) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	wxSex := loginForm.WxSex
	wxImage := loginForm.WxImage

	if len(wxId) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(wxName) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	wxSexEnum := userservice.Sex(wxSex)
	if !wxSexEnum.Valid() {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	us := userservice.UserServiceInContext(req.Context())
	if us == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	t, expiredTime, err := us.WxLogin(wxId, deviceMac, wxName, wxSexEnum, wxImage)
	if err != nil {
		log.Println("wx login error ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	lr := &LoginResponse{}
	lr.Token = t
	lr.ExpireTime = expiredTime

	rr := RestResult{}
	rr.ErrorCode = 0
	rr.Result = lr
	httputils.WriteJSON(rw, http.StatusOK, rr)
}
