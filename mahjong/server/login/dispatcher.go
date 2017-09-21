package login

import (
	messagetypepb "game/mahjong/pb/messagetype"
	loginhandler "game/mahjong/server/login/handler"
	"game/mahjong/server/mahjong"
)

func InitDispatcher(d *mahjong.Dispatcher) {
	d.Register(int32(messagetypepb.MessageType_CGLoginType), mahjong.MessageHandlerFunc(loginhandler.HandleLogin))
}
