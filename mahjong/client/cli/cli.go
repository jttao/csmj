package cli

import (
	//"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/websocket"

	"game/basic/pb"
	"game/mahjong/client/client"
	clienthandler "game/mahjong/client/handler"
	roomhandler "game/mahjong/client/room/handler"
	loginpb "game/mahjong/pb/login"
	messagetypepb "game/mahjong/pb/messagetype"
	pkghttputils "game/pkg/httputils"
	userapi "game/user/api"

	"github.com/codegangsta/cli"
	"github.com/golang/protobuf/proto"
	//chatmessage "github.com/mahjong/chat/message"
)

var (
	debug        = false
	path         = "/room"
	serverAddr   = "127.0.0.1"
	deviceMac    = "test"
	timeInterval = time.Second
	timeout      = false
)

func Start() {
	app := cli.NewApp()
	app.Name = "room client"
	app.Usage = "room [global options] command [command options] [arguments...]."

	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "debug",
			Destination: &debug,
		},

		cli.StringFlag{
			Name:        "path,p",
			Value:       path,
			Usage:       "path",
			Destination: &path,
		},
		cli.StringFlag{
			Name:        "server",
			Value:       serverAddr,
			Usage:       "server",
			Destination: &serverAddr,
		},
		cli.StringFlag{
			Name:        "deviceMac",
			Value:       deviceMac,
			Usage:       "deviceMac",
			Destination: &deviceMac,
		},

		cli.DurationFlag{
			Name:        "time",
			Value:       timeInterval,
			Usage:       "time",
			Destination: &timeInterval,
		},
		cli.BoolFlag{
			Name:        "timeout",
			Usage:       "timeout",
			Destination: &timeout,
		},
	}

	//app.Action = start
	app.Commands = commands
	app.Run(os.Args)
}

// func start(c *cli.Context) {

// 	//登陆
// 	visitorPath := "http://" + serverAddr + "/api/user/visitor"
// 	visitorMap := userapi.VisitorLoginForm{
// 		DeviceMac: deviceMac,
// 	}

// 	result, err := pkghttputils.PostJson(visitorPath, visitorMap)
// 	if err != nil {
// 		log.Fatalln("login failed ,err", err)
// 	}
// 	fmt.Println("%#v", result)

// 	lr, ok := result.(map[string]interface{})
// 	if !ok {
// 		log.Fatalln("type assert failed")
// 	}
// 	token := lr["token"].(string)

// 	roomAddr := serverAddr + ":3000"

// 	cl, err := net.Dial("tcp", roomAddr)
// 	if err != nil {
// 		log.Fatalf("dial err [%s]", err.Error())
// 	}

// 	conn, err := websocket.NewClient(newConfig(roomAddr, path), cl)
// 	if err != nil {
// 		log.Fatalf("WebSocket handshake error: %v", err)
// 		return
// 	}

// 	//加入房间

// 	playerId := 0
// 	go func() {
// 		msg := &pb.Message{}
// 		t := int32(messagetypepb.MessageType_CGLoginType)
// 		msg.MessageType = &t
// 		var pId int64 = int64(playerId)
// 		token := token

// 		cgLogin := &loginpb.CGLogin{}
// 		cgLogin.PlayerId = &pId
// 		cgLogin.Token = &token
// 		err := proto.SetExtension(msg, loginpb.E_CgLogin, cgLogin)
// 		if err != nil {
// 			log.Printf("set extension error [%v]", err)
// 		}

// 		content, err := proto.Marshal(msg)
// 		if err != nil {
// 			log.Printf("marshal error [%v]", err)
// 		}

// 		log.Printf("%v", msg)
// 		err = websocket.Message.Send(conn, content)
// 		if err != nil {
// 			log.Printf("send err [%s]", err.Error())
// 			conn.Close()
// 		}
// 	}()

// 	initDispatch()

// 	tc := client.NewClient(int64(playerId), conn)
// 	for {
// 		var msg []byte
// 		err = websocket.Message.Receive(conn, &msg)
// 		if err != nil {
// 			log.Printf("receive error message [%s]", err.Error())
// 			break
// 		}
// 		var p *pb.Message = &pb.Message{}
// 		err := proto.Unmarshal(msg, p)
// 		if err != nil {
// 			log.Printf("unmarshall error [%s]", err.Error())
// 		}

// 		log.Printf("receive %#v", p)
// 		err = d.Handle(tc, p)
// 		if err != nil {
// 			log.Println("handler msg error ", err.Error())
// 		}
// 	}

// 	log.Println("client close")
// }

func newConfig(serverAddr string, path string) *websocket.Config {
	config, _ := websocket.NewConfig(fmt.Sprintf("ws://%s%s", serverAddr, path), fmt.Sprintf("http://%s%s", serverAddr, path))
	return config
}

func login(deviceMac string) (token string, err error) {
	visitorPath := "http://" + serverAddr + "/api/user/visitor"
	visitorMap := userapi.VisitorLoginForm{
		DeviceMac: deviceMac,
	}

	result, err := pkghttputils.PostJson(visitorPath, nil, visitorMap)
	if err != nil {
		return "", err
	}

	lr, ok := result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("type assert failed")
	}
	token = lr["token"].(string)
	return token, nil
}

func connectTargetServer(roomAddr string, token string) {

	cl, err := net.Dial("tcp", roomAddr)
	if err != nil {
		log.Error("dial err [%s]", err.Error())
		return
	}

	conn, err := websocket.NewClient(newConfig(roomAddr, path), cl)
	if err != nil {
		log.Error("WebSocket handshake error: %v", err)
		return
	}

	//加入房间

	playerId := 0
	if !timeout {
		go func() {
			msg := &pb.Message{}
			t := int32(messagetypepb.MessageType_CGLoginType)
			msg.MessageType = &t
			var pId int64 = int64(playerId)
			token := token

			cgLogin := &loginpb.CGLogin{}
			cgLogin.PlayerId = &pId
			cgLogin.Token = &token
			err := proto.SetExtension(msg, loginpb.E_CgLogin, cgLogin)
			if err != nil {
				log.Printf("set extension error [%v]", err)
			}

			content, err := proto.Marshal(msg)
			if err != nil {
				log.Printf("marshal error [%v]", err)
			}

			log.Printf("%v", msg)
			err = websocket.Message.Send(conn, content)
			if err != nil {
				log.Printf("send err [%s]", err.Error())
				conn.Close()
			}
		}()
	}

	d := clienthandler.NewDispatcher()
	roomhandler.InitDispatch(d)

	tc := client.NewClient(int64(playerId), conn)
	for {
		var msg []byte
		err = websocket.Message.Receive(conn, &msg)
		if err != nil {
			log.Printf("receive error message [%s]", err.Error())
			break
		}
		var p *pb.Message = &pb.Message{}
		err := proto.Unmarshal(msg, p)
		if err != nil {
			log.Printf("unmarshall error [%s]", err.Error())
		}

		log.Printf("receive %#v", p)
		err = d.Handle(tc, p)
		if err != nil {
			log.Println("handler msg error ", err.Error())
		}
	}

	log.Println("client close")
}

func connectChatServer(roomAddr string, token string) {

	cl, err := net.Dial("tcp", roomAddr)
	if err != nil {
		log.Fatalf("dial err [%s]", err.Error())
	}

	conn, err := websocket.NewClient(newConfig(roomAddr, "/chat"), cl)
	if err != nil {
		log.Fatalf("WebSocket handshake error: %v", err)
		return
	}
/****
	//加入房间
	go func() { 
		cgChatLogin := &chatmessage.CGChatLoginMesssage{}
		cgChatLogin.Token = token
		msg := &chatmessage.Message{}
		msg.Typ = int(chatmessage.CGChatLoginMesssageType)

		cgChatLoginBytes, err := json.Marshal(cgChatLogin)
		if err != nil {
			log.Println("marshal error", err.Error())
			conn.Close()
			return
		}
		msg.Body = cgChatLoginBytes
		content, err := json.Marshal(msg)
		if err != nil {
			log.Println("marshal error", err.Error())
			conn.Close()
			return
		}
		log.Printf("%v", msg)
		err = websocket.Message.Send(conn, content)
		if err != nil {
			log.Printf("send err [%s]", err.Error())
			conn.Close()
		}

		//定时发消息
		for {
			<-time.After(time.Second * 10)
			cgChat := &chatmessage.CGChatMessage{}

			cgChat.Msg = []byte(fmt.Sprintf("time %d", time.Now().UnixNano()/int64(time.Second)))
			fmt.Printf(string(cgChat.Msg))
			chatMsg := &chatmessage.Message{}
			chatMsg.Typ = int(chatmessage.CGChatMessageType)

			cgChatBytes, err := json.Marshal(cgChat)
			if err != nil {
				log.Println("marshal error", err.Error())
				conn.Close()
				return
			}
			chatMsg.Body = cgChatBytes
			chatContent, err := json.Marshal(chatMsg)
			if err != nil {
				log.Println("marshal error", err.Error())
				conn.Close()
				return
			}
			log.Printf("send %v", chatMsg)
			err = websocket.Message.Send(conn, chatContent)
			if err != nil {
				log.Printf("send err [%s]", err.Error())
				conn.Close()
				return
			}
		}
	}()
		***/
	for {
		var msg []byte
		err = websocket.Message.Receive(conn, &msg)
		if err != nil {
			log.Printf("receive error message [%s]", err.Error())
			break
		}
		log.Printf("receive %#v\n", string(msg))
	}

	log.Println("client close")
}
