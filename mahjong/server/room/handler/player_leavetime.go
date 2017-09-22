package handler

import (
	changshapb "game/mahjong/pb/changsha" 

	log "github.com/Sirupsen/logrus"

	"game/session"

	"game/basic/pb"

	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"

	"github.com/golang/protobuf/proto"
)

//玩家离开
func handlePlayerLeaveTime(s session.Session, msg *pb.Message) error {

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
			"玩家id":   pl.Id(),
			"房间id":   r.RoomId(),
			"房间主人id": r.OwnerId(),
		}).Debug("接收消息,玩家自主离开")

	extend, err := proto.GetExtension(msg, changshapb.E_CgLeaveTime)
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

	cgLeaveTime := extend.(*changshapb.CGLeaveTime)
	flag := r.PlayerLeaveTime(roomPlayer, cgLeaveTime.GetFlag())

	if flag { 
		log.WithFields(
			log.Fields{
				"玩家id":   pl.Id(),
				"房间id":   r.RoomId(),
				"房间主人id": r.OwnerId(),
			}).Debug("准备自离开主房间") 
		return nil
	}
	return nil
}
