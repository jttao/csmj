package cli

import (
	"fmt"
	"log"

	changsharoomapi "game/mahjong/changsha/room/api"
	gamepkghttputils "game/pkg/httputils"

	"github.com/codegangsta/cli"
)

var (
	createCmd = cli.Command{
		Name:        "create",
		Usage:       "create",
		Description: "create room",
		Action:      create,
		Flags:       []cli.Flag{},
	}
)

func init() {
	appendCmd(createCmd)
}

func create(c *cli.Context) {
	token, err := login(deviceMac)
	if err != nil {
		log.Fatalln("login failed ", err)
	}
	createApi := "http://" + serverAddr + ":81" + "/api/changsha/room/create"
	form := &changsharoomapi.CreateForm{}
	form.PeopleId = 1
	form.RoundId = 1
	form.ZhuangXian = false
	form.ZhuaNiaoId = 1
	header := make(map[string]string)
	header["Authorization"] = "BEARER " + token
	result, err := gamepkghttputils.PostJson(createApi, header, form)
	if err != nil {
		log.Fatalln("create room failed", err)
	}
	tresult := result.(map[string]interface{})
	roomId := tresult["roomId"]
	host := tresult["host"]
	port := int(tresult["port"].(float64))
	log.Println("create room successful,roomId ", roomId)
	// go func() {
	// 	chatAddr := serverAddr + ":5000"
	// 	connectChatServer(chatAddr, token)
	// }()

	roomAddr := fmt.Sprintf("%s:%d", host, port)
	connectTargetServer(roomAddr, token)

}
