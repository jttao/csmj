package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	userPath = "/user"
)

func Router(r *mux.Router) {
	sr := r.PathPrefix(userPath).Subrouter()

	sr.Path("/wxlogin").HandlerFunc(http.HandlerFunc(handleWxLogin))
	sr.Path("/visitor").HandlerFunc(http.HandlerFunc(handleVisitorLogin))

	sr.Path("/logout").HandlerFunc(logoutHandler())
	sr.Path("/refresh_token").HandlerFunc(refreshTokenHandler())
}

type LoginResponse struct {
	Token      string `json:"token"`
	ExpireTime int64  `json:"expireTime"`
}

type RestResult struct {
	ErrorCode int         `json:"error_code"`
	Result    interface{} `json:"result"`
}
