package handler

import (
	"game/mahjong/card"
	changsha "game/mahjong/changsha"
	"game/mahjong/client/room"
	changshapb "game/mahjong/pb/changsha"
)

func playerFromPlayerInfo(gcPlayerInfo *changshapb.GCPlayerInfo) room.Player {
	pl := room.NewPlayer(*gcPlayerInfo.PlayerId, int(*gcPlayerInfo.Position))
	pl.SetCardsNum(gcPlayerInfo.GetCardNum())

	pl.SetCards(cardsFromInts(gcPlayerInfo.GetCardList()))
	pl.SetComposes(composeListFromComposeInfoList(gcPlayerInfo.GetComposeList()))
	pl.SetPlayedCards(cardsFromInts(gcPlayerInfo.GetPlayCardList()))
	return pl
}

func roomFromRoomInfo(gcRoomInfo *changshapb.GCRoomInfo) room.Room {
	r := room.NewRoom(*gcRoomInfo.RoomId)
	return r
}

func cardsFromInts(vals []int32) (cs []*card.Card) {
	for _, v := range vals {
		c := card.NewCardValue(v)
		cs = append(cs, c)
	}
	return
}

func xiaoHuOpreationsFromXiaoHu(xiaohu *changshapb.XiaoHu) *changsha.XiaoHuOperation {
	xhp := &changsha.XiaoHuOperation{}
	xhp.XiaoHuType = changsha.XiaoHuType(xiaohu.GetXiaoHuType())
	xhp.Cards = cardsFromInts(xiaohu.GetCards())
	return xhp
}

func composeListFromComposeInfoList(composeList []*changshapb.Compose) []*changsha.Compose {
	cs := make([]*changsha.Compose, 0, len(composeList))
	for _, compose := range composeList {
		c := composeFromComposeInfo(compose)
		cs = append(cs, c)
	}
	return cs
}
func composeFromComposeInfo(compose *changshapb.Compose) *changsha.Compose {
	c := &changsha.Compose{}
	c.ComposeType = changsha.ComposeType(compose.GetComposeType())
	c.CardList = cardsFromInts(compose.GetCards())
	return c
}
