package handler

import (
	log "github.com/Sirupsen/logrus"

	changshapb "game/mahjong/pb/changsha"

	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"
	"game/session"

	"game/basic/pb"

	"github.com/golang/protobuf/proto"
)

func handlePlayerPlayCard(s session.Session, msg *pb.Message) error {
	log.Println("receive player play")
	pl := player.PlayerInContext(s.Context())
	mahjongContext := mahjong.MahjongInContext(s.Context())
	r := mahjongContext.RoomManager.GetRoomById(pl.RoomId())
	if r == nil {
		log.WithFields(
			log.Fields{
				"玩家id": pl.Id(),
			}).Warn("房间不存在")
		return nil
	}
	extend, err := proto.GetExtension(msg, changshapb.E_CgPlayerPlayCard)
	if err != nil {
		return err
	}
	cgPlayerPlayCard := extend.(*changshapb.CGPlayerPlayCard)

	roomPlayer := r.RoomPlayerManager().GetPlayerById(pl.Id())
	if roomPlayer == nil {
		log.WithFields(
			log.Fields{
				"玩家id": pl.Id(),
			}).Warn("玩家不存在")
		return nil
	}
	r.Play(roomPlayer, cgPlayerPlayCard.GetCard())
	return nil
}
