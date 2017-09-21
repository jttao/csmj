package util

import (
	"game/basic/pb"
	commonpb "game/mahjong/pb/common"
	messagetypepb "game/mahjong/pb/messagetype"
	"game/session"

	"github.com/golang/protobuf/proto"
)

func CloseWithError(s session.Session, errorCode int32) error {
	gcError := &commonpb.GCError{}
	gcError.ErrorCode = &errorCode
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCError)
	msg.MessageType = &msgType
	err := proto.SetExtension(msg, commonpb.E_GcError, gcError)
	if err != nil {
		return err
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = s.Send(msgBytes)
	if err != nil {
		return err
	}
	return s.Close()
}
