package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	gamedb "game/db"
	gameredis "game/redis"
	"math/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

type RoomManageService interface {
	AutoRoom(roomType RoomType, pId int64, maxPlayers int, round int, cost int, roomConfig string,forbidIp int, ip string,openRoomType OpenRoomType) (Room, error)
	CreateRoom(roomType RoomType, pId int64, maxPlayers int, round int, cost int, roomConfig string,forbidIp int, ip string,openRoomType OpenRoomType) (r Room, err error)
	JoinRoom(rid int64, pId int64,ip string) Room
	LeaveRoom(rId int64, pId int64)
	DestroyRoom(roomId int64)
	GetRoomById(roomId int64) Room
	GetPlayerById(playerId int64) Player
	AddPlayer(p Player) bool
	RemovePlayer(p Player)
	Rooms() []Room
	Debug() bool
	IsFull() bool
	IfFree() (bool, error)
	IfCheck() (bool, error)
	IfForbidIp(id int64,ip string) bool
	GetServerByServerId(serverId string) *ServerConfig
	
	ListAgentRooms(agentId int64) []Room
	AgentCloseRooms(agentId int64) (err error)
	AgentCloseSingleRoom(agentId int64, roomId int64) (err error)

}

const (
	roomManageKey = "roomManage"
)

var (
	PlayerAlreadyInRoomError = fmt.Errorf("player already in room")
)

func RoomManageInContext(ctx context.Context) RoomManageService {
	c, ok := ctx.Value(roomManageKey).(RoomManageService)
	if !ok {
		return nil
	}
	return c
}

func WithRoomManageService(ctx context.Context, csrs RoomManageService) context.Context {
	return context.WithValue(ctx, roomManageKey, csrs)
}

type ServerConfig struct {
	Id   string `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Room int    `json:"room"`
}

type RoomManageConfig struct {
	Servers []*ServerConfig `json:"servers"`
	Debug   bool            `json:"debug"`
	ExpiredTime   int64 `json:"expiredTime"`
	UserRoomTime  int64 `json:"userRoomTime"`
	AgentRoomTime int64 `json:"agentRoomTime"`
	MaxRoomTime   int64 `json:"maxRoomTime"`
}

const (
	gameConfigFreeKey = "game.config.free"
)

const (
	gameConfigCheckKey = "game.config.check"
)

type roomManageService struct {
	config *RoomManageConfig
	db     gamedb.DBService
	rs     gameredis.RedisService
	rwm    sync.RWMutex
	//房间列表
	roomsMap map[int64]Room
	// //玩家列表
	playersMap map[int64]Player
	playersRwm sync.RWMutex
	//服务器列表
	servers map[string]int
	//房间号
	roomNums []int32
	//回收的房间号
	recycleRoomNums []int32

	done            chan struct{}
}

func (rms *roomManageService) Debug() bool {
	return rms.config.Debug
}

func (rms *roomManageService) IfFree() (bool, error) {
	pool := rms.rs.Pool()
	conn := pool.Get()
	if conn.Err() != nil {
		return false, conn.Err()
	}
	defer conn.Close()

	flag, err := redis.Bool(conn.Do("get", gameConfigFreeKey))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	return flag, nil
}

func (rms *roomManageService) IfCheck() (bool, error) {
	pool := rms.rs.Pool()
	conn := pool.Get()
	if conn.Err() != nil {
		return false, conn.Err()
	}
	defer conn.Close()

	flag, err := redis.Bool(conn.Do("get", gameConfigCheckKey))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	return flag, nil
}

func (rms *roomManageService) IfForbidIp(rid int64,ip string) bool {
	
	rms.rwm.Lock()
	defer rms.rwm.Unlock()
	r, ok := rms.roomsMap[rid]
	if !ok {
		return false
	}   

	flag := r.IfForbidIp(ip)

	if flag { 
		return true 
	}

	return false
}

func (rms *roomManageService) Start() error {
	go func() {
		log.Info("房间服务开始")
	Loop:
		for {
			select {
			case <-time.After(time.Minute * 5):
				{
					rms.tick()
				}
			case <-rms.done:
				{
					break Loop
				}
			}
		}
		log.Info("房间服务结束")
	}()
	return nil
}

func (rms *roomManageService) tick() {
	rms.rwm.Lock()
	defer rms.rwm.Unlock()
	log.Info("房间服务心跳") 
	for _, r := range rms.roomsMap {
		if r.IfShouldRemove() {
			rms.DestroyRoom(r.Id())
		}
	} 
}

func (rms *roomManageService) Stop() {
	log.Info("房间服务正在结束")
	rms.done <- struct{}{}
}

func (rms *roomManageService) AutoRoom(roomType RoomType, pId int64, maxPlayers int, round int, cost int, roomConfig string,forbidIp int,ip string,openRoomType OpenRoomType) (Room, error) {
	rms.rwm.Lock()
	defer rms.rwm.Unlock()

	owner := NewPlayer(pId, int64(0),ip)

	for _, r := range rms.roomsMap {
		if r.Full() {
			continue
		}
		flag := r.JoinPlayer(owner)
		if !flag {
			continue
		}
		rms.playersMap[owner.Id()] = owner
		return r, nil
	}
	//创建 新的
	// rm := &roommanagemodel.RoomModel{}
	// rm.OwnerId = owner.Id()
	// rm.RoomType = int(roomType)
	// rm.RoomConfig = roomConfig
	// rm.Round = round
	// rm.Cost = cost
	// tdb := rms.db.DB().Save(rm)
	// if tdb.Error != nil {
	// 	return nil, tdb.Error
	// }

	roomId, err := rms.getUsableRoomNum()
	if err != nil {
		return nil, err
	}
	//查找服务器
	bestServer := rms.findBestServer()
	if bestServer == "" {
		return nil, nil
	}
	
	roomTime := rms.config.AgentRoomTime
	maxRoomTime := rms.config.MaxRoomTime   

	if openRoomType==OpenRoomTypeAgent {
		roomTime = rms.config.AgentRoomTime
		maxRoomTime = rms.config.MaxRoomTime   
	}
	
	r := NewRoom(roomTime,maxRoomTime,int64(roomId), bestServer, roomType, pId, cost, maxPlayers, round, roomConfig , forbidIp,openRoomType)

	if openRoomType==OpenRoomTypeUser { 
		r.JoinPlayer(owner)
	}
	
	//缓存
	rms.roomsMap[r.Id()] = r

	rms.playersMap[owner.Id()] = owner
	num, ok := rms.servers[r.ServerId()]
	if !ok {
		rms.servers[r.ServerId()] = 1
	} else {
		rms.servers[r.ServerId()] = num + 1
	}
	return r, nil
}

//创建房间
func (rms *roomManageService) CreateRoom(roomType RoomType, ownerId int64, maxPlayers int, round int, cost int, roomConfig string,forbidIp int,ip string,openRoomType OpenRoomType) (r Room, err error) {

	rms.rwm.Lock()
	defer rms.rwm.Unlock()

	//查找服务器
	bestServer := rms.findBestServer()
	if bestServer == "" {
		return nil, nil
	}
	
	owner := NewPlayer(ownerId, int64(0) , ip)
	// //创建房间
	// rm := &roommanagemodel.RoomModel{}
	// rm.OwnerId = owner.Id()
	// rm.RoomType = int(roomType)
	// rm.RoomConfig = roomConfig
	// rm.Round = round
	// rm.Cost = cost
	// tdb := rms.db.DB().Save(rm)
	// if tdb.Error != nil {
	// 	return nil, tdb.Error
	// }
	//获取房间号
	roomId, err := rms.getUsableRoomNum()
	if err != nil {
		return nil, err
	}

	roomTime := rms.config.AgentRoomTime
	maxRoomTime := rms.config.MaxRoomTime   
	
	if openRoomType==OpenRoomTypeAgent {
		roomTime = rms.config.AgentRoomTime
		maxRoomTime = rms.config.MaxRoomTime   
	}
	
	r = NewRoom(roomTime,maxRoomTime,int64(roomId), bestServer, roomType, ownerId, cost, maxPlayers, round, roomConfig , forbidIp,openRoomType)
	
	if openRoomType==OpenRoomTypeUser { 
		r.JoinPlayer(owner)
	}
	
	//缓存
	rms.roomsMap[r.Id()] = r
	rms.playersMap[owner.Id()] = owner

	num, ok := rms.servers[r.ServerId()]
	if !ok {
		rms.servers[r.ServerId()] = 1
	} else {
		rms.servers[r.ServerId()] = num + 1
	}
	
	return r, nil
}

//加入房间
func (rms *roomManageService) JoinRoom(rid int64, pId int64,ip string) Room {
	rms.rwm.Lock()
	defer rms.rwm.Unlock()
	r, ok := rms.roomsMap[rid]
	if !ok {
		return nil
	}

	p := NewPlayer(pId, int64(0),ip)  

	flag := r.JoinPlayer(p)

	if flag {
		rms.playersMap[p.Id()] = p
		return r
	} else {
		return nil
	}

}

//离开房间
func (rms *roomManageService) LeaveRoom(rid int64, pId int64) {
	rms.rwm.Lock()
	defer rms.rwm.Unlock()
	delete(rms.playersMap, pId)
	r, ok := rms.roomsMap[rid]
	if !ok {
		return
	}
	r.LeavePlayer(pId)
	return
}

//删除房间
func (rms *roomManageService) DestroyRoom(roomId int64) {
	r := rms.GetRoomById(roomId)
	if r == nil {
		log.WithField("房间号", roomId).Warn("摧毁不存在的房间")
		return
	}

	rms.rwm.Lock()
	defer rms.rwm.Unlock()
	//移除所有玩家
	for _, p := range r.Players() {
		delete(rms.playersMap, p.Id())
	}
	delete(rms.roomsMap, roomId)
	rms.recycleRoomNums = append(rms.recycleRoomNums, int32(roomId))
	num, ok := rms.servers[r.ServerId()]
	if !ok {
		return
	}
	rms.servers[r.ServerId()] = num - 1
	return
}


//获取房间
func (rms *roomManageService) GetRoomById(roomId int64) Room {
	rms.rwm.RLock()
	defer rms.rwm.RUnlock()
	r, ok := rms.roomsMap[roomId]
	if !ok {
		return nil
	}
	return r
}

func (rms *roomManageService) GetPlayerById(playerId int64) Player {
	rms.playersRwm.RLock()
	defer rms.playersRwm.RUnlock()
	p, ok := rms.playersMap[playerId]
	if !ok {
		return nil
	}
	return p
}

func (rms *roomManageService) AddPlayer(p Player) bool {
	rms.playersRwm.Lock()
	defer rms.playersRwm.Unlock()
	_, exist := rms.playersMap[p.Id()]
	if exist {
		return false
	}
	rms.playersMap[p.Id()] = p
	return true
}

func (rms *roomManageService) RemovePlayer(p Player) {
	rms.playersRwm.Lock()
	defer rms.playersRwm.Unlock()
	delete(rms.playersMap, p.Id())
}

func (rms *roomManageService) Rooms() []Room {
	trs := make([]Room, 0, len(rms.roomsMap))
	for _, tr := range rms.roomsMap {
		trs = append(trs, tr)
	}
	return trs
}

func (rms *roomManageService) IsFull() bool {
	rms.rwm.RLock()
	defer rms.rwm.RUnlock()
	for _, sc := range rms.config.Servers {
		num, ok := rms.servers[sc.Id]
		if !ok {
			return false
		}
		if num < sc.Room {
			return false
		}
	}
	return true
}

func (rms *roomManageService) GetServerByServerId(serverId string) *ServerConfig {
	rms.rwm.RLock()
	defer rms.rwm.RUnlock()
	for _, s := range rms.config.Servers {
		if s.Id == serverId {
			return s
		}
	}
	return nil
}

func (rms *roomManageService) findBestServer() string {
	for _, sc := range rms.config.Servers {
		num, ok := rms.servers[sc.Id]
		if !ok {
			return sc.Id
		}
		if num < sc.Room {
			return sc.Id
		}
	}
	return ""
}

//洗牌
func (rms *roomManageService) initRooms() {
	rms.roomNums = make([]int32, 0, 899999)

	for i := 100000; i < 1000000; i++ {
		rms.roomNums = append(rms.roomNums, int32(i))
	}
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(rms.roomNums))
	for i, v := range perm {
		rms.roomNums[i], rms.roomNums[v] = rms.roomNums[v], rms.roomNums[i]
	}
}

//获取可用的房间号
func (rms *roomManageService) getUsableRoomNum() (int32, error) {
	if len(rms.roomNums) == 0 {
		rms.roomNums = append(rms.roomNums, rms.recycleRoomNums...)
		rms.recycleRoomNums = nil
	}
	if len(rms.roomNums) == 0 {
		return 0, errors.New("no usable room number")
	}
	roomId := rms.roomNums[0]
	rms.roomNums = rms.roomNums[1:]
	return roomId, nil
}






//---------


const (
	maxAgentRooms = 20
)


func (rms *roomManageService) AgentCloseRooms(agentId int64) (err error) {
	rms.rwm.Lock()
	defer rms.rwm.Unlock()
	rooms := rms.listAgentRooms(agentId)
	for _, room := range rooms {
		if room.IfCanRemove() {
			rms.DestroyRoom(room.Id())
		}
	}
	return
}

func (rms *roomManageService) AgentCloseSingleRoom(agentId int64, roomId int64) (err error) {
	rms.rwm.Lock()
	defer rms.rwm.Unlock()  

	r, ok := rms.roomsMap[roomId]

	if !ok { 
		return
	}

	if r == nil {
		return
	}
	
	if r.OwnerId() != agentId || r.getOpenRoomType() != OpenRoomTypeAgent {
		//err = ErrorAgentCloseRoomNoSelf
		return 
	}

	if !r.IfCanRemove() {
		//err = ErrorAgentCloseRoomInProcess
		return
	} 
	
	rms.DestroyRoom(roomId)
	
	return
}

//列出代理的房间
func (rms *roomManageService) ListAgentRooms(agentId int64) []Room {
	rms.rwm.Lock()
	defer rms.rwm.Unlock()
	return rms.listAgentRooms(agentId)
}

//列出代理的房间
func (rms *roomManageService) listAgentRooms(agentId int64) (rooms []Room) {
	for _, r := range rms.roomsMap {
		if r.getOpenRoomType() == OpenRoomTypeAgent {
			if r.OwnerId() == agentId {
				rooms = append(rooms, r)
			}
		}
	}	
	return
}

//---------
func NewRoomManageService(config *RoomManageConfig, db gamedb.DBService, rs gameredis.RedisService) RoomManageService {
	rms := &roomManageService{}
	rms.db = db
	rms.config = config
	rms.rs = rs
	rms.roomsMap = make(map[int64]Room)
	rms.playersMap = make(map[int64]Player)
	rms.servers = make(map[string]int)
	rms.initRooms()
	rms.done = make(chan struct{})
	return rms
}
