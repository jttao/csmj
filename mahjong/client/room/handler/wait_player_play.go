package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleWaitPlayerPlayCard(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcWaitPlayerPlayCard)
	if err != nil {
		return err
	}
	gcWaitPlayerPlayCard, ok := val.(*changshapb.GCWaitPlayerPlayCard)
	if !ok {
		panic("no gc wait player play card")
	}

	r := c.Room()
	pId := gcWaitPlayerPlayCard.GetPlayerId()
	r.WaitPlayerPlay(pId)

	if pId == c.Id() {
		c.PreaprePlay()
	}

	return nil
}
