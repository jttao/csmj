package handler

import clienthandler "game/mahjong/client/handler"
import "game/mahjong/pb/messagetype"

func InitDispatch(d clienthandler.Dispatcher) {
	d.Register((int)(messagetype.MessageType_GCLoginType), clienthandler.HandlerFunc(handleLogin))
	d.Register(int(messagetype.MessageType_GCJoinRoomType), clienthandler.HandlerFunc(handleJoinRoom))
	d.Register(int(messagetype.MessageType_GCJoinRoomBroadcastType), clienthandler.HandlerFunc(handleJoinRoomBroadcast))

	d.Register(int(messagetype.MessageType_GCSelectBankerBroadcastType), clienthandler.HandlerFunc(handleBankerSelect))
	d.Register(int(messagetype.MessageType_GCDealCardsType), clienthandler.HandlerFunc(handleDealCards))

	d.Register(int(messagetype.MessageType_GCXiaoHuListType), clienthandler.HandlerFunc(handleXiaoHuList))
	d.Register(int(messagetype.MessageType_GCWaitPlayerPlayCardType), clienthandler.HandlerFunc(handleWaitPlayerPlayCard))
	d.Register(int(messagetype.MessageType_GCPlayerPlayCardType), clienthandler.HandlerFunc(handlePlayerPlayCard))
	d.Register(int(messagetype.MessageType_GCPlayerOperationsType), clienthandler.HandlerFunc(handlePlayerOpreations))
	d.Register(int(messagetype.MessageType_GCPlayerChiType), clienthandler.HandlerFunc(handlePlayerChi))
	d.Register(int(messagetype.MessageType_GCPlayerPengType), clienthandler.HandlerFunc(handlePlayerPeng))
	d.Register(int(messagetype.MessageType_GCPlayerBuType), clienthandler.HandlerFunc(handlePlayerBu))
	d.Register(int(messagetype.MessageType_GCPlayerGangType), clienthandler.HandlerFunc(handlePlayerGang))

	d.Register(int(messagetype.MessageType_GCPlayerMoCardType), clienthandler.HandlerFunc(handlePlayerMo))
	d.Register(int(messagetype.MessageType_GCPlayerMoCardBroadcastType), clienthandler.HandlerFunc(handlePlayerMoBroadcast))
	d.Register(int(messagetype.MessageType_GCPlayerHuType), clienthandler.HandlerFunc(handleSettle))
	d.Register(int(messagetype.MessageType_GCClearType), clienthandler.HandlerFunc(handleClear))
	d.Register(int(messagetype.MessageType_GCPlayerGangMoType), clienthandler.HandlerFunc(handlePlayerGangMo))
	d.Register(int(messagetype.MessageType_GCPlayerLeaveType), clienthandler.HandlerFunc(handleLeaveRoom))
	d.Register(int(messagetype.MessageType_GCPlayerDissolveType), clienthandler.HandlerFunc(handlePlayerDissolve))
	d.Register(int(messagetype.MessageType_GCTotalSettleType), clienthandler.HandlerFunc(handleTotalSettle))
	d.Register(int(messagetype.MessageType_GCXiaohuType), clienthandler.HandlerFunc(handleXiaoHu))
	d.Register(int(messagetype.MessageType_GCHaiDiAskType), clienthandler.HandlerFunc(handleHaidiAsk))
}
