package room

import (
	"encoding/json"

	"game/db"
	"game/mahjong/card"
	"game/mahjong/changsha"
	changshapb "game/mahjong/pb/changsha"
	messagetypepb "game/mahjong/pb/messagetype"
	recordmodel "game/mahjong/record/model"
	"game/mahjong/server/mahjong"

	"game/basic/pb"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

type RoomDelegate struct {
	rm *changsha.RoomManager
	ds db.DBService
}

func (rd *RoomDelegate) OnRoomStart(r *changsha.Room) {

}

func (rd *RoomDelegate) OnRoomPlayerJoin(r *changsha.Room, player changsha.Player) {

	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCJoinRoomType)
	msg.MessageType = &msgType

	gcJoinRoom := buildGCJoinRoom(r, player)

	err := proto.SetExtension(msg, changshapb.E_GcJoinRoom, gcJoinRoom)
	if err != nil {
		panic(err)
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	send(player, msgBytes)

	msgBroadcast := &pb.Message{}
	msgBroadcastType := int32(messagetypepb.MessageType_GCJoinRoomBroadcastType)
	msgBroadcast.MessageType = &msgBroadcastType
	gcJoinRoomBroadcast := &changshapb.GCJoinRoomBroadcast{}
	gcJoinRoomBroadcast.GcPlayerInfo = buildGCPlayerInfo(player, false)
	err = proto.SetExtension(msgBroadcast, changshapb.E_GcJoinRoomBroadcast, gcJoinRoomBroadcast)
	if err != nil {
		panic(err)
	}
	msgBroadcastBytes, err := proto.Marshal(msgBroadcast)
	if err != nil {
		panic(err)
	}

	broadcastExcept(r, msgBroadcastBytes, player)
}

func (rd *RoomDelegate) OnRoomPlayerStart(r *changsha.Room, player changsha.Player) {

	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerStartType)
	msg.MessageType = &msgType

	gcPlayerStart := &changshapb.GCPlayerStart{}
	pId := player.Id()
	gcPlayerStart.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcPlayerStart, gcPlayerStart)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerLeave(r *changsha.Room, player changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerLeaveType)
	msg.MessageType = &msgType

	gcPlayerLeave := &changshapb.GCPlayerLeave{}
	pId := player.Id()
	gcPlayerLeave.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcPlayerLeave, gcPlayerLeave)
	if err != nil {
		panic(err)
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerDissolve(r *changsha.Room, player changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerDissolveType)
	msg.MessageType = &msgType

	gcPlayerDissolve := &changshapb.GCPlayerDissolve{}
	pId := player.Id()
	gcPlayerDissolve.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcPlayerDissolve, gcPlayerDissolve)
	if err != nil {
		panic(err)
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerDissolveAgree(r *changsha.Room, player changsha.Player, flag bool) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerDissolveAgreeType)
	msg.MessageType = &msgType

	gcPlayerDissolveAgree := &changshapb.GCPlayerDissolveAgree{}
	pId := player.Id()
	gcPlayerDissolveAgree.PlayerId = &pId
	gcPlayerDissolveAgree.Agree = &flag
	err := proto.SetExtension(msg, changshapb.E_GcPlayerDissolveAgree, gcPlayerDissolveAgree)
	if err != nil {
		panic(err)
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomReconnect(r *changsha.Room, player changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCJoinRoomType)
	msg.MessageType = &msgType

	gcJoinRoom := buildGCJoinRoom(r, player)

	err := proto.SetExtension(msg, changshapb.E_GcJoinRoom, gcJoinRoom)
	if err != nil {
		panic(err)
	}

	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	send(player, msgBytes)

	msgBroadcast := &pb.Message{}
	msgBroadcastType := int32(messagetypepb.MessageType_GCReconnectPlayerType)
	msgBroadcast.MessageType = &msgBroadcastType

	gcReconncetPlayer := &changshapb.GCReconnectPlayer{}
	pId := player.Id()
	gcReconncetPlayer.PlayerId = &pId
	err = proto.SetExtension(msgBroadcast, changshapb.E_GcReconnectPlayer, gcReconncetPlayer)
	if err != nil {
		panic(err)
	}
	msgBroadcastBytes, err := proto.Marshal(msgBroadcast)
	if err != nil {
		panic(err)
	}

	broadcastExcept(r, msgBroadcastBytes, player)
}

func (rd *RoomDelegate) OnRoomDisconnect(r *changsha.Room, player changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCDisconnectPlayerType)
	msg.MessageType = &msgType

	gcDisconncetPlayer := &changshapb.GCDisconnectPlayer{}
	pId := player.Id()
	gcDisconncetPlayer.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcDisconnectPlayer, gcDisconncetPlayer)
	if err != nil {
		panic(err)
	}
	msgBroadcastBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	broadcast(r, msgBroadcastBytes)
}

func (rd *RoomDelegate) OnLevaeTime(r *changsha.Room, player changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCLeaveTimeType)
	msg.MessageType = &msgType

	gcLeaveTime := &changshapb.GCLeaveTime{}
	pId := player.Id()  
	flag := (player.State()==2)
	
	gcLeaveTime.PlayerId = &pId 
	gcLeaveTime.Flag = &flag

	err := proto.SetExtension(msg, changshapb.E_GcLeaveTime, gcLeaveTime)
	if err != nil {
		panic(err)
	}
	msgBroadcastBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	
	broadcast(r, msgBroadcastBytes)
}

func (rd *RoomDelegate) OnRoomSelectBanker(r *changsha.Room, bankPos int) {

	newBankPos := int32(bankPos)
	gcSelectBankerBroadcast := &changshapb.GCSelectBankerBroadcast{}
	gcSelectBankerBroadcast.BankerPos = &newBankPos
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCSelectBankerBroadcastType)
	msg.MessageType = &msgType
	if err := proto.SetExtension(msg, changshapb.E_GcSelectBankerBroadcast, gcSelectBankerBroadcast); err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	//广播消息
	broadcast(r, msgBytes)

	//纪录日志
}

//发送牌
func (rd *RoomDelegate) OnRoomDealCards(r *changsha.Room) {
	for _, rp := range r.RoomPlayerManager().Players() {

		gcDealCards := &changshapb.GCDealCards{}
		gcDealCards.Cards = CardValues(rp.Cards())
		msg := &pb.Message{}
		msgType := int32(messagetypepb.MessageType_GCDealCardsType)
		msg.MessageType = &msgType
		if err := proto.SetExtension(msg, changshapb.E_GcDealCards, gcDealCards); err != nil {
			panic(err)
		}
		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
		send(rp, msgBytes)
	}
}

//玩家小胡
func (rd *RoomDelegate) OnRoomPlayerXiaoHu(r *changsha.Room, player changsha.Player, xho *changsha.XiaoHuOperation) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCXiaohuType)
	msg.MessageType = &msgType
	gcXiaoHu := &changshapb.GCXiaoHu{}
	xiaoHu := buildXiaoHu(xho)
	gcXiaoHu.XiaoHu = xiaoHu
	pId := player.Id()
	gcXiaoHu.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcXiaohu, gcXiaoHu)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

//玩家小胡过
func (rd *RoomDelegate) OnRoomPlayerXiaoHuPass(room *changsha.Room, player changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCXiaoHuPassType)
	msg.MessageType = &msgType
	gcXiaoHuPass := &changshapb.GCXiaoHuPass{}
	pId := player.Id()
	gcXiaoHuPass.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcXiaoHuPass, gcXiaoHuPass)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	send(player, msgBytes)
}

//发小胡列表
func (rd *RoomDelegate) OnRoomXiaoHu(r *changsha.Room) {

	for _, pl := range r.RoomPlayerManager().Players() {
		if len(pl.XiaoHus()) == 0 {
			continue
		}
		if pl.Player() == nil {
			continue
		}
		msg := &pb.Message{}
		msgType := int32(messagetypepb.MessageType_GCXiaoHuListType)
		msg.MessageType = &msgType
		gcXiaoHuList := buildGCXiaoHuList(pl)
		err := proto.SetExtension(msg, changshapb.E_GcXiaoHuList, gcXiaoHuList)
		if err != nil {
			panic(err)
		}
		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
		send(pl, msgBytes)
	}

	// msg := &pb.Message{}
	// msgType := int32(messagetypepb.)
	// msg.MessageType = &msgType

	// gcXiaoHuBroadcast := &changshapb.GCXiaoHuBroadcast{}
	// gcXiaoHuBroadcast.XiaohuPlayerList = make([]*changshapb.XiaohuPlayer, 0, 1)

	// for _, pl := range r.RoomPlayerManager().Players() {
	// 	if len(pl.XiaoHus()) == 0 {
	// 		continue
	// 	}
	// 	xhp := buildXiaoHuPlayer(pl)
	// 	gcXiaoHuBroadcast.XiaohuPlayerList = append(gcXiaoHuBroadcast.XiaohuPlayerList, xhp)
	// }
	// err := proto.SetExtension(msg, changshapb.E_GcXiaoHuBroadcast, gcXiaoHuBroadcast)
	// if err != nil {
	// 	panic(err)
	// }
	// msgBytes, err := proto.Marshal(msg)
	// if err != nil {
	// 	panic(err)
	// }
	// broadcast(r, msgBytes)
}

//玩家摸牌
func (rd *RoomDelegate) OnRoomPlayerMo(r *changsha.Room, player changsha.Player, c *card.Card) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerMoCardType)
	msg.MessageType = &msgType

	gcPlayerMoCard := &changshapb.GCPlayerMoCard{}
	cv := card.Value(c)
	gcPlayerMoCard.Card = &cv
	err := proto.SetExtension(msg, changshapb.E_GcPlayerMoCard, gcPlayerMoCard)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	send(player, msgBytes)

	msgBroadcast := &pb.Message{}
	msgBroadcastType := int32(messagetypepb.MessageType_GCPlayerMoCardBroadcastType)
	msgBroadcast.MessageType = &msgBroadcastType
	gcPlayerMoCardBroadcast := &changshapb.GCPlayerMoCardBroadcast{}
	pId := player.Id()
	gcPlayerMoCardBroadcast.PlayerId = &pId
	err = proto.SetExtension(msgBroadcast, changshapb.E_GcPlayerMoCardBroadcast, gcPlayerMoCardBroadcast)
	if err != nil {
		panic(err)
	}
	msgBroadcastBytes, err := proto.Marshal(msgBroadcast)
	if err != nil {
		panic(err)
	}
	broadcastExcept(r, msgBroadcastBytes, player)
}

//等候玩家打牌
func (rd *RoomDelegate) OnRoomWaitPlayerPlay(r *changsha.Room, player changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCWaitPlayerPlayCardType)
	msg.MessageType = &msgType

	gcWaitPlayerPlayCard := &changshapb.GCWaitPlayerPlayCard{}
	pId := player.Id()
	gcWaitPlayerPlayCard.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcWaitPlayerPlayCard, gcWaitPlayerPlayCard)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	broadcast(r, msgBytes)
}

//等候玩家行为
func (rd *RoomDelegate) OnRoomWaitPlayerAction(r *changsha.Room) {
	for _, pl := range r.RoomPlayerManager().Players() {

		if len(pl.PossibleOperations()) == 0 {
			continue
		}
		msg := &pb.Message{}
		msgType := int32(messagetypepb.MessageType_GCPlayerOperationsType)
		msg.MessageType = &msgType

		gcPlayerOpreations := &changshapb.GCPlayerOperations{}
		gcPlayerOpreations.OperationList = buildPlayerOperations(pl)
		err := proto.SetExtension(msg, changshapb.E_GcPlayerOperations, gcPlayerOpreations)
		if err != nil {
			panic(err)
		}
		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}

		send(pl, msgBytes)
	}
}

func (rd *RoomDelegate) OnRoomPlayerPlayCard(r *changsha.Room, player changsha.Player, c *card.Card) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerPlayCardType)
	msg.MessageType = &msgType
	gcPlayerPlayCard := &changshapb.GCPlayerPlayCard{}
	pId := player.Id()
	cv := card.Value(c)
	gcPlayerPlayCard.PlayerId = &pId
	gcPlayerPlayCard.Card = &cv
	err := proto.SetExtension(msg, changshapb.E_GcPlayerPlayCard, gcPlayerPlayCard)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerChi(r *changsha.Room, player changsha.Player, c *card.Card, cardValues []int32) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerChiType)
	msg.MessageType = &msgType
	gcPlayerChi := &changshapb.GCPlayerChi{}
	pId := player.Id()
	gcPlayerChi.PlayerId = &pId
	cv := card.Value(c)
	gcPlayerChi.Card = &cv
	gcPlayerChi.CardList = cardValues

	err := proto.SetExtension(msg, changshapb.E_GcPlayerChi, gcPlayerChi)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerPeng(r *changsha.Room, player changsha.Player, c *card.Card) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerPengType)
	msg.MessageType = &msgType
	gcPlayerPeng := &changshapb.GCPlayerPeng{}
	pId := player.Id()
	gcPlayerPeng.PlayerId = &pId
	cv := card.Value(c)
	gcPlayerPeng.Card = &cv

	err := proto.SetExtension(msg, changshapb.E_GcPlayerPeng, gcPlayerPeng)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerBu(r *changsha.Room, player changsha.Player, c *card.Card, gt changsha.GangType) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerBuType)
	msg.MessageType = &msgType
	gcPlayerBu := &changshapb.GCPlayerBu{}
	pId := player.Id()
	gcPlayerBu.PlayerId = &pId
	cv := card.Value(c)
	gcPlayerBu.Card = &cv
	gtInt := int32(gt)
	gcPlayerBu.BuType = &gtInt

	err := proto.SetExtension(msg, changshapb.E_GcPlayerBu, gcPlayerBu)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerGang(r *changsha.Room, player changsha.Player, c *card.Card, gt changsha.GangType) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerGangType)
	msg.MessageType = &msgType
	gcPlayerGang := &changshapb.GCPlayerGang{}
	pId := player.Id()
	gcPlayerGang.PlayerId = &pId
	cv := card.Value(c)
	gcPlayerGang.Card = &cv
	gtInt := int32(gt)
	gcPlayerGang.GangType = &gtInt

	err := proto.SetExtension(msg, changshapb.E_GcPlayerGang, gcPlayerGang)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomLiuJu(r *changsha.Room) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCLiuJuType)
	msg.MessageType = &msgType
	gcLiuJu := &changshapb.GCLiuJu{}
	gcLiuJu.SettlePlayerList = buildRoundSettlePlayerList(r)

	err := proto.SetExtension(msg, changshapb.E_GcLiuJu, gcLiuJu)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)

}

func (rd *RoomDelegate) OnRoomPlayersSettle(r *changsha.Room, players []changsha.Player, dianPaoId int64, c *card.Card) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerHuType)
	msg.MessageType = &msgType
	gcPlayerHu := &changshapb.GCPlayerHu{}
	gcPlayerHu.HuList = buildHuList(players)
	gcPlayerHu.SettlePlayerList = buildRoundSettlePlayerList(r)
	tDianPaoId := dianPaoId
	gcPlayerHu.DianPaoPlayerId = &tDianPaoId
	if len(r.NiaoPais()) != 0 {
		gcPlayerHu.NiaoPais = card.Values(r.NiaoPais())
		gcPlayerHu.NiaoPaiPlayers = r.NiaoPaiPlayerIds()
	}
	var cv int32 = int32(0)
	if c != nil {
		cv = card.Value(c)
	}
	gcPlayerHu.Card = &cv

	err := proto.SetExtension(msg, changshapb.E_GcPlayerHu, gcPlayerHu)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)

}

func (rd *RoomDelegate) OnRoomClear(r *changsha.Room) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCClearType)
	msg.MessageType = &msgType
	gcClear := &changshapb.GCClear{}
	err := proto.SetExtension(msg, changshapb.E_GcClear, gcClear)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomPlayerGangMo(r *changsha.Room, player changsha.Player, cs []*card.Card) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCPlayerGangMoType)
	msg.MessageType = &msgType
	gcPlayerGangMo := &changshapb.GCPlayerGangMo{}
	gcPlayerGangMo.Card = buildCards(cs)
	pId := player.Id()
	gcPlayerGangMo.PlayerId = &pId
	err := proto.SetExtension(msg, changshapb.E_GcPlayerGangMo, gcPlayerGangMo)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomHaiDiAsk(r *changsha.Room, pl changsha.Player) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCHaiDiAskType)
	msg.MessageType = &msgType
	gcHaiDiAsk := &changshapb.GCHaiDiAsk{}
	pId := pl.Id()
	gcHaiDiAsk.PlayerId = &pId

	err := proto.SetExtension(msg, changshapb.E_GcHaiDiAsk, gcHaiDiAsk)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomHaiDiAnswer(r *changsha.Room, pl changsha.Player, c *card.Card, flag bool) {
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_GCHaiDiAnswerType)
	msg.MessageType = &msgType
	gcHaiDiAnswer := &changshapb.GCHaiDiAnswer{}
	pId := pl.Id()
	gcHaiDiAnswer.PlayerId = &pId
	gcHaiDiAnswer.Flag = &flag
	if c != nil {
		cv := card.Value(c)
		gcHaiDiAnswer.Card = &cv
	}
	err := proto.SetExtension(msg, changshapb.E_GcHaiDiAnswer, gcHaiDiAnswer)
	if err != nil {
		panic(err)
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	broadcast(r, msgBytes)
}

func (rd *RoomDelegate) OnRoomSettle(r *changsha.Room) {
	err := rd.saveRound(r)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("save round error")
	}
}

func (rd *RoomDelegate) OnRoomEnd(r *changsha.Room,start bool) {
	
	if start {

		msg := &pb.Message{}
		msgType := int32(messagetypepb.MessageType_GCTotalSettleType)
		msg.MessageType = &msgType
		gcTotalSettle := buildGCTotalSettle(r)
	
		err := proto.SetExtension(msg, changshapb.E_GcTotalSettle, gcTotalSettle)
		if err != nil {
			panic(err)
		}
	
		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
	
		broadcast(r, msgBytes)
	
		err = rd.saveRoom(r)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Warn("save room error")
		}    
	}
	
	rp := mahjong.RoomProcessorInContext(r.Context())
	rp.Stop()

	//玩家推出游戏服
	players := r.RoomPlayerManager().Players()

	for _, pl := range r.RoomPlayerManager().Players() {
		if pl.Player() == nil {
			continue
		} 
		err := pl.Player().Session().Close()
		if err != nil {
			log.Println("close with error", err.Error())
		}
	} 
	
	mahjongContext := mahjong.MahjongInContext(r.Context())
	//移除房间
	mahjongContext.RoomManager.RemoveRoom(r)

	//远程摧毁房间
	log.WithFields(
		log.Fields{
			"房间id":   r.RoomId(),
			"房间主人id": r.OwnerId(),
		}).Debug("准备远程摧毁房间")
	
	refund := false
	if r.CurrentRound() <= 1 {
		refund = true
	}
	
	err := mahjongContext.RoomManageClient.Destroy(r.RoomId(), refund)
	
	if err != nil {
		log.WithFields(
			log.Fields{
				"房间id":   r.RoomId(),
				"房间主人id": r.OwnerId(),
				"error":  err.Error(),
			}).Debug("远程摧毁房间失败")
	}
	
	//设置玩家任务
	if start {

		for _, pl := range players {
			
			if pl.Player() == nil {
				continue
			} 

			playerId := pl.Id()
			score := pl.Score() 
				
			err = mahjongContext.RoomManageClient.GameEnd(playerId,score) 
			if err != nil {
				log.WithFields(
					log.Fields{
						"playerId":   playerId,
						"score": score,
						"error":  err.Error(),
					}).Debug("游戏结束，玩家任务失败")
			} 			
		} 

	}

}

func (rd *RoomDelegate) saveRoom(r *changsha.Room) error {
	rou := &recordmodel.RoomRecordModel{}
	rou.RoomId = r.RoomId()
	rou.OwnerId = r.OwnerId()
	for i := 0; i < len(r.RoomPlayerManager().Players()); i++ {
		pId := r.RoomPlayerManager().Players()[i].Id()
		switch i {
		case 0:
			rou.Player1 = pId
		case 1:
			rou.Player2 = pId
		case 2:
			rou.Player3 = pId
		case 3:
			rou.Player4 = pId
		}
	}
	settleInfo := make([]map[string]interface{}, 0, len(r.RoomPlayerManager().Players()))

	for _, pl := range r.RoomPlayerManager().Players() {
		playerSettleInfo := make(map[string]interface{})
		playerSettleInfo["playerId"] = pl.Id()
		playerSettleInfo["score"] = pl.Score()
		playerSettleInfo["name"] = pl.Name()
		settleInfo = append(settleInfo, playerSettleInfo)
	}

	settleBytes, err := json.Marshal(settleInfo)
	if err != nil {
		return err
	}
	rou.Settle = string(settleBytes)
	rou.CreateTime = r.CreateTime()
	tdb := rd.ds.DB().Save(rou)
	if tdb.Error != nil {
		return tdb.Error
		//log.Error("save round error")
	} 
	return nil
}

func (rd *RoomDelegate) saveRound(r *changsha.Room) error {
	rou := &recordmodel.RoundModel{}
	rou.RoomId = r.RoomId()
	rou.Round = r.CurrentRound()
	rou.TotoalRound = r.TotalRound()
	logContent, err := json.Marshal(r.LogList())
	if err != nil {
		//log.Error("json Marshal error")
		return err
	}
	rou.Logs = string(logContent)

	settleInfo := make([]map[string]interface{}, 0, len(r.RoomPlayerManager().Players()))

	for _, pl := range r.RoomPlayerManager().Players() {
		playerSettleInfo := make(map[string]interface{})
		playerSettleInfo["playerId"] = pl.Id()
		playerSettleInfo["score"] = pl.CurrentScore()
		playerSettleInfo["name"] = pl.Name()
		settleInfo = append(settleInfo, playerSettleInfo)
	}
	settleBytes, err := json.Marshal(settleInfo)
	if err != nil {
		return err
	}
	configBytes, err := json.Marshal(r.CustomRoomConfig())
	rou.Config = string(configBytes)
	rou.Settle = string(settleBytes)
	rou.CreateTime = r.RoundStartTime()
	tdb := rd.ds.DB().Save(rou)
	if tdb.Error != nil {
		return tdb.Error
		//log.Error("save round error")
	}
	return nil
}

func send(p changsha.Player, content []byte) {
	if p.Player() == nil {
		return
	}
	p.Player().Send(content)

}

func broadcast(r *changsha.Room, content []byte) {
	for _, pl := range r.RoomPlayerManager().Players() {
		if pl.Player() == nil {
			continue
		}
		pl.Player().Send(content)
	}
	log.Println("广播消息")
}

func broadcastExcept(r *changsha.Room, content []byte, p changsha.Player) {
	for _, pl := range r.RoomPlayerManager().Players() {
		if pl == p {
			continue
		}
		if pl.Player() == nil {
			continue
		}
		pl.Player().Send(content)

	}
	log.Println("多播消息")
}

func NewRoomDelegate(rm *changsha.RoomManager, ds db.DBService) *RoomDelegate {
	rd := &RoomDelegate{}
	rd.ds = ds
	rd.rm = rm
	return rd
}

func CardValues(cs []*card.Card) []int32 {
	vs := make([]int32, 0, len(cs))
	for _, c := range cs {
		vs = append(vs, card.Value(c))
	}
	return vs
}

