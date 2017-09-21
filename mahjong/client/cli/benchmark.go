package cli

import (
	"fmt"
	"sync"

	changsharoomapi "game/mahjong/changsha/room/api"
	gamepkghttputils "game/pkg/httputils"

	log "github.com/Sirupsen/logrus"

	"github.com/codegangsta/cli"
)

var (
	wg           sync.WaitGroup
	benchmarkNum = 1
)
var (
	benchmarkCommand = cli.Command{
		Name:        "benchmark",
		Usage:       "benchmark",
		Description: "benchmark",
		Action:      benchmark,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "number",
				Value:       benchmarkNum,
				Usage:       "number",
				Destination: &benchmarkNum,
			},
		},
	}
)

func init() {
	appendCmd(benchmarkCommand)
}

func benchmark(c *cli.Context) {
	wg.Add(benchmarkNum)
	actulNum := 0
	for i := 0; i < benchmarkNum; i++ {
		tempDeviceMac := fmt.Sprintf("%s:%d", deviceMac, i)
		token, err := login(tempDeviceMac)
		if err != nil {
			log.Info("login failed ", err)
			wg.Done()
			continue
		}

		roomId, host, port, err := hallLogin(token)
		if err != nil {
			log.Error("hall login", err)
			continue
		}
		roomAddr := ""
		if roomId == 0 {
			createApi := "http://" + serverAddr + ":81" + "/api/changsha/room/auto"
			form := &changsharoomapi.CreateForm{}
			form.PeopleId = 1
			form.RoundId = 1
			form.ZhuangXian = false
			form.ZhuaNiaoId = 1
			header := make(map[string]string)
			header["Authorization"] = "BEARER " + token
			result, err := gamepkghttputils.PostJson(createApi, header, form)
			if err != nil {
				log.Info("auto room failed", err)
				wg.Done()
				continue
			}
			tresult := result.(map[string]interface{})

			host := tresult["host"]
			port := int(tresult["port"].(float64))
			roomAddr = fmt.Sprintf("%s:%d", host, port)
		} else {
			roomAddr = fmt.Sprintf("%s:%d", host, port)
		}

		go func() {
			log.Println("auto room successful")
			// go func() {
			// 	chatAddr := serverAddr + ":5000"
			// 	connectChatServer(chatAddr, token)
			// }()

			connectTargetServer(roomAddr, token)
			wg.Done()
		}()
		actulNum++
	}
	log.Printf("bengin %d client\n", actulNum)
	wg.Wait()
	log.Printf("end %d client\n", actulNum)
}
