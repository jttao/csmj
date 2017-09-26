package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	gamedb "game/db"
	hallapi "game/hall/api"
	news "game/hall/news"
	notice "game/hall/notice"
	tasks "game/hall/tasks"

	changsharoomapi "game/mahjong/changsha/room/api"
	changsharoomservice "game/mahjong/changsha/room/service"
	changshatemplateapi "game/mahjong/changsha/template/api"
	changshatemplateservice "game/mahjong/changsha/template/service"
	recordapi "game/mahjong/record/api"
	recordservice "game/mahjong/record/service"
	gameredis "game/redis"
	roommanageclient "game/roommanage/client"
	userservice "game/user/service"
	hallclient "game/hall/client"

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
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return nil
}

type serverConfig struct {
	Host       string                                   `json:"host"`
	Port       int                                      `json:"port"`
	ChangSha   *changShaConfig                          `json:"changsha"`
	RoomManage *roommanageclient.RoomManageClientConfig `json:"roomManage"` 
	User       *userservice.UserConfig                  `json:"user"`
	Redis      *gameredis.RedisConfig                   `json:"redis"`
	DB         *gamedb.DbConfig                         `json:"db"`
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

	//消息服务
	ns := news.NewNewsService(db, rs)
	//跑马灯服务
	noticeService := notice.NewNoticeService(db, rs)
	//任务服务器
	ts := tasks.NewTaskService(db, rs)

	//初始化房间查询服务
	rmc := roommanageclient.NewRoomManageClient(sc.RoomManage)
	
	//初始化长沙房间服务
	csrs := changsharoomservice.NewChangShaRoomService(sc.ChangSha.Room, csrts)

	//初始化用户服务
	us, err := userservice.NewUserService(sc.User, db, rs)
	if err != nil {
		log.Fatalln("init changsha user service failed:", err)
	}

	//大厅服务
	hc := hallclient.NewHallClient(ts,us) 
	
	//初始化录像服务
	recordS := recordservice.NewRecordService(db)

	addr := fmt.Sprintf("%s:%d", sc.Host, sc.Port)
	router := mux.NewRouter().StrictSlash(true)

	apiSubRouter := router.PathPrefix(apiPath).Subrouter()
	hallapi.Router(apiSubRouter)

	changShaPath := apiPath + "/changsha"
	subrouter := router.PathPrefix(changShaPath).Subrouter()

	changshatemplateapi.Router(subrouter)
	changsharoomapi.Router(subrouter)

	//录像api
	recordapi.Router(apiSubRouter)
	n.UseFunc(setupRecordServiceHandler(recordS))
	n.UseFunc(setupNewsServiceHandler(ns))
	n.UseFunc(setupNoticeServiceHandler(noticeService))
	n.UseFunc(setupRoomManageClientHandler(rmc))
	n.UseFunc(setupTasksServiceHandler(ts))
	
	n.UseFunc(setupHallClientHandler(hc))
	

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

func setupNoticeServiceHandler(ns notice.NoticeService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := notice.WithNoticeService(ctx, ns)
		nreq := req.WithContext(nctx)

		hf(rw, nreq)
	})
}

func setupNewsServiceHandler(ns news.NewsService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := news.WithNewsService(ctx, ns)
		nreq := req.WithContext(nctx)

		hf(rw, nreq)
	})
}

func setupTasksServiceHandler(ns tasks.TaskService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := tasks.WithTaskService(ctx, ns)
		nreq := req.WithContext(nctx)

		hf(rw, nreq)
	})
}

func setupHallClientHandler(hc hallclient.HallClient) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := hallclient.WithHallClient(ctx, hc)
		nreq := req.WithContext(nctx)
		
		hf(rw, nreq)
	})
}

func setupRoomManageClientHandler(rmc roommanageclient.RoomManageClient) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := roommanageclient.WithRoomManageClient(ctx, rmc)
		nreq := req.WithContext(nctx)

		hf(rw, nreq)
	})
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

func setupRecordServiceHandler(rs recordservice.RecordService) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, req *http.Request, hf http.HandlerFunc) {
		ctx := req.Context()
		nctx := recordservice.WithRecordService(ctx, rs)
		nreq := req.WithContext(nctx)
		hf(rw, nreq)
	})
}
