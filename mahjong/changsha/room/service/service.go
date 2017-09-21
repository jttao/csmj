package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"game/mahjong/changsha"
	changshatemplateservice "game/mahjong/changsha/template/service"
	gamepkghttputils "game/pkg/httputils"
)

type ChangShaRoomService interface {
	AutoRoom(ownerId int64, round int, cost int, people int, zhuaXian bool, zhuaNiao int, zhuaNiaoAlg bool,forbidIp int,location string,openRoomType int) (*gamepkghttputils.RestResult, error)

	CreateRoom(ownerId int64, round int, cost int, people int, zhuaXian bool, zhuaNiao int, zhuaNiaoAlg bool,forbidIp int,location string,openRoomType int) (*gamepkghttputils.RestResult, error)
	
	JoinRoom(playerId int64, roomId int64 , location string) (*gamepkghttputils.RestResult, error)
}

const (
	key = "changsha_room_service"
)

func WithChangShaRoomService(ctx context.Context, csrts ChangShaRoomService) context.Context {
	return context.WithValue(ctx, key, csrts)
}

func ChangShaRoomServiceInContext(ctx context.Context) ChangShaRoomService {
	us, ok := ctx.Value(key).(ChangShaRoomService)
	if !ok {
		return nil
	}
	return us
}

type ChangShaRoomConfig struct {
	RoomManageCenter string `roomManageCenter`
}

type changshaRoomService struct {
	config         *ChangShaRoomConfig
	tempateService changshatemplateservice.ChangShaRoomTemplateService
}

func (csrs *changshaRoomService) CreateRoom(ownerId int64, round int, cost int, people int, zhuangXian bool, zhuaNiao int, zhuaNiaoAlg bool,forbidIp int,location string,openRoomType int) (*gamepkghttputils.RestResult, error) {
	url := csrs.config.RoomManageCenter + "/api/roommanage/create"

	formValue := make(map[string]interface{})

	formValue["roomType"] = 0
	formValue["ownerId"] = ownerId
	formValue["maxPlayers"] = people
	formValue["round"] = round
	formValue["cost"] = cost
	formValue["forbidIp"] = forbidIp
	formValue["location"] = location
	formValue["openRoomType"] = openRoomType
	
	customRoomConfig := &changsha.CustomRoomConfig{}
	customRoomConfig.ZhuaNiao = zhuaNiao
	customRoomConfig.ZhuangXian = zhuangXian
	customRoomConfig.ZhuaNiaoAlg = zhuaNiaoAlg
	configBytes, err := json.Marshal(customRoomConfig)
	if err != nil {
		return nil, err
	}
	formValue["roomConfig"] = string(configBytes)
	content, err := json.Marshal(formValue)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	restResult := &gamepkghttputils.RestResult{}
	err = json.Unmarshal(respBody, restResult)
	if err != nil {
		return nil, err
	}

	return restResult, nil
}

func (csrs *changshaRoomService) JoinRoom(playerId int64, roomId int64, location string ) (*gamepkghttputils.RestResult, error) {
	url := csrs.config.RoomManageCenter + "/api/roommanage/join"

	formValue := make(map[string]interface{})

	formValue["playerId"] = playerId
	formValue["roomId"] = roomId
	formValue["location"] = location 

	content, err := json.Marshal(formValue)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	restResult := &gamepkghttputils.RestResult{}
	err = json.Unmarshal(respBody, restResult)
	if err != nil {
		return nil, err
	}

	return restResult, nil
}

func (csrs *changshaRoomService) AutoRoom(ownerId int64, round int, cost int, people int, zhuangXian bool, zhuaNiao int, zhuaNiaoAlg bool,forbidIp int,location string,openRoomType int) (*gamepkghttputils.RestResult, error) {
	url := csrs.config.RoomManageCenter + "/api/roommanage/auto"

	formValue := make(map[string]interface{})

	formValue["roomType"] = 0
	formValue["ownerId"] = ownerId
	formValue["maxPlayers"] = people
	formValue["round"] = round
	formValue["cost"] = cost
	formValue["forbidIp"] = forbidIp
	formValue["location"] = location
	formValue["openRoomType"] = openRoomType
	
	customRoomConfig := &changsha.CustomRoomConfig{}
	customRoomConfig.ZhuaNiao = zhuaNiao
	customRoomConfig.ZhuangXian = zhuangXian
	customRoomConfig.ZhuaNiaoAlg = zhuaNiaoAlg
	configBytes, err := json.Marshal(customRoomConfig)
	if err != nil {
		return nil, err
	}
	formValue["roomConfig"] = string(configBytes)
	content, err := json.Marshal(formValue)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	restResult := &gamepkghttputils.RestResult{}
	err = json.Unmarshal(respBody, restResult)
	if err != nil {
		return nil, err
	}

	return restResult, nil
}

func NewChangShaRoomService(config *ChangShaRoomConfig, templateService changshatemplateservice.ChangShaRoomTemplateService) ChangShaRoomService {
	csrs := &changshaRoomService{}
	csrs.config = config
	csrs.tempateService = templateService
	return csrs
}
