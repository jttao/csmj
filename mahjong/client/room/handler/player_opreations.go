package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerOpreations(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerOperations)
	if err != nil {
		return err
	}
	gcPlayerOpreations, ok := val.(*changshapb.GCPlayerOperations)
	if !ok {
		panic("no gc player play card")
	}

	c.PreparePlayerOpereations(gcPlayerOpreations.GetOperationList())
	return nil
}
