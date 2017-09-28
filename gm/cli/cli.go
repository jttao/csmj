package cli

import (
	"encoding/json"
	"fmt"
	gamedb "game/db"
	"game/gm/api"

	gameredis "game/redis"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	gmservice "game/gm/service"

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
	apiPath = "/api"
)

func Start() {
	app := cli.NewApp()
	app.Name = "gm"
	app.Usage = "gm [global options]"

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
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	return nil
}

type serverConfig struct {
	Host  string                 `json:"host"`
	Port  int                    `json:"port"`
	Redis *gameredis.RedisConfig `json:"redis"`
	DB    *gamedb.DbConfig       `json:"db"`
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
	us := gmservice.NewUserService(db)
	rcs := gmservice.NewRoomCardService(db, rs)
	ts := gmservice.NewTaskService(db)
	
	router := mux.NewRouter()
	subrouter := router.PathPrefix(apiPath).Subrouter()
	api.Router(subrouter)
	n.Use(SetupUserServiceHandler(us))
	n.Use(SetupRoomCardServiceHandler(rcs))
	n.Use(SetupTaskServiceHandler(ts))

	n.UseHandler(router)
	//register interruput handler
	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	hooker := osutils.NewInterruptHooker()
	hooker.AddHandler(osutils.InterruptHandlerFunc(stop))
	log.Println("trying to listen ", addr)
	n.Run(addr)
}

func stop() {
	log.Println("server stop")
}

func SetupUserServiceHandler(us gmservice.UserService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := gmservice.WithUserService(ctx, us)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}

func SetupRoomCardServiceHandler(us gmservice.RoomCardService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := gmservice.WithRoomCardService(ctx, us)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}

func SetupTaskServiceHandler(us gmservice.TaskService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := gmservice.WithTaskService(ctx, us)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}