package handler

import (
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"

	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleBankerSelect(c *client.Client, msg *pb.Message) error {
	log.Println("receive bank select  msg")

	val, err := proto.GetExtension(msg, changshapb.E_GcSelectBankerBroadcast)
	if err != nil {
		return err
	}
	gcSelectBankerBroadcast, ok := val.(*changshapb.GCSelectBankerBroadcast)
	if !ok {
		panic("no gc selec banker broadcast")
	}

	r := c.Room()

	r.BankerSelect(gcSelectBankerBroadcast.GetBankerPos())

	return nil
}
