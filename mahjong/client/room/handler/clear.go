package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleClear(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcClear)
	if err != nil {
		return err
	}
	_, ok := val.(*changshapb.GCClear)
	if !ok {
		panic("no gc player play card")
	}
	r := c.Room()
	r.Clear()
	return nil
}
