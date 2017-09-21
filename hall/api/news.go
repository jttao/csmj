package api

import (
	"net/http"

	news "game/hall/news"
	gamepkghttputils "game/pkg/httputils"
	
	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type NewsResponse struct {
	News string `json:"news"`
}

func handleNews(rw http.ResponseWriter, req *http.Request) {

	ns := news.NewsServiceInContext(req.Context())
	news, err := ns.GetNews()
	if err != nil {
		log.Error("get news  error ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := &NewsResponse{
		News: news,
	}
	rr := gamepkghttputils.NewSuccessResult(result)
	httputils.WriteJSON(rw, http.StatusOK, rr)
}
