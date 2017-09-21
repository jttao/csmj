package mahjong

import (
	"context"
	"game/mahjong/changsha"
	"time"

	log "github.com/Sirupsen/logrus"
)

type RoomProcessor struct {
	r          *changsha.Room
	messages   chan *Message
	tick       time.Duration
	dispatcher *Dispatcher
}

func (p *RoomProcessor) Receive(msg *Message) {
	p.messages <- msg
}

func (p *RoomProcessor) Start() error {
	go func() {
		log.WithFields(
			log.Fields{
				"房间id": p.r.RoomId(),
			}).Info("房间处理器开始")
	Loop:
		for {
			select {
			case msg, flag := <-p.messages:
				{
					if !flag {
						break Loop
					}
					err := p.dispatcher.Handle(msg.Session(), msg.Msg())
					if err != nil {
						log.WithFields(
							log.Fields{
								"房间id":  p.r.RoomId(),
								"error": err,
							}).Error("房间处理器,错误")
					}
				}
			case <-time.After(p.tick):
				{
					p.r.Tick()
				}
			}
		}
		log.WithFields(
			log.Fields{
				"房间id": p.r.RoomId(),
			}).Info("房间处理器结束")
	}()
	return nil
}

func (p *RoomProcessor) Stop() {
	close(p.messages)
}

func NewRoomProcessor(r *changsha.Room, queueSize int, tick time.Duration, dispatcher *Dispatcher) *RoomProcessor {
	p := &RoomProcessor{}
	p.r = r
	p.messages = make(chan *Message, queueSize)
	p.tick = tick
	p.dispatcher = dispatcher

	return p
}

const (
	roomProcessorKey = "game.roomprocessor"
)

func WithRoomProcessor(ctx context.Context, rp *RoomProcessor) context.Context {
	return context.WithValue(ctx, roomProcessorKey, rp)
}

func RoomProcessorInContext(ctx context.Context) *RoomProcessor {
	m, ok := ctx.Value(roomProcessorKey).(*RoomProcessor)
	if !ok {
		return nil
	}
	return m
}
