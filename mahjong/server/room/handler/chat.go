package handler

import (
	log "github.com/Sirupsen/logrus"

	messagetypepb "game/mahjong/pb/messagetype"

	"game/basic/pb"

	"github.com/golang/protobuf/proto"

	changshapb "game/mahjong/pb/changsha"

	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"
	"game/session"
)

//聊天信息
func handleChat(s session.Session, msg *pb.Message) error {

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
		}).Debug("接收聊天")

	extend, err := proto.GetExtension(msg, changshapb.E_CgChat)
	if err != nil {
		log.WithFields(
			log.Fields{
				"error": err,
			}).Warn("接收聊天,错误")
		return nil
	}

	//TODO 字数限制
	//转发聊天消息
	cgChat := extend.(*changshapb.CGChat)
	gcChat := &changshapb.GCChat{}
	gcChat.Typ = cgChat.Typ
	gcChat.Content = cgChat.Content

	senderId := pl.Id()
	gcChat.Sender = &senderId
	chatMsg := &pb.Message{}
	chatMsgType := int32(messagetypepb.MessageType_GCChatType)
	chatMsg.MessageType = &chatMsgType
	err = proto.SetExtension(chatMsg, changshapb.E_GcChat, gcChat)

	if err != nil {
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Error("receive msg chat error")
		return err
	}

	msgBytes, err := proto.Marshal(chatMsg)
	if err != nil {
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Error("receive msg chat error")
		return err
	}

	for _, pl := range r.RoomPlayerManager().Players() {
		if pl.Player() == nil {
			continue
		}
		pl.Player().Send(msgBytes)

	}
	return nil
}
