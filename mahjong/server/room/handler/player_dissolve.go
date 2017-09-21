package handler

import (
	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"
	"game/session"

	"game/basic/pb"

	log "github.com/Sirupsen/logrus"
)

//申请解散中
func handlePlayerDissolve(s session.Session, msg *pb.Message) error {

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

	log.WithFields(
		log.Fields{
			"玩家id": pl.Id(),
			"房间id": r.RoomId(),
		}).Debug("玩家申请解散")

	roomPlayer := r.RoomPlayerManager().GetPlayerById(pl.Id())
	if roomPlayer == nil {
		log.WithFields(
			log.Fields{
				"玩家id": pl.Id(),
			}).Warn("玩家不存在")
		return nil
	}
	r.PlayerDissolve(roomPlayer)
	return nil
}
