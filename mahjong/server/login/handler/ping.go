package handler

import (
	"game/basic/pb"
	loginpb "game/mahjong/pb/login"
	messagetypepb "game/mahjong/pb/messagetype"
	"game/mahjong/server/player"
	"game/session"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

func HandlePing(s session.Session, msg *pb.Message) error {
	log.Debug("处理ping消息")
	pl := player.PlayerInContext(s.Context())
	if pl == nil {
		log.WithFields(
			log.Fields{
				"sessionId": s.Id(),
			}).Warn("ping,玩家还没登陆")
		return nil
	}
	pl.Ping()
	m := &pb.Message{}
	pingMsgType := int32(messagetypepb.MessageType_GCPingType)
	m.MessageType = &pingMsgType
	gcPing := &loginpb.GCPing{}
	now := time.Now().UnixNano() / int64(time.Millisecond)
	gcPing.Now = &now
	err := proto.SetExtension(m, loginpb.E_GcPing, gcPing)
	if err != nil {
		return err
	}
	msgBytes, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	pl.Send(msgBytes)
	return nil
}
