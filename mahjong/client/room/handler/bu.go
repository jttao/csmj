package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerBu(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerBu)
	if err != nil {
		return err
	}
	gcPlayerBu, ok := val.(*changshapb.GCPlayerBu)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	pId := gcPlayerBu.GetPlayerId()
	r.PlayerBu(pId, gcPlayerBu.GetCard(), gcPlayerBu.GetBuType())
	return nil
}
