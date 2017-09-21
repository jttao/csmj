package room

import (
	"log"

	"game/mahjong/card"
	"game/mahjong/changsha"
)

type Room interface {
	PlayerManager() PlayerManager
	PlayerLeave(p Player)
	PlayerJoin(p Player)
	BankerSelect(pos int32)
	DealCards(playerId int64, cards []*card.Card)
	XiaoHu(playerId int64, operations []*changsha.XiaoHuOperation)
	WaitPlayerPlay(pId int64)
	PlayerPlay(pId int64, cv int32)
	PlayerChi(pId int64, cv int32, cvs []int32)
	PlayerPeng(pId int64, cv int32)
	PlayerBu(playerId int64, cv int32, buType int32)
	PlayerMoCard(pId int64, cv int32)
	PlayerGangMo(pId int64, cvs []int32)
	PlayerMo(pId int64)
	PlayerGang(pId int64, cv int32, buType int32)
	Clear()
}

type room struct {
	roomId        int64
	bankerPos     int32
	playerManager PlayerManager
}

func (r *room) PlayerManager() PlayerManager {
	return r.playerManager
}
func (r *room) PlayerJoin(p Player) {
	flag := r.playerManager.AddPlayer(p)
	if !flag {
		log.Println("join player failed")
		return
	}
	log.Printf("join player %#v\n", p)
}

func (r *room) PlayerLeave(p Player) {
	flag := r.playerManager.RemovePlayer(p)
	if !flag {
		log.Println("remove player failed")
		return
	}
	log.Printf("remove player %#v\n", p)
}

func (r *room) BankerSelect(pos int32) {

	r.bankerPos = pos
	log.Printf("select banker pos %d", pos)
}

func (r *room) DealCards(playerId int64, cards []*card.Card) {
	log.Printf("player [%d] deal cards [%v]\n", playerId, cards)
	pl := r.PlayerManager().GetPlayerById(playerId)
	pl.DealCards(cards)
}

func (r *room) XiaoHu(playerId int64, operations []*changsha.XiaoHuOperation) {
	log.Printf("玩家[%d]小胡\n", playerId)
	for _, o := range operations {
		log.Println(o)
	}
}

func (r *room) WaitPlayerPlay(playerId int64) {
	log.Printf("等候玩家[%d] 打牌\n", playerId)
}

func (r *room) PlayerPlay(playerId int64, cv int32) {
	pl := r.playerManager.GetPlayerById(playerId)
	c := pl.PlayCard(cv)
	log.Printf("玩家［%d］打牌[%s]", playerId, c)
}

func (r *room) PlayerChi(playerId int64, cv int32, cvs []int32) {
	pl := r.playerManager.GetPlayerById(playerId)
	c := pl.Chi(cv, cvs)
	log.Printf("玩家［%d］吃 [%s]", c)

}

func (r *room) PlayerPeng(playerId int64, cv int32) {
	pl := r.playerManager.GetPlayerById(playerId)
	c := pl.Peng(cv)
	log.Printf("玩家[%d]碰 [%s]", playerId, c)
}

func (r *room) PlayerBu(playerId int64, cv int32, buType int32) {

	pl := r.playerManager.GetPlayerById(playerId)
	c := pl.Bu(cv, buType)
	log.Printf("玩家[%d]补[%s]", playerId, c)
}

func (r *room) PlayerMoCard(playerId int64, cv int32) {
	pl := r.playerManager.GetPlayerById(playerId)
	c := card.NewCardValue(cv)
	pl.MoCard(c)
	log.Printf("玩家[%d]摸牌[%s]", playerId, c)
}

func (r *room) PlayerMo(playerId int64) {
	pl := r.playerManager.GetPlayerById(playerId)
	pl.Mo()
}

func (r *room) PlayerGang(playerId int64, cv int32, buType int32) {
	pl := r.playerManager.GetPlayerById(playerId)
	c := pl.Gang(cv, buType)
	log.Printf("玩家[%d]杠[%s]", playerId, c)
}

func (r *room) PlayerGangMo(playerId int64, cv []int32) {
	log.Printf("玩家[%d]杠摸[%s]", playerId, card.NewCardValues(cv))
}

func (r *room) Clear() {
	for _, pl := range r.playerManager.Players() {
		pl.Clear()
	}
}

func NewRoom(roomId int64) Room {
	r := &room{}
	r.roomId = roomId
	r.playerManager = NewPlayerManager()
	return r
}
