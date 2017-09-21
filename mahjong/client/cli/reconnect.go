package cli

import (
	"fmt"
	"log"

	gamepkghttputils "game/pkg/httputils"

	"github.com/codegangsta/cli"
)

var (
	reconnectCmd = cli.Command{
		Name:        "reconnect",
		Usage:       "reconnect ${roomId}",
		Description: "reconnect room",
		Action:      reconnect,
		Flags:       []cli.Flag{},
	}
)

func init() {
	appendCmd(reconnectCmd)
}

func reconnect(c *cli.Context) {

	token, err := login(deviceMac)
	if err != nil {
		log.Fatalln("login failed ", err)
	}
	roomId, host, port, err := hallLogin(token)
	if err != nil {
		log.Fatalln("hall login failed", err)
	}
	if roomId == 0 {
		log.Fatalln("player no in room")
	}
	// go func() {
	// 	chatAddr := serverAddr + ":5000"
	// 	connectChatServer(chatAddr, token)
	// }()
	log.Println("join room successful,roomId ", roomId)

	roomAddr := fmt.Sprintf("%s:%d", host, port)

	connectTargetServer(roomAddr, token)

}

func hallLogin(token string) (roomId int64, host string, port int, err error) {
	loginApi := "http://" + serverAddr + ":81" + "/api/hall/login"

	header := make(map[string]string)
	header["Authorization"] = "BEARER " + token
	result, err := gamepkghttputils.PostJson(loginApi, header, nil)
	if err != nil {
		return
	}
	resultMap := result.(map[string]interface{})
	roomId = int64(resultMap["roomId"].(float64))
	host = resultMap["host"].(string)
	port = int(resultMap["port"].(float64))
	return
}
