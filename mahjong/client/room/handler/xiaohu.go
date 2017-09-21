package handler

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleXiaoHu(c *client.Client, msg *pb.Message) error {
	log.Println("receive xiao hu msg")

	val, err := proto.GetExtension(msg, changshapb.E_GcXiaohu)
	if err != nil {
		return err
	}
	gcXiaoHu, ok := val.(*changshapb.GCXiaoHu)
	if !ok {
		panic("no gc xiaohu list")
	}
	fmt.Printf("%s", gcXiaoHu)
	//r := c.Room()

	return nil
}
