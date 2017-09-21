package handler

import (
	"fmt"
	"log"

	"game/basic/pb"
	"game/mahjong/client/client"
)

type HandlerFunc func(c *client.Client, msg *pb.Message) error

func (hf HandlerFunc) Handle(c *client.Client, msg *pb.Message) error {
	return hf(c, msg)
}

type Handler interface {
	Handle(c *client.Client, msg *pb.Message) error
}

type Dispatcher interface {
	Register(messageType int, h Handler)
	Handler
}

type dispatcher struct {
	handlerMap map[int]Handler
}

func (d *dispatcher) Register(messageType int, h Handler) {
	d.handlerMap[messageType] = h
}

func (d *dispatcher) Handle(c *client.Client, msg *pb.Message) error {
	msgType := int(*msg.MessageType)
	log.Printf("receive msg type [%d]\n", msgType)
	h, exist := d.handlerMap[msgType]
	if !exist {
		return fmt.Errorf("handler no exist")
	}
	return h.Handle(c, msg)
}

func NewDispatcher() Dispatcher {
	d := &dispatcher{}
	d.handlerMap = make(map[int]Handler)
	return d
}
