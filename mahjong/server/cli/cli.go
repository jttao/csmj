package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gamedb "game/db"
	"game/mahjong/changsha"
	"game/mahjong/server/mahjong"
	gameredis "game/redis"
	roommanageclient "game/roommanage/client"
	websocketsession "game/session/websocket"
	userservice "game/user/service"
	"net/http"

	"game/mahjong/server/login"
	roomhandler "game/mahjong/server/room/handler"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"golang.org/x/net/websocket"
)

var (
	debug      = false
	path       = "/room"
	configFile = "./config/config.json"
)

func Start() {
	app := cli.NewApp()
	app.Name = "room server"
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
			Name:        "config,c",
			Value:       configFile,
			Usage:       "config file",
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "path,p",
			Value:       path,
			Usage:       "path",
			Destination: &path,
		},
	}

	app.Action = start
	app.Run(os.Args)
}

// var sh session.SessionHandler
// var dispatcher handler.Dispatcher

type serverConfig struct {
	ServerId   int                                      `json:"serverId"`
	Host       string                                   `json:"host"`
	Port       int                                      `json:"port"`
	DB         *gamedb.DbConfig                         `json:db`
	Redis      *gameredis.RedisConfig                   `json:redis`
	User       *userservice.UserConfig                  `json:"user"`
	RoomManage *roommanageclient.RoomManageClientConfig `json:"roomManage"`
	Room       *changsha.RoomConfig                     `json:"room"`
}

func initConfig(configFile string) (sc *serverConfig, err error) {
	abs, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, err
	}
	sc = &serverConfig{}
	if err = json.Unmarshal(bs, sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func start(c *cli.Context) {
	//设置日志格式
	log.SetLevel(log.DebugLevel)

	sc, err := initConfig(configFile)
	if err != nil {
		log.Fatalln("init config file failed ", err)
	}

	//初始化db
	db, err := gamedb.NewDBService(sc.DB)
	if err != nil {
		log.Fatalln("init db service failed:", err)
	}

	//初始化redis
	rs, err := gameredis.NewRedisService(sc.Redis)
	if err != nil {
		log.Fatalln("init redis service failed:", err)
	}

	//初始化用户服务
	us, err := userservice.NewUserService(sc.User, db, rs)
	if err != nil {
		log.Fatalln("init user service failed:", err)
	}
	//初始化房间管理客户端
	rmc := roommanageclient.NewRoomManageClient(sc.RoomManage)
	//初始化用户管理
	pm := mahjong.NewPlayerManager()
	//房间管理
	rm := changsha.NewRoomManager()
	//初始化心跳goroutine
	mah := &mahjong.Mahjong{}
	mah.UserService = us
	mah.DB = db
	mah.RS = rs
	mah.RoomManageClient = rmc
	mah.RoomManager = rm
	mah.PlayerManager = pm
	dis := mahjong.NewDispatch()
	login.InitDispatcher(dis)
	roomDis := mahjong.NewDispatch()
	roomhandler.InitDispatcher(roomDis)

	hb := mahjong.NewHeartBeat(mah)
	hb.Start()

	gp := mahjong.NewGlobalProcessor(1000, dis)
	gp.Start()
	mah.GlobalProcessor = gp
	mah.Dispatcher = roomDis
	serverCfg := &mahjong.ServerConfig{}
	serverCfg.Room = sc.Room
	mah.ServerCfg = serverCfg

	sessionOpener := mahjong.AuthTimeoutMiddleware(mahjong.NewSessionOpener(mah))
	sessionCloser := mahjong.NewSessionCloser()
	sessionProcessor := mahjong.NewSessionProcessor()
	http.Handle("/room", websocket.Handler(websocketsession.NewWebsocketHandler(sessionOpener, sessionCloser, sessionProcessor, nil).Handle))
	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	log.Infof("listen %s", addr)
	http.ListenAndServe(addr, nil)

}
