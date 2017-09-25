package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	gamepkghttputils "game/pkg/httputils" 

	"game/hall/api" 
	log "github.com/Sirupsen/logrus"

)

type HallClient interface { 
	GameEnd(playerId int64, flag bool) error
}

type HallClientConfig struct {
	HallClientCenter       string      `json:"hallClientCenter"` 
}

type hallClient struct {
	config *HallClientConfig
}

func (rmc *hallClient) GameEnd(playerId int64, maxwin bool) error {
	
	apiPath := "/api/hall/task_finish"
	apiPath = "http://" + rmc.config.HallClientCenter + apiPath
	
	log.WithFields(
		log.Fields{
			"playerId":   playerId, 
			"maxwin":	maxwin, 
		}).Debug("开始设置玩家任务> ")

	//每日游戏任务
	taskId := int32(2)
	state := true 
	form := api.TaskFinishForm {
		PlayerId : playerId,
		TaskShareId: taskId,
		State: state,
	} 	
	
	_, err := gamepkghttputils.PostJson(apiPath, nil, form)
	if err != nil {
		return err
	}

	//每日连续赢任务
	taskId = int32(2)
	state = maxwin 
	form = api.TaskFinishForm {
		PlayerId : playerId,
		TaskShareId: taskId,
		State: state,
	}   

	_, err1 := gamepkghttputils.PostJson(apiPath, nil, form)

	if err1 != nil {
		return err1
	}
	
	return nil
}

func NewHallClient(config *HallClientConfig) HallClient {
	rmc := &hallClient{}
	rmc.config = config
	return rmc
}

func postForm(apiPath string, form interface{}) (result interface{}, err error) {

	bodyBytes, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(apiPath, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status failed")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rr := &gamepkghttputils.RestResult{}
	err = json.Unmarshal(respBody, rr)
	if err != nil {
		return nil, err
	}
	if rr.ErrorCode != 0 {
		return rr.Result, fmt.Errorf("error_code %d", rr.ErrorCode)
	}
	return rr.Result, nil
}

const (
	key = "hall_client"
)

func WithHallClient(ctx context.Context, rmc HallClient) context.Context {
	return context.WithValue(ctx, key, rmc)
}

func HallClientInContext(ctx context.Context) HallClient {
	us, ok := ctx.Value(key).(HallClient)
	if !ok {
		return nil
	}
	return us
}
