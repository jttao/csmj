package handler

import (
	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleJoinRoomBroadcast(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcJoinRoomBroadcast)
	if err != nil {
		return err
	}
	gcJoinRoomBroadcast := val.(*changshapb.GCJoinRoomBroadcast)
	pl := playerFromPlayerInfo(gcJoinRoomBroadcast.GcPlayerInfo)

	r := c.Room()
	r.PlayerJoin(pl)
	return nil
}
