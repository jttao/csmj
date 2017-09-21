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

func handlePlayerXiaoHu(s session.Session, msg *pb.Message) error {
	log.Println("receive player xiaohu")
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
	extend, err := proto.GetExtension(msg, changshapb.E_CgXiaoHu)
	if err != nil {
		return err
	}
	cgXiaoHu := extend.(*changshapb.CGXiaoHu)
	r.PlayerXiaoHu(roomPlayer, cgXiaoHu.GetXiaoHu())
	return nil
}
