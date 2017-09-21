package handler

import (
	log "github.com/Sirupsen/logrus"

	"game/mahjong/changsha"
	changshapb "game/mahjong/pb/changsha"
	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"

	"game/session"

	"game/basic/pb"

	"github.com/golang/protobuf/proto"
)

func handlePlayerOpreate(s session.Session, msg *pb.Message) error {
	log.Println("receive player operate")
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
	extend, err := proto.GetExtension(msg, changshapb.E_CgPlayerOpreate)
	if err != nil {
		return err
	}
	cgPlayerOperate := extend.(*changshapb.CGPlayerOperate)

	roomPlayer := r.RoomPlayerManager().GetPlayerById(pl.Id())
	if roomPlayer == nil {
		log.WithFields(
			log.Fields{
				"玩家id": pl.Id(),
			}).Warn("玩家不存在")
		return nil
	}
	r.Operate(roomPlayer, cgPlayerOperate.GetOperation().GetTargetIndex(), changsha.OperationType(cgPlayerOperate.GetOperation().GetOperationType()), cgPlayerOperate.GetOperation().GetTargetCard(), cgPlayerOperate.GetOperation().GetCardList())
	return nil
}
