package handler

import (
	log "github.com/Sirupsen/logrus"

	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"
	"game/session"

	"game/basic/pb"
)

func handlePlayerStart(s session.Session, msg *pb.Message) error {
	log.Println("receive player start")
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
	roomPlayer := r.RoomPlayerManager().GetPlayerById(pl.Id())
	if roomPlayer == nil {
		log.WithFields(
			log.Fields{
				"玩家id": pl.Id(),
			}).Warn("玩家不存在")
		return nil
	}
	r.PrepareStart(roomPlayer)
	return nil
}
