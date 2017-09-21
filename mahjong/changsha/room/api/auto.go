package api

import (
	"net/http"

	changsharoomservice "game/mahjong/changsha/room/service"
	changshatemplateservice "game/mahjong/changsha/template/service"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

func handleAuto(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求自动加入长沙房间")

	form := &CreateForm{}

	if err := httputils.Bind(req, form); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Warn("请求创建长沙房间 解析错误")
		return
	}

	csrs := changsharoomservice.ChangShaRoomServiceInContext(req.Context())

	ownerId := userservice.UserInContext(req.Context())

	csts := changshatemplateservice.ChangShaTemplateServiceInContext(req.Context())

	roundId := form.RoundId

	roundTemplate := csts.GetRoundTemplateById(roundId)

	if roundTemplate == nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(
			log.Fields{}).Warn("请求创建长沙房间 盘数模版不存在")
		return
	}

	round := roundTemplate.Round
	cost := roundTemplate.Cost

	peopleId := form.PeopleId
	peopleTemplate := csts.GetPeopleTemplateById(peopleId)
	if peopleTemplate == nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(
			log.Fields{}).Warn("请求创建长沙房间 人数模版不存在")
		return
	}
	people := peopleTemplate.People

	zhuaNiaoId := form.ZhuaNiaoId
	// zhuaNiaoTemplate := csts.GetZhuaNiaoTemplateById(zhuaNiaoId)
	// zhuaNiao := 0
	// if zhuaNiaoTemplate != nil {
	// 	zhuaNiao = zhuaNiaoTemplate.ZhuaNiao
	// }

	zhuanXian := form.ZhuangXian
	zhuaNiaoAlg := form.ZhuaNiaoAlg
	if !valid(zhuanXian, zhuaNiaoId) {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(
			log.Fields{}).Warn("请求创建长沙房间 抓鸟不对")
		return
	}
	
	forbidIp := form.ForbidIp
	location := form.Location
	openRoomType := form.OpenRoomType
	
	result, err := csrs.AutoRoom(ownerId, round, cost, people, zhuanXian, zhuaNiaoId, zhuaNiaoAlg,forbidIp,location,openRoomType)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Error("请求创建长沙房间 失败")
		return
	}

	httputils.WriteJSON(rw, http.StatusOK, result)
	log.Debug("请求创建长沙房间 成功")
}
