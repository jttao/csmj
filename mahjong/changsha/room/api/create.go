package api

import (
	"net/http"

	changsharoomservice "game/mahjong/changsha/room/service"
	changshatemplateservice "game/mahjong/changsha/template/service"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/xozrc/pkg/httputils"
)

type CreateForm struct {
	RoundId     int  `form:"roundId" json:"roundId"`
	PeopleId    int  `form:"peopleId" json:"peopleId"`
	ZhuangXian  bool `form:"zhuangXian" json:"zhuangXian"`
	ZhuaNiaoId  int  `form:"zhuaNiaoId" json:"zhuaNiaoId"`
	ZhuaNiaoAlg bool `form:"zhuaNiaoAlg" json:"zhuaNiaoAlg"`
	Location  string `form:"location" json:"location"`
	ForbidIp     int `form:"forbidIp" json:"forbidIp"`
	OpenRoomType int `form:"openRoomType" json:"openRoomType"`
}

func handleCreate(rw http.ResponseWriter, req *http.Request) {
	log.Debug("请求创建长沙房间")
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

	location := form.Location
	forbidIp := form.ForbidIp
	openRoomType := form.OpenRoomType

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

	// zhuaNiaoTemplate := csts.GetZhuaNiaoTemplateById(zhuaNiaoId)
	// zhuaNiao := 0
	// if zhuaNiaoTemplate != nil {
	// 	zhuaNiao = zhuaNiaoTemplate.ZhuaNiao
	// }

	zhuanXian := form.ZhuangXian
	zhuaNiaoAlg := form.ZhuaNiaoAlg
	zhuaNiaoId := form.ZhuaNiaoId
	

	if !valid(zhuanXian, zhuaNiaoId) {
		rw.WriteHeader(http.StatusBadRequest)
		log.WithFields(
			log.Fields{}).Warn("请求创建长沙房间 抓鸟不对")
		return
	}
	
	result, err := csrs.CreateRoom(ownerId, round, cost, people, zhuanXian, zhuaNiaoId, zhuaNiaoAlg,forbidIp,location,openRoomType)
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

func valid(zhuanXian bool, zhuaNiaoId int) bool {
	if zhuanXian {
		switch zhuaNiaoId {
		case 0:
		case 1:
		case 2:
		default:
			return false
		}
		return true
	} else {
		switch zhuaNiaoId {
		case 0:
		case 2:
		case 4:
		case 6:
		default:
			return false
		}
		return true
	}
}
