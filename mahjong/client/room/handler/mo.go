package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerMo(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerMoCard)
	if err != nil {
		return err
	}
	
	gcPlayerMoCard, ok := val.(*changshapb.GCPlayerMoCard)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	r.PlayerMoCard(c.Id(),gcPlayerMoCard.GetCard())
	return nil
}
