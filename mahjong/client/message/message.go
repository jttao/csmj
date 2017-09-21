package message

import (
	"github.com/coreos/etcd/client"
	"game/basic/pb"
)

type Message interface {
	Client() *client.Client
	Message() *pb.Message
}

type message struct {
    
}

func NewMessage(c *client.Client, msg *pb.Message) Message {

}
