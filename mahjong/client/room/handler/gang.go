package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerGang(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerGang)
	if err != nil {
		return err
	}
	gcPlayerGang, ok := val.(*changshapb.GCPlayerGang)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	pId := gcPlayerGang.GetPlayerId()
	r.PlayerGang(pId, gcPlayerGang.GetCard(), gcPlayerGang.GetGangType())
	return nil
}
