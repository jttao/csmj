package player

import (
	"context"
	"game/session"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type PlayerState int

const (
	PlayerStateInit PlayerState = iota
	PlayerStateInRoom
)

type Player interface {
	Id() int64
	RoomId() int64
	PlayerState() PlayerState
	Session() session.Session
	Start()
	Stop()
	Send(msg []byte)
	Ping()
	LastPing() int64
}

type player struct {
	rwm         sync.RWMutex
	id          int64
	roomId      int64
	playerState PlayerState
	session     session.Session
	msgs        chan []byte
	done        chan struct{}
	lastPing    int64
}

func (p *player) Id() int64 {
	return p.id
}

func (p *player) RoomId() int64 {
	return p.roomId
}

func (p *player) PlayerState() PlayerState {
	return p.playerState
}

func (p *player) Session() session.Session {
	return p.session
}

func (p *player) SetSession(s session.Session) {
	p.session = s 
}

func (p *player) Ping() {
	p.rwm.Lock()
	defer p.rwm.Unlock()
	p.lastPing = time.Now().UnixNano() / int64(time.Millisecond)
}

func (p *player) LastPing() int64 {
	p.rwm.RLock()
	defer p.rwm.RUnlock()
	return p.lastPing
}

func (p *player) Start() {
	go func() {
		log.WithFields(
			log.Fields{
				"玩家id": p.id,
			}).Info("玩家,开始")
	Loop:
		for {
			select {
			case <-p.done:
				break Loop
			case msg, ok := <-p.msgs:
				{
					if !ok {
						break
					}
					err := p.session.Send(msg)
					if err != nil {
						log.WithFields(
							log.Fields{
								"玩家id": p.id,
							}).Error("玩家发送消息,错误")
					}
				}
			}
		}
		log.WithFields(
			log.Fields{
				"玩家id": p.id,
			}).Info("玩家,结束")
	}()
}

func (p *player) Send(msg []byte) {
	p.msgs <- msg
}

func (p *player) Stop() {
	p.done <- struct{}{}
}

func NewPlayer(id int64, roomId int64, s session.Session) Player {
	p := &player{
		id:          id,
		roomId:      roomId,
		playerState: PlayerStateInit,
		session:     s,
	}
	p.msgs = make(chan []byte, 1000)
	p.done = make(chan struct{})
	p.Ping()
	return p
}

const (
	playerKey = "game.mahjong.player"
)

func PlayerInContext(ctx context.Context) Player {
	p, ok := ctx.Value(playerKey).(Player)
	if !ok {
		return nil
	}
	return p
}

func WithPlayer(ctx context.Context, pl Player) context.Context {
	return context.WithValue(ctx, playerKey, pl)
}
