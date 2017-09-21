package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	gamedb "game/db"
	changsharoomapi "game/mahjong/changsha/room/api"

	changsharoomservice "game/mahjong/changsha/room/service"
	changshatemplateapi "game/mahjong/changsha/template/api"
	changshatemplateservice "game/mahjong/changsha/template/service"
	gameredis "game/redis"
	userservice "game/user/service"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/xozrc/pkg/osutils"
)

var (
	debug      = false
	configFile = "./config/config.json"
)

const (
	apiPath = "/api/changsha"
)

func Start() {
	app := cli.NewApp()
	app.Name = "changsha"
	app.Usage = "changsha [global options]"

	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Value:       configFile,
			Usage:       "config file",
			Destination: &configFile,
		},
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "debug",
			Destination: &debug,
		},
	}
	app.Before = before
	app.Action = start
	app.Run(os.Args)
}

func before(ctx *cli.Context) error {
	return nil
}

type serverConfig struct {
	Host     string                  `json:"host"`
	Port     int                     `json:"port"`
	ChangSha *changShaConfig         `json:"changsha"`
	User     *userservice.UserConfig `json:"user"`
	Redis    *gameredis.RedisConfig  `json:"redis"`
	DB       *gamedb.DbConfig        `json:"db"`
}

type changShaConfig struct {
	Template string                                  `json:"template"`
	Room     *changsharoomservice.ChangShaRoomConfig `json:"room"`
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
	db, err := gamedb.NewDBService(sc.DB)
	if err != nil {
		log.Fatalln("init db service failed:", err)
	}
	rs, err := gameredis.NewRedisService(sc.Redis)
	if err != nil {
		log.Fatalln("init redis service failed:", err)
	}

	changshaConfigFile, err := filepath.Abs(sc.ChangSha.Template)
	if err != nil {
		log.Fatalln("filepath abs failed:", err)
	}

	changshaConfigBody, err := ioutil.ReadFile(changshaConfigFile)
	if err != nil {
		log.Fatalln("read changsha config  failed:", err)
	}
	changShaConfig := &changshatemplateservice.ChangShaRoomTemplateConfig{}
	err = json.Unmarshal(changshaConfigBody, changShaConfig)
	if err != nil {
		log.Fatalln("unmarshal changsha config  failed:", err)
	}

	csrts := changshatemplateservice.NewChangShaRoomTemplateService(changShaConfig)

	if err != nil {
		log.Fatalln("init changsha template service failed:", err)
	}

	//初始化长沙房间服务
	csrs := changsharoomservice.NewChangShaRoomService(sc.ChangSha.Room, csrts)

	//初始化用户服务
	us, err := userservice.NewUserService(sc.User, db, rs)
	if err != nil {
		log.Fatalln("init changsha user service failed:", err)
	}

	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	router := mux.NewRouter().StrictSlash(true)

	subrouter := router.PathPrefix(apiPath).Subrouter()
	changshatemplateapi.Router(subrouter)
	changsharoomapi.Router(subrouter)

	n.UseFunc(userservice.SetupUserServiceHandler(us))
	n.UseFunc(setupChangShaTemplateServiceHandler(csrts))
	n.UseFunc(setupChangShaRoomServiceHandler(csrs))
	//认证服务
	n.UseFunc(userservice.AuthHandler())
	n.UseHandler(router)
	//register interruput handler
	hooker := osutils.NewInterruptHooker()
	hooker.AddHandler(osutils.InterruptHandlerFunc(stop))
	log.Println("trying to listen ", addr)
	n.Run(addr)
}

func stop() {

}

func setupChangShaTemplateServiceHandler(csrts changshatemplateservice.ChangShaRoomTemplateService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := changshatemplateservice.WithChangShaTemplateService(ctx, csrts)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}

func setupChangShaRoomServiceHandler(csrs changsharoomservice.ChangShaRoomService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := changsharoomservice.WithChangShaRoomService(ctx, csrs)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}
