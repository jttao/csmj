package handler

import (
	messagetypepb "game/mahjong/pb/messagetype"
	loginhandler "game/mahjong/server/login/handler"
	"game/mahjong/server/mahjong"
)

//初始化分发器
func InitDispatcher(d *mahjong.Dispatcher) {
	d.Register(int32(messagetypepb.MessageType_CGPingType), mahjong.MessageHandlerFunc(loginhandler.HandlePing))

	d.Register(int32(messagetypepb.MessageType_CGPlayerPlayCardType), mahjong.MessageHandlerFunc(handlePlayerPlayCard))
	d.Register(int32(messagetypepb.MessageType_CGPlayerOperateType), mahjong.MessageHandlerFunc(handlePlayerOpreate))
	d.Register(int32(messagetypepb.MessageType_CGChatType), mahjong.MessageHandlerFunc(handleChat))
	d.Register(int32(messagetypepb.MessageType_CGPlayerLeaveType), mahjong.MessageHandlerFunc(handlePlayerLeave))
	d.Register(int32(messagetypepb.MessageType_CGPlayerStartType), mahjong.MessageHandlerFunc(handlePlayerStart))
	d.Register(int32(messagetypepb.MessageType_CGPlayerDissolveType), mahjong.MessageHandlerFunc(handlePlayerDissolve))
	d.Register(int32(messagetypepb.MessageType_CGPlayerDissolveAgreeType), mahjong.MessageHandlerFunc(handlePlayerDissolveAgree))
	d.Register(int32(messagetypepb.MessageType_CGXiaoHuType), mahjong.MessageHandlerFunc(handlePlayerXiaoHu))
	d.Register(int32(messagetypepb.MessageType_CGXiaoHuPassType), mahjong.MessageHandlerFunc(handlePlayerXiaoHuPass))
	d.Register(int32(messagetypepb.MessageType_CGHaiDiAnswerType), mahjong.MessageHandlerFunc(handlePlayerHaidiAnswer))
	d.Register(int32(messagetypepb.MessageType_CGLeaveTimeType), mahjong.MessageHandlerFunc(handlePlayerLeaveTime))
}
