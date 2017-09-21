package api

import (
	"net/http"

	changshatemplateservice "game/mahjong/changsha/template/service"
	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleTemplate(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求长沙房间模版数据")
	csts := changshatemplateservice.ChangShaTemplateServiceInContext(req.Context())

	c := csts.GetAll()

	result := gamepkghttputils.NewSuccessResult(c)

	httputils.WriteJSON(rw, http.StatusOK, result)
	log.Debug("请求长沙房间模版数据 成功")
}
