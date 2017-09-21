package cli

import (
	"fmt"
	"log"
	"strconv"

	changsharoomapi "game/mahjong/changsha/room/api"
	gamepkghttputils "game/pkg/httputils"

	"github.com/codegangsta/cli"
)

var (
	joinCmd = cli.Command{
		Name:        "join",
		Usage:       "join ${roomId}",
		Description: "join room",
		Action:      join,
		Flags:       []cli.Flag{},
	}
)

func init() {
	appendCmd(joinCmd)
}

func join(c *cli.Context) {
	if c.NArg() < 1 {
		log.Println("need add room id")
		return
	}
	roomIdStr := c.Args().First()
	roomId, err := strconv.ParseInt(roomIdStr, 10, 64)
	if err != nil {
		log.Fatalln("room id is not integer", err)
	}

	token, err := login(deviceMac)
	if err != nil {
		log.Fatalln("login failed ", err)
	}
	joinApi := "http://" + serverAddr + ":81" + "/api/changsha/room/join"
	form := &changsharoomapi.JoinForm{}
	form.RoomId = roomId
	header := make(map[string]string)
	header["Authorization"] = "BEARER " + token
	result, err := gamepkghttputils.PostJson(joinApi, header, form)
	if err != nil {
		log.Fatalln("join room failed", err)
	}
	tresult := result.(map[string]interface{})
	log.Println("join room successful,roomId ", roomId)
	// go func() {
	// 	chatAddr := serverAddr + ":5000"
	// 	connectChatServer(chatAddr, token)
	// }()
	host := tresult["host"]
	port := int(tresult["port"].(float64))
	roomAddr := fmt.Sprintf("%s:%d", host, port)
	connectTargetServer(roomAddr, token)

}
