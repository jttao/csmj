package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	
	gamepkghttputils "game/pkg/httputils"
	"game/roommanage/api"

)

type RoomManageClient interface {
	Create(RoomType int, ownerId int, maxPlayers int) (int64, error)
	Query(playerId int64) (roomId int64, serverId string, host string, port int, maxPlayers int, round int, roomConfig string,forbidIp int,openRoomType int,createTime int64,forbidJoinTime int64,lastGameTime int64,ownerId int64,err error)
	Leave(playerId int64, roomId int64) error
	Destroy(roomId int64, refund bool) error 
}

type RoomManageClientConfig struct {
	RoomManageCenter string `json:"roomManageCenter"`
}

type roomManageClient struct {
	config *RoomManageClientConfig
}

func (rmc *roomManageClient) Create(roomType int, ownerId int, maxPlayers int) (roomId int64, err error) {
	// apiPath := "/api/roommanage/query"
	// apiPath = "http://" + rmc.config.RoomManageCenter + apiPath
	// form := api.CreateRoomForm{
	// 	RoomType: roomType,
	// }
	// result, err := gamepkghttputils.PostJson(apiPath, nil, form)
	// if err != nil {
	// 	return 0, err
	// }

	// query, ok := result.(map[string]interface{})
	// if !ok {
	// 	return 0, fmt.Errorf("response error")
	// }

	// roomId = int64(query["roomId"].(float64))

	// return roomId, nil
	return
}

func (rmc *roomManageClient) Query(playerId int64) (roomId int64, serverId string, host string, port int, maxPlayer int, round int, roomConfig string,forbidIp int,openRoomType int,createTime int64,forbidJoinTime int64,lastGameTime int64,ownerId int64,err error) {
	apiPath := "/api/roommanage/query"
	apiPath = "http://" + rmc.config.RoomManageCenter + apiPath
	form := api.QueryForm{
		PlayerId: playerId,
	}
	result, err := gamepkghttputils.PostJson(apiPath, nil, form)
	if err != nil {
		return 0, "", "", 0, 0, 0, "",0,0,0,0,0,0, err
	}

	query, ok := result.(map[string]interface{})
	if !ok {
		return 0, "", "", 0, 0, 0, "",0,0,0,0,0,0, fmt.Errorf("response error")
	}

	roomId = int64(query["roomId"].(float64))
	if roomId == 0 {
		return 0, "", "", 0, 0, 0, "",0,0,0,0,0,0, nil
	}
	if query["round"] == nil {
		return 0, "", "", 0, 0, 0, "",0,0,0,0,0,0, fmt.Errorf("round error")
	}
	round = int(query["round"].(float64))
	if query["roomConfig"] == nil {
		return 0, "", "", 0, 0, 0, "",0,0,0,0,0,0, fmt.Errorf("round config error")
	}
	roomConfig = query["roomConfig"].(string)
	if query["maxPlayers"] == nil {
		return 0, "", "", 0, 0, 0, "",0,0,0,0,0,0, fmt.Errorf("max players error")
	}
	maxPlayer = int(query["maxPlayers"].(float64))

	serverId = query["serverId"].(string)
	host = query["host"].(string)
	port = int(query["port"].(float64))
	forbidIp = int(query["forbidIp"].(float64))
	openRoomType = int(query["openRoomType"].(float64)) 
	createTime = int64(query["createTime"].(float64)) 
	forbidJoinTime = int64(query["forbidJoinTime"].(float64))
	lastGameTime = int64(query["lastGameTime"].(float64))
	ownerId = int64(query["ownerId"].(float64))
	
	return roomId, serverId, host, port, maxPlayer, round, roomConfig,forbidIp,openRoomType,createTime,forbidJoinTime,lastGameTime,ownerId,nil
}

func (rmc *roomManageClient) Leave(playerId int64, roomId int64) error {
	apiPath := "/api/roommanage/leave"
	apiPath = "http://" + rmc.config.RoomManageCenter + apiPath
	form := api.LeaveRoomForm{
		PlayerId: playerId,
		RoomId:   roomId,
	}
	_, err := gamepkghttputils.PostJson(apiPath, nil, form)
	if err != nil {
		return err
	}

	return nil
}

func (rmc *roomManageClient) Destroy(roomId int64, refund bool) error {
	apiPath := "/api/roommanage/destroy"
	apiPath = "http://" + rmc.config.RoomManageCenter + apiPath
	form := api.DestroyRoomForm{
		RoomId: roomId,
		Refund: refund,
	}
	_, err := gamepkghttputils.PostJson(apiPath, nil, form)
	if err != nil {
		return err
	}
	return nil
}

func NewRoomManageClient(config *RoomManageClientConfig) RoomManageClient {
	rmc := &roomManageClient{}
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
	key = "room_manage_client"
)

func WithRoomManageClient(ctx context.Context, rmc RoomManageClient) context.Context {
	return context.WithValue(ctx, key, rmc)
}

func RoomManageClientInContext(ctx context.Context) RoomManageClient {
	us, ok := ctx.Value(key).(RoomManageClient)
	if !ok {
		return nil
	}
	return us
}
