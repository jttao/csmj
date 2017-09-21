package handler

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleTotalSettle(c *client.Client, msg *pb.Message) error {

	_, err := proto.GetExtension(msg, changshapb.E_GcTotalSettle)
	if err != nil {
		return err
	}

	fmt.Println("总结算")

	return nil
}
