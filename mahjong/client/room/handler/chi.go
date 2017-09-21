package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerChi(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerChi)
	if err != nil {
		return err
	}
	gcPlayerChi, ok := val.(*changshapb.GCPlayerChi)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	pId := gcPlayerChi.GetPlayerId()
	r.PlayerChi(pId, gcPlayerChi.GetCard(),gcPlayerChi.GetCardList())
	return nil
}
