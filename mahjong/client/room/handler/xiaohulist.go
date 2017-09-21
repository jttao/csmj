package handler

import (
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleXiaoHuList(c *client.Client, msg *pb.Message) error {
	log.Println("receive xiao hu list msg")

	val, err := proto.GetExtension(msg, changshapb.E_GcXiaoHuList)
	if err != nil {
		return err
	}
	gcXiaoHuList, ok := val.(*changshapb.GCXiaoHuList)
	if !ok {
		panic("no gc xiaohu list")
	}

	//r := c.Room()

	for _, xh := range gcXiaoHuList.GetXiaoHus() {
		c.PrepareXiaoHu(xh)
	}

	return nil
}
