package handler

import (
	log "github.com/Sirupsen/logrus"

	"game/session"

	"game/basic/pb"

	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"
)

//玩家离开
func handlePlayerLeave(s session.Session, msg *pb.Message) error { 
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
			"房间主人id": r.OwnerPlayer().Id(),
		}).Debug("接收消息,玩家准备离开")

	roomPlayer := r.RoomPlayerManager().GetPlayerById(pl.Id())
	if roomPlayer == nil {
		log.WithFields(
			log.Fields{
				"玩家id": pl.Id(),
			}).Warn("玩家不存在")
		return nil
	}
	
	flag := r.LeavePlayer(roomPlayer)	
	//flag := r.PlayerLeave(roomPlayer)
	if flag {
		/***
		if roomPlayer == r.OwnerPlayer() {
			rp := mahjong.RoomProcessorInContext(r.Context())
			rp.Stop()

			mahjongContext.RoomManager.RemoveRoom(r)
			log.WithFields(
				log.Fields{
					"玩家id":   pl.Id(),
					"房间id":   r.RoomId(),
					"房间主人id": r.OwnerPlayer().Id(),
				}).Debug("准备远程摧毁房间")
			err := mahjongContext.RoomManageClient.Destroy(r.RoomId(), true)
			if err != nil {
				log.WithFields(
					log.Fields{
						"玩家id":   pl.Id(),
						"房间id":   r.RoomId(),
						"房间主人id": r.OwnerPlayer().Id(),
						"error":  err.Error(),
					}).Debug("远程摧毁房间失败")
			}
			return nil
		}	
		***/
		log.WithFields(
			log.Fields{
				"玩家id":   pl.Id(),
				"房间id":   r.RoomId(),
				"房间主人id": r.OwnerPlayer().Id(),
			}).Debug("准备远程离开房间")
		err := mahjongContext.RoomManageClient.Leave(pl.Id(), r.RoomId())
		if err != nil {
			log.WithFields(
				log.Fields{
					"playerId": pl.Id(),
					"roomId":   r.RoomId(),
					"error":    err.Error(),
				}).Error("远程离开房间失败")
		}
		return nil
	}
	return nil 

}
