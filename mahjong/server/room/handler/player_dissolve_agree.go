package handler

import (
	changshapb "game/mahjong/pb/changsha"
	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"

	"game/session"

	"game/basic/pb"
	
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

//玩家同意解散
func handlePlayerDissolveAgree(s session.Session, msg *pb.Message) error {

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
		}).Debug("玩家请求同意解散")
	extend, err := proto.GetExtension(msg, changshapb.E_CgPlayerDissolveAgree)
	if err != nil {
		return err
	}
	roomPlayer := r.RoomPlayerManager().GetPlayerById(pl.Id())

	if roomPlayer == nil {
		log.WithFields(
			log.Fields{
				"玩家id": pl.Id(),
			}).Warn("玩家不存在")
		return nil
	}
	cgPlayerDissolveAgree := extend.(*changshapb.CGPlayerDissolveAgree)
	r.PlayerDissolveAgree(roomPlayer, cgPlayerDissolveAgree.GetAgree())
	return nil
}
