package service

import (
	"sync"
	"time"
	timeutils "game/pkg/timeutils"
)

type RoomType int

const (
	RoomTypeChangSha RoomType = iota
)

type OpenRoomType int

const (
	OpenRoomTypeUser OpenRoomType = iota
	OpenRoomTypeAgent
)

func (rt RoomType) Valid() bool {
	switch rt {
	case RoomTypeChangSha:
	default:
		return false
	}
	return true
}

type Room interface {
	Id() int64
	ServerId() string
	OwnerId() int64
	RoomType() RoomType
	Round() int
	RoomConfig() string
	MaxPlayers() int
	JoinPlayer(player Player) bool
	LeavePlayer(id int64)
	GetPlayerById(id int64) Player
	Players() map[int64]Player
	Full() bool
	Cost() int 
	ForbidIp() int
	IfForbidIp(ip string) bool
	GetOpenRoomType() OpenRoomType
	CanJoinPlayer() bool
	IfCanRemove() bool
	IfShouldRemove() bool
	IfExpiredJoin() bool
	IfExpired() bool
	CreateTime() int64
	ForbidJoinTime() int64 
	LastGameTime() int64 
} 

type room struct {
	rwm        sync.RWMutex
	id         int64
	serverId   string
	ownerId    int64
	maxPlayers int
	roomType   RoomType
	round      int
	roomConfig string
	cost       int
	players    map[int64]Player
	forbidIp   int
	openRoomType OpenRoomType 
	createTime int64 
	forbidJoinTime int64 
	lastGameTime int64 
}

func (r *room) Id() int64 {
	return r.id
}

func (r *room) ServerId() string {
	return r.serverId
}
func (r *room) OwnerId() int64 {
	return r.ownerId
}

func (r *room) RoomType() RoomType {
	return r.roomType
}

func (r *room) Round() int {
	return r.round
}

func (r *room) Cost() int {
	return r.cost
}

func (r *room) RoomConfig() string {
	return r.roomConfig
}

func (r *room) MaxPlayers() int {
	return r.maxPlayers
}

func (r *room) Players() map[int64]Player {
	return r.players
}

func (r *room) Full() bool {
	return len(r.players) >= int(r.maxPlayers) 
}

func (r *room) ForbidIp() int {
	return r.forbidIp
}

func (r *room) GetOpenRoomType() OpenRoomType {
	return r.openRoomType
}

func (r *room) CreateTime() int64 {
	return r.createTime
}

func (r *room) ForbidJoinTime() int64 {
	return r.forbidJoinTime
}

func (r *room) LastGameTime() int64 {
	return r.lastGameTime
}

func (r *room) IfForbidIp(ip string) bool  {
	if r.ForbidIp()==1 { 
		//检查是否相同IP禁止
		for _, user := range r.players {
			if (user.Ip()==ip) {
				return true
			}
		}   
	}
	return false 
} 

func (r *room) JoinPlayer(player Player) bool {

	tp := r.GetPlayerById(player.Id())
	if tp != nil { 
		return false
	}
	r.rwm.Lock()
	defer r.rwm.Unlock()
	
	//检查是否人数已满
	if len(r.players) >= r.maxPlayers { 
		return false
	}
	
	player.SetRoomId(r.Id())
	r.players[player.Id()] = player
	return true
}

func (r *room) LeavePlayer(id int64) {
	tp := r.GetPlayerById(id)
	if tp == nil {
		return
	}
	r.rwm.Lock()
	defer r.rwm.Unlock()
	delete(r.players, id)
	return
}

func (r *room) GetPlayerById(id int64) Player {
	r.rwm.RLock()
	defer r.rwm.RUnlock()
	p, ok := r.players[id]
	if !ok {
		return nil
	}
	return p
}

//是否可以加玩家
func (r *room) CanJoinPlayer() bool {
	now := timeutils.TimeToMillisecond(time.Now())
	if now >= r.forbidJoinTime {
		return false
	}
	return true
}

//是否可以可以移除了
func (r *room) IfCanRemove() bool {
	if len(r.players) == 0 {
		return true
	}
	if r.IfExpired() {
		return true
	}
	return false
}

//是否应该移除
func (r *room) IfShouldRemove() bool {
	//超出最大时间
	if r.IfExpired() {
		return true
	}
	//还没开始
	if len(r.players) < int(r.maxPlayers) && r.IfExpiredJoin() {
		return true
	}
	return false
}

//是否过期加入了
func (r *room) IfExpiredJoin() bool {
	now := timeutils.TimeToMillisecond(time.Now())
	if now >= r.forbidJoinTime {
		return true
	}
	return false
}

//是否超过时间了
func (r *room) IfExpired() bool {
	now := timeutils.TimeToMillisecond(time.Now())
	if now >= r.lastGameTime {
		return true
	}
	return false
}

func NewRoom(roomTime int64, maxRoomTime int64, id int64, serverId string, roomType RoomType, ownerId int64, cost int, maxPlayers int, round int, roomConfig string,forbidIp int,openRoomType OpenRoomType) Room {
	r := &room{}
	r.players = make(map[int64]Player)
	r.id = id
	r.serverId = serverId
	r.ownerId = ownerId
	r.roomType = roomType
	r.maxPlayers = maxPlayers
	r.round = round
	r.roomConfig = roomConfig
	r.cost = cost
	r.forbidIp = forbidIp 
	r.openRoomType = openRoomType   
	r.createTime = time.Now().UnixNano() / int64(time.Millisecond)
	r.forbidJoinTime = r.createTime + roomTime*int64(time.Second/time.Millisecond) 
	r.lastGameTime = r.createTime + maxRoomTime*int64(time.Second/time.Millisecond)
	
	return r
}
