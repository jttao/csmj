package handler

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handlePlayerDissolve(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerDissolve)
	if err != nil {
		return err
	}
	gcPlayerDissolve, ok := val.(*changshapb.GCPlayerDissolve)
	if !ok {
		panic("no gc player dissolve")
	}
	uId := gcPlayerDissolve.GetPlayerId()
	if uId == c.Id() {
		fmt.Println("自己申请解散中")
		return nil
	}
	fmt.Println("玩家申请解散中")

	c.PrepareDissolve()
	return nil
}
