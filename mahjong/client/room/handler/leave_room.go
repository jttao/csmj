package handler

import (
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleLeaveRoom(c *client.Client, msg *pb.Message) error {
	log.Println("receive leave room msg")

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerLeave)
	if err != nil {
		return err
	}
	gcPlayerLeave, ok := val.(*changshapb.GCPlayerLeave)
	if !ok {
		panic("no gc leave room")
	}

	pId := gcPlayerLeave.GetPlayerId()
	r := c.Room()
	pl := r.PlayerManager().GetPlayerById(pId)
	r.PlayerLeave(pl)
	if pId == c.Id() {
		c.Close()
	}
	return nil
}
