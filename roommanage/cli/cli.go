package cli

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"io/ioutil"
	"path/filepath"

	gamedb "game/db"
	gameredis "game/redis"
	roommanageapi "game/roommanage/api"
	roommanageservice "game/roommanage/service"
	userservice "game/user/service"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/xozrc/pkg/osutils"
)

var (
	debug      = false
	configFile = "./config/config.json"
)

func Start() {

	app := cli.NewApp()
	app.Name = "room manage"
	app.Usage = "roommanage [global options]"

	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "debug ",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "config,c",
			Value:       configFile,
			Usage:       "config file",
			Destination: &configFile,
		},
	}
	app.Before = before
	app.Action = start
	app.Run(os.Args)
}

func before(ctx *cli.Context) error {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return nil
}

type serverConfig struct {
	Host       string                              `json:"host"`
	Port       int                                 `json:"port"`
	Redis      *gameredis.RedisConfig              `json:"redis"`
	DB         *gamedb.DbConfig                    `json:"db"`
	User       *userservice.UserConfig             `json:"user"`
	RoomManage *roommanageservice.RoomManageConfig `json:"roommanage"`
}

func newServerConfig(configFile string) (sc *serverConfig, err error) {
	c, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	sc = &serverConfig{}
	err = json.Unmarshal(c, sc)
	if err != nil {
		return nil, err
	}
	return
}

var (
	apiPath = "/api"
)

func start(ctx *cli.Context) {
	log.SetLevel(log.DebugLevel)
	config, err := filepath.Abs(configFile)
	if err != nil {
		log.Fatalln("filepath abs failed:", err)
	}

	sc, err := newServerConfig(config)
	if err != nil {
		log.Fatalln("read config file failed:", err)
	}
	n := negroni.Classic()
	//初始化db
	dbService, err := gamedb.NewDBService(sc.DB)
	if err != nil {
		log.Fatalln("init db service failed:", err)
	}
	redisService, err := gameredis.NewRedisService(sc.Redis)
	if err != nil {
		log.Fatalln("init redis service failed:", err)
	}

	us, err := userservice.NewUserService(sc.User, dbService, redisService)

	if err != nil {
		log.Fatalln("init user service failed:", err)
	}

	csrs := roommanageservice.NewRoomManageService(sc.RoomManage, dbService, redisService)
	err = rms.Start() 
	if err != nil {
		log.Fatalln("init room mananage failed:", err)
	}
	
	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	router := mux.NewRouter()
	subrouter := router.PathPrefix(apiPath).Subrouter()
	roommanageapi.Router(subrouter)

	n.UseFunc(setupUserServiceHandler(us))
	n.UseFunc(setupChangshaRoomServiceHandler(csrs))
	n.UseHandler(router)
	//register interruput handler
	hooker := osutils.NewInterruptHooker()
	hooker.AddHandler(osutils.InterruptHandlerFunc(stop))
	log.Println("trying to listen ", addr)
	n.Run(addr)
}

func stop() {
	log.Println("stop server")
}

//设置房间管理服务
func setupChangshaRoomServiceHandler(csrs roommanageservice.RoomManageService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := roommanageservice.WithRoomManageService(ctx, csrs)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}

//设置用户服务
func setupUserServiceHandler(us userservice.UserService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := userservice.WithUserService(ctx, us)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}
