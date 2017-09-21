package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerGangMo(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerGangMo)
	if err != nil {
		return err
	}
	gcPlayerGangMo, ok := val.(*changshapb.GCPlayerGangMo)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	r.PlayerGangMo(gcPlayerGangMo.GetPlayerId(), gcPlayerGangMo.GetCard())
	return nil
}
