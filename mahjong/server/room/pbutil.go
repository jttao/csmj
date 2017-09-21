package room

import (
	"fmt"

	"game/mahjong/card"
	"game/mahjong/changsha"
	changshapb "game/mahjong/pb/changsha"
)

func buildGCJoinRoom(r *changsha.Room, player changsha.Player) *changshapb.GCJoinRoom {
	gcJoinRoom := &changshapb.GCJoinRoom{}
	gcJoinRoom.RoomInfo = buildGCRoomInfo(r, player)
	gcJoinRoom.PlayerInfoList = make([]*changshapb.GCPlayerInfo, 0, r.RoomPlayerManager().NumPlayers())
	for _, pl := range r.RoomPlayerManager().Players() {
		var gcPlayerInfo *changshapb.GCPlayerInfo
		if pl == player {
			gcPlayerInfo = buildGCPlayerInfo(pl, true)
		} else {
			gcPlayerInfo = buildGCPlayerInfo(pl, false)
		}
		gcJoinRoom.PlayerInfoList = append(gcJoinRoom.PlayerInfoList, gcPlayerInfo)
	}
	return gcJoinRoom
}

func buildGCRoomInfo(r *changsha.Room, player changsha.Player) *changshapb.GCRoomInfo {
	gcRoomInfo := &changshapb.GCRoomInfo{}
	roomId := r.RoomId()
	gcRoomInfo.RoomId = &roomId
	state := int32(r.State())
	gcRoomInfo.State = &state
	remains := int32(r.Deck().Remains())
	gcRoomInfo.Remains = &remains
	round := r.CurrentRound()
	gcRoomInfo.CurrentRound = &round
	name := r.Name()
	gcRoomInfo.Name = &name
	ownerId := r.OwnerPlayer().Id()
	gcRoomInfo.OwnerId = &ownerId
	totalRound := r.TotalRound()
	gcRoomInfo.TotalRound = &totalRound

	dissolveStartTime := r.DissolveStartTime()
	gcRoomInfo.DissolveStartTime = &dissolveStartTime

	zhuaNiao := int32(r.CustomRoomConfig().ZhuaNiao)
	zhuangXian := r.CustomRoomConfig().ZhuangXian
	zhuaNiaoAlg := r.CustomRoomConfig().ZhuaNiaoAlg
	gcRoomInfo.ZhuaNiao = &zhuaNiao
	gcRoomInfo.ZhuangXian = &zhuangXian
	gcRoomInfo.ZhuaNiaoAlg = &zhuaNiaoAlg

	forbidIp := r.ForbidIp() 
	gcRoomInfo.ForbidIp = &forbidIp

	applyPlayer := r.ApplyDissolvePlayer()
	if applyPlayer != nil {
		dissolvePlayers := make([]int64, 0, len(r.AgreeDissolvePlayers()))
		dissolvePlayers = append(dissolvePlayers, applyPlayer.Id())
		for _, tp := range r.AgreeDissolvePlayers() {
			if tp == applyPlayer {
				continue
			}
			dissolvePlayers = append(dissolvePlayers, tp.Id())
		}
		gcRoomInfo.AggreeDissolvePlayers = dissolvePlayers
	}

	bankerPos := int32(r.BankerPos())
	gcRoomInfo.BankerPos = &bankerPos

	switch r.State() {
	case changsha.RoomStateInit:
	case changsha.RoomStateWait:
	case changsha.RoomStatePrepare:
	case changsha.RoomStateWaitPlayerXiaoHuAction:
	case changsha.RoomStateSettle:
	case changsha.RoomStateWaitPlayerPlay:
		{
			if r.CurrentPlayer() != nil {
				pId := r.CurrentPlayer().Id()
				gcRoomInfo.CurrentPlayerId = &pId
			}
		}
	case changsha.RoomStateWaitPlayerAction:
		{
			if r.CurrentPlayer() != nil {
				pId := r.CurrentPlayer().Id()
				gcRoomInfo.CurrentPlayerId = &pId
				gcRoomInfo.CurrentCards = buildCards(r.CurrentCards())
				currentCardMethod := int32(r.CurrentCardPlayMethod())
				gcRoomInfo.CurrentCardMethod = &currentCardMethod
			}

			if len(player.PossibleOperations()) != 0 && !player.IfOperationsFinish(r.CurrentCards()) {
				gcRoomInfo.CurrentOperationList = buildPlayerOperations(player)
			}

		}
	}
	return gcRoomInfo
}

func buildGCPlayerInfo(player changsha.Player, isSelf bool) *changshapb.GCPlayerInfo {
	gcPlayerInfo := &changshapb.GCPlayerInfo{}
	img := player.Image()
	gcPlayerInfo.Image = &img
	ip := player.Ip()
	gcPlayerInfo.Ip = &ip
	id := player.Id()
	gcPlayerInfo.PlayerId = &id
	pos := int32(player.Position())
	gcPlayerInfo.Position = &pos
	score := player.Score()
	gcPlayerInfo.Score = &score
	state := int32(player.State())
	gcPlayerInfo.State = &state
	cardNum := int32(len(player.Cards()))
	gcPlayerInfo.CardNum = &cardNum
	gcPlayerInfo.ComposeList = buildComposeList(player.ComposeList())
	ifListen := player.IsListen()
	gcPlayerInfo.IfListen = &ifListen
	name := player.Name()
	gcPlayerInfo.Name = &name
	sex := int32(player.Sex())
	gcPlayerInfo.Sex = &sex
	location := player.Location()
	gcPlayerInfo.Location = &location
	
	if isSelf {
		gcPlayerInfo.CardList = buildCards(player.Cards())

		remainXiaohu := player.RemainXiaoHus()
		xiaoHus := make([]int32, 0, len(remainXiaohu))
		for _, xh := range remainXiaohu {
			xiaoHus = append(xiaoHus, int32(xh.XiaoHuType))
		}
		gcPlayerInfo.XiaoHus = xiaoHus
	}
	gcPlayerInfo.PlayCardList = buildCards(player.PlayedCards())
	actualState := int32(player.ActualState())
	gcPlayerInfo.ActualState = &actualState

	return gcPlayerInfo
}

func buildXiaoHuPlayer(player changsha.Player) *changshapb.XiaohuPlayer {
	xhp := &changshapb.XiaohuPlayer{}
	pId := player.Id()
	xhp.PlayerId = &pId
	xhp.XiaoHuList = make([]*changshapb.XiaoHu, 0, len(player.XiaoHus()))
	for _, xho := range player.XiaoHus() {
		xh := buildXiaoHu(xho)
		xhp.XiaoHuList = append(xhp.XiaoHuList, xh)
	}
	return xhp
}

func buildGCXiaoHuList(player changsha.Player) *changshapb.GCXiaoHuList {
	xhl := &changshapb.GCXiaoHuList{}

	xhl.XiaoHus = make([]int32, 0, len(player.XiaoHus()))
	for _, xho := range player.XiaoHus() {
		xh := int32(xho.XiaoHuType)
		xhl.XiaoHus = append(xhl.XiaoHus, xh)
	}
	return xhl
}

func buildXiaoHu(xhp *changsha.XiaoHuOperation) *changshapb.XiaoHu {
	xh := &changshapb.XiaoHu{}
	xh.Cards = buildCards(xhp.Cards)
	t := int32(xhp.XiaoHuType)
	xh.XiaoHuType = &t
	return xh
}

func buildPlayerOperations(player changsha.Player) []*changshapb.Operation {
	ops := make([]*changshapb.Operation, 0, len(player.PossibleOperations()))
	for _, pop := range player.PossibleOperations() {
		op := buildOperation(pop)
		ops = append(ops, op)
	}
	return ops
}

func buildOperation(op *changsha.Operation) *changshapb.Operation {
	top := &changshapb.Operation{}
	opType := int32(op.OperationType)
	top.OperationType = &opType
	top.CardList = buildCards(op.Cards)
	targetCardValue := card.Value(op.TargetCard)
	top.TargetCard = &targetCardValue
	targetIndex := int32(op.TargetIndex)
	top.TargetIndex = &targetIndex
	return top
}

func buildHuList(players []changsha.Player) []*changshapb.Hu {
	if len(players) == 0 {
		return nil
	}
	hs := make([]*changshapb.Hu, 0, len(players))
	for _, p := range players {
		h := buildHu(p)
		hs = append(hs, h)
	}
	return hs
}

func buildHu(p changsha.Player) *changshapb.Hu {
	h := &changshapb.Hu{}
	pId := p.Id()
	h.PlayerId = &pId
	result := p.Result()
	fmt.Printf("build hu %d\n", result)
	h.HuType = &result
	return h
}

func buildRoundSettlePlayerList(r *changsha.Room) []*changshapb.RoundSettlePlayerInfo {
	infos := make([]*changshapb.RoundSettlePlayerInfo, 0, len(r.RoomPlayerManager().Players()))
	for i := 0; i < r.RoomPlayerManager().NumPlayers(); i++ {
		p := r.RoomPlayerManager().GetPlayerByPos(i)
		info := buildRoundSettlePlayer(p)
		infos = append(infos, info)
	}
	return infos
}

func buildRoundSettlePlayer(p changsha.Player) *changshapb.RoundSettlePlayerInfo {
	info := &changshapb.RoundSettlePlayerInfo{}
	pId := p.Id()
	info.PlayerId = &pId
	info.ComposeList = buildComposeList(p.ComposeList())
	info.CardList = buildCards(p.Cards())
	cs := p.CurrentScore()
	info.Score = &cs

	xiaoHus := make([]int32, 0, len(p.CurrentXiaoHus()))
	for _, xiaohu := range p.CurrentXiaoHus() {
		xiaoHus = append(xiaoHus, int32(xiaohu.XiaoHuType))
	}
	info.XiaoHuList = xiaoHus
	return info
}

func buildComposeList(composeList []*changsha.Compose) []*changshapb.Compose {
	cl := make([]*changshapb.Compose, 0, len(composeList))
	for _, compose := range composeList {
		c := buildCompose(compose)
		cl = append(cl, c)
	}
	return cl
}

func buildCompose(compose *changsha.Compose) *changshapb.Compose {
	c := &changshapb.Compose{}
	c.Cards = buildCards(compose.CardList)
	cType := int32(compose.ComposeType)
	c.ComposeType = &cType
	return c
}

func buildCards(cs []*card.Card) []int32 {
	cvs := make([]int32, 0, len(cs))
	for _, c := range cs {
		cvs = append(cvs, card.Value(c))
	}
	return cvs
}

func buildGCTotalSettle(r *changsha.Room) *changshapb.GCTotalSettle {
	gcTotalSettle := &changshapb.GCTotalSettle{}
	infos := make([]*changshapb.TotalSettlePlayerInfo, 0, len(r.RoomPlayerManager().Players()))
	for _, p := range r.RoomPlayerManager().Players() {
		info := buildTotalSettlePlayer(p)
		infos = append(infos, info)
	}
	gcTotalSettle.TotalSettlePlayerInfoList = infos
	return gcTotalSettle
}

func buildTotalSettlePlayer(player changsha.Player) *changshapb.TotalSettlePlayerInfo {
	totalSettlePlayerInfo := &changshapb.TotalSettlePlayerInfo{}
	pId := player.Id()
	totalSettlePlayerInfo.PlayerId = &pId
	score := int32(player.Score())
	totalSettlePlayerInfo.Totalscore = &score
	settleInfoList := make([]*changshapb.SettleInfo, 0, len(player.Settles()))
	for settleType, num := range player.Settles() {
		settleInfo := &changshapb.SettleInfo{}
		st := int32(settleType)
		stp := &st
		settleInfo.SettleType = stp
		np := num
		npp := &np
		settleInfo.Num = npp
		settleInfoList = append(settleInfoList, settleInfo)
	}
	totalSettlePlayerInfo.SettleInfoList = settleInfoList
	return totalSettlePlayerInfo
}
