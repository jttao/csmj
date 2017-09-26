package mahjong

import (
	"context"
	"fmt"
	gamedb "game/db"
	"game/mahjong/changsha"
	"game/mahjong/server/player"
	gameredis "game/redis"
	roommanageclient "game/roommanage/client"
	hallclient "game/hall/client"
	
	userserice "game/user/service" 
	"sync"
)

type ServerConfig struct {
	Room *changsha.RoomConfig `json:"room"`
}

type Mahjong struct {
	ServerCfg        *ServerConfig
	DB               gamedb.DBService
	RS               gameredis.RedisService
	UserService      userserice.UserService 
	RoomManageClient roommanageclient.RoomManageClient
	HallClient		 hallclient.HallClient
	RoomManager      *changsha.RoomManager
	GlobalProcessor  *GlobalProcessor
	Dispatcher       *Dispatcher
	PlayerManager    *PlayerManager
}

const (
	mahjongKey = "game.mahjong"
)

func WithMahjong(ctx context.Context, m *Mahjong) context.Context {
	return context.WithValue(ctx, mahjongKey, m)
}

func MahjongInContext(ctx context.Context) *Mahjong {
	m, ok := ctx.Value(mahjongKey).(*Mahjong)
	if !ok {
		return nil
	}
	return m
}

type PlayerManager struct {
	rwm     sync.RWMutex
	players map[int64]player.Player
}

func (pm *PlayerManager) AddPlayer(p player.Player) error {
	pm.rwm.Lock()
	defer pm.rwm.Unlock()
	_, exist := pm.players[p.Id()]
	if exist {
		return fmt.Errorf("玩家已经存在")
	}
	pm.players[p.Id()] = p
	return nil
}

func (pm *PlayerManager) RemovePlayer(p player.Player) {
	pm.rwm.Lock()
	defer pm.rwm.Unlock()
	delete(pm.players, p.Id())
}

func (pm *PlayerManager) Players() map[int64]player.Player {
	pm.rwm.RLock()
	defer pm.rwm.RUnlock()
	return pm.players
}

func (pm *PlayerManager) GetPlayerById(id int64) player.Player {
	pm.rwm.RLock()
	defer pm.rwm.RUnlock()
	pl, exist := pm.players[id]
	if !exist {
		return nil
	}
	return pl
}

func NewPlayerManager() *PlayerManager {
	pm := &PlayerManager{}
	pm.players = make(map[int64]player.Player)
	return pm
}
