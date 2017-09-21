package handler

import (
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/changsha"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleXiaoHuBroadcast(c *client.Client, msg *pb.Message) error {
	log.Println("receive xiao hu broadcast msg")

	val, err := proto.GetExtension(msg, changshapb.E_GcXiaoHuBroadcast)
	if err != nil {
		return err
	}
	gcXiaoHuBroadcast, ok := val.(*changshapb.GCXiaoHuBroadcast)
	if !ok {
		panic("no gc join room")
	}

	r := c.Room()

	for _, xhp := range gcXiaoHuBroadcast.GetXiaohuPlayerList() {
		pId := xhp.GetPlayerId()
		xhos := make([]*changsha.XiaoHuOperation, 0, len(xhp.GetXiaoHuList()))
		for _, xh := range xhp.GetXiaoHuList() {
			xhos = append(xhos, xiaoHuOpreationsFromXiaoHu(xh))
		}
		r.XiaoHu(pId, xhos)
	}

	return nil
}
