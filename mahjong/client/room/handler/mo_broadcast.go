package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerMoBroadcast(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerMoCardBroadcast)
	if err != nil {
		return err
	}
	gcPlayerMoCardBroadcast, ok := val.(*changshapb.GCPlayerMoCardBroadcast)
	if !ok {
		panic("no gc player play card")
	}

	r := c.Room()
	r.PlayerMo(gcPlayerMoCardBroadcast.GetPlayerId())
	return nil
}
