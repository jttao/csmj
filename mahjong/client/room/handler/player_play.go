package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerPlayCard(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerPlayCard)
	if err != nil {
		return err
	}
	gcPlayerPlayCard, ok := val.(*changshapb.GCPlayerPlayCard)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	pId := gcPlayerPlayCard.GetPlayerId()
	r.PlayerPlay(pId, gcPlayerPlayCard.GetCard())
	return nil
}
