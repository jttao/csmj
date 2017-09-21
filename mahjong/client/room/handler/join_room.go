package handler

import (
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleJoinRoom(c *client.Client, msg *pb.Message) error {
	log.Println("receive join room msg")

	val, err := proto.GetExtension(msg, changshapb.E_GcJoinRoom)
	if err != nil {
		return err
	}
	gcJoinRoom, ok := val.(*changshapb.GCJoinRoom)
	if !ok {
		panic("no gc join room")
	}

	r := roomFromRoomInfo(gcJoinRoom.RoomInfo)
	c.SetRoom(r)

	for _, playerInfo := range gcJoinRoom.PlayerInfoList {
		pl := playerFromPlayerInfo(playerInfo)
		r.PlayerJoin(pl)
	}

	pl := r.PlayerManager().GetPlayerById(c.Id())
	c.SetPlayer(pl)

	if gcJoinRoom.RoomInfo.GetState() == 5 {

		if len(gcJoinRoom.RoomInfo.GetCurrentOperationList()) != 0 {
			c.PreparePlayerOpereations(gcJoinRoom.RoomInfo.GetCurrentOperationList())
		}
	}

	if gcJoinRoom.RoomInfo.GetState() == 4 {
		if gcJoinRoom.RoomInfo.GetCurrentPlayerId() == c.Id() {

			c.PreaprePlay()

		}
	}

	return nil
}
