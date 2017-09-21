package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerPeng(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerPeng)
	if err != nil {
		return err
	}
	gcPlayerPeng, ok := val.(*changshapb.GCPlayerPeng)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	pId := gcPlayerPeng.GetPlayerId()
	r.PlayerPeng(pId, gcPlayerPeng.GetCard())
	return nil
}
