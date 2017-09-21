package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleHaidiAsk(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcHaiDiAsk)
	if err != nil {
		return err
	}
	gcHaiDiAsk, ok := val.(*changshapb.GCHaiDiAsk)
	if !ok {
		panic("no gc haidi ask")
	}

	pId := gcHaiDiAsk.GetPlayerId()

	if pId == c.Id() {
		c.PrepareHaidiAnswer()
		return nil
	}

	return nil
}
