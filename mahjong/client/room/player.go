package room

import (
	"fmt"
	"sort"

	"game/mahjong/card"
	"game/mahjong/changsha"
)

type Player interface {
	Id() int64
	SetId(id int64)
	Position() int
	DealCards(cards []*card.Card)
	Cards() []*card.Card
	PlayCard(cv int32) *card.Card
	Chi(cv int32, cvs []int32) *changsha.Compose
	Peng(cv int32) *changsha.Compose
	Mo()
	MoCard(c *card.Card)
	Clear()
	Bu(cv int32, buType int32) *changsha.Compose
	Gang(cv int32, gangType int32) *changsha.Compose
	SetCardsNum(cardsNum int32)
	SetCards(cs []*card.Card)
	SetPlayedCards(cs []*card.Card)
	SetComposes(cs []*changsha.Compose)
	Location() string
	SetLocation(location string)
}

type player struct {
	id          int64
	position    int
	cardsNum    int32
	cards       []*card.Card
	playedCards []*card.Card
	composes    []*changsha.Compose
	location 	string
}

func (p *player) Location() string {
	return p.location
}

func (p *player) SetLocation(location string) {
	p.location = location
}

func (p *player) Id() int64 {
	return p.id
}

func (p *player) SetId(id int64) {
	p.id = id
}

func (p *player) Position() int {
	return p.position
}

func (p *player) SetCardsNum(cardsNum int32) {
	p.cardsNum = cardsNum
}

func (p *player) SetCards(cs []*card.Card) {
	p.cards = cs
}

func (p *player) SetPlayedCards(cs []*card.Card) {
	p.playedCards = cs
}

func (p *player) SetComposes(cs []*changsha.Compose) {
	p.composes = cs
}

func (p *player) DealCards(cards []*card.Card) {
	p.cards = append(p.cards, cards...)
}

func (p *player) Cards() []*card.Card {
	return p.cards
}

func (p *player) PlayCard(cv int32) *card.Card {
	if !p.isClient() {
		c := card.NewCardValue(cv)
		p.playedCards = append(p.playedCards, c)
		return c
	}
	index := -1
	for i, c := range p.cards {
		if card.Value(c) == cv {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	c := p.cards[index]
	p.playedCards = append(p.playedCards, p.cards[index])
	p.cards = append(p.cards[:index], p.cards[index+1:]...)
	return c

}

func (p *player) Chi(cv int32, cvs []int32) *changsha.Compose {
	compose := &changsha.Compose{}
	compose.ComposeType = changsha.ComposeTypeChi
	compose.CardList = make([]*card.Card, 0, 3)
	p.composes = append(p.composes, compose)
	if !p.isClient() {

		compose.CardList = append(compose.CardList, card.NewCardValue(cv))
		for _, tcv := range cvs {
			compose.CardList = append(compose.CardList, card.NewCardValue(tcv))
		}

		return compose
	}
	cs := p.takeCards(cvs)

	compose.CardList = append(compose.CardList, cs...)
	compose.CardList = append(compose.CardList, card.NewCardValue(cv))
	return compose
}

func (p *player) Peng(cv int32) *changsha.Compose {
	compose := &changsha.Compose{}
	p.composes = append(p.composes, compose)
	compose.ComposeType = changsha.ComposeTypePeng
	compose.CardList = make([]*card.Card, 0, 3)
	if !p.isClient() {
		for i := 0; i < 3; i++ {
			compose.CardList = append(compose.CardList, card.NewCardValue(cv))
		}
		return compose
	}
	cs := p.takeNCard(cv, 2)
	compose.CardList = append(compose.CardList, card.NewCardValue(cv))
	compose.CardList = append(compose.CardList, cs...)
	return compose
}

func (p *player) Gang(cv int32, buType int32) *changsha.Compose {
	switch buType {
	//明杠
	case 0:
		{
			compose := &changsha.Compose{}
			compose.ComposeType = changsha.ComposeTypeGang
			compose.CardList = make([]*card.Card, 0, 4)
			if !p.isClient() {
				for i := 0; i < 4; i++ {
					compose.CardList = append(compose.CardList, card.NewCardValue(cv))
				}
				//少3张牌
				return compose
			}
			cs := p.takeNCard(cv, 3)
			compose.CardList = append(compose.CardList, card.NewCardValue(cv))
			compose.CardList = append(compose.CardList, cs...)
			return compose
		}
		//摸的明杠
	case 1, 3:
		{
			fmt.Println("compose %s", p.composes)
			for _, c := range p.composes {
				if c.ComposeType == changsha.ComposeTypePeng {
					if card.Value(c.CardList[0]) == cv {
						if p.isClient() {
							p.takeNCard(cv, 1)
						}
						c.ComposeType = changsha.ComposeTypeGang
						c.CardList = append(c.CardList, card.NewCardValue(cv))
						return c
					}
				}
			}
		}
		//暗杠
	case 2:
		{
			compose := &changsha.Compose{}
			compose.ComposeType = changsha.ComposeTypeAnGang
			compose.CardList = make([]*card.Card, 0, 4)
			if !p.isClient() {
				for i := 0; i < 4; i++ {
					compose.CardList = append(compose.CardList, card.NewCardValue(cv))
				}
				//少4张牌
				return compose
			}
			cs := p.takeNCard(cv, 4)
			compose.CardList = append(compose.CardList, cs...)
			return compose
		}

	}
	fmt.Println("gang type ", buType)
	panic("never reach here")
}

func (p *player) Bu(cv int32, buType int32) *changsha.Compose {

	switch buType {
	//明杠
	case 0:
		{
			compose := &changsha.Compose{}
			compose.ComposeType = changsha.ComposeTypeGang
			compose.CardList = make([]*card.Card, 0, 4)
			if !p.isClient() {
				for i := 0; i < 4; i++ {
					compose.CardList = append(compose.CardList, card.NewCardValue(cv))
				}
				//少3张牌
				return compose
			}
			cs := p.takeNCard(cv, 3)
			compose.CardList = append(compose.CardList, card.NewCardValue(cv))
			compose.CardList = append(compose.CardList, cs...)
			return compose
		}
		//摸的明杠
	case 1, 3:
		{
			fmt.Println("compose %s", p.composes)
			for _, c := range p.composes {
				if c.ComposeType == changsha.ComposeTypePeng {
					if card.Value(c.CardList[0]) == cv {
						if p.isClient() {
							p.takeNCard(cv, 1)
						}
						c.ComposeType = changsha.ComposeTypeGang
						c.CardList = append(c.CardList, card.NewCardValue(cv))
						return c
					}
				}
			}
		}
		//暗杠
	case 2:
		{
			compose := &changsha.Compose{}
			compose.ComposeType = changsha.ComposeTypeAnGang
			compose.CardList = make([]*card.Card, 0, 4)
			if !p.isClient() {
				for i := 0; i < 4; i++ {
					compose.CardList = append(compose.CardList, card.NewCardValue(cv))
				}
				//少4张牌
				return compose
			}
			cs := p.takeNCard(cv, 4)
			compose.CardList = append(compose.CardList, cs...)
			return compose
		}

	}
	fmt.Println("gang type ", buType)
	panic("never reach here")
}

func (p *player) Mo() {
	p.cardsNum += 1
}

func (p *player) MoCard(c *card.Card) {
	p.cards = append(p.cards, c)
	sort.Sort(card.CardList(p.cards))
}

func (p *player) takeCard(cv int32) (c *card.Card) {
	index := -1
	for i, c := range p.cards {
		if card.Value(c) == cv {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	c = p.cards[index]
	p.playedCards = append(p.playedCards, p.cards[index])
	p.cards = append(p.cards[:index], p.cards[index+1:]...)
	return c
}

func (p *player) takeNCard(cv int32, n int) (cs []*card.Card) {
	for i := 0; i < n; i++ {
		c := p.takeCard(cv)
		cs = append(cs, c)
	}
	return
}

func (p *player) takeCards(cvs []int32) (cs []*card.Card) {

	for _, cv := range cvs {
		c := p.takeCard(cv)
		cs = append(cs, c)
	}

	return cs
}

func (p *player) Clear() {
	p.cardsNum = 0
	p.cards = nil
	p.playedCards = nil
	p.composes = nil
}

func (p *player) isClient() bool {
	return len(p.cards) != 0
}

func NewPlayer(id int64, position int) Player {
	p := &player{}
	p.id = id
	p.position = position 
	return p
}

type PlayerManager interface {
	AddPlayer(p Player) bool
	RemovePlayer(p Player) bool
	Players() []Player
	GetPlayerByPos(pos int) Player
	GetPlayerById(id int64) Player
}

type playerManager struct {
	players []Player
}

func (pm *playerManager) Players() []Player {
	return pm.players
}

func (pm *playerManager) AddPlayer(p Player) bool {
	pm.players = append(pm.players, p)
	return true
}
func (pm *playerManager) RemovePlayer(p Player) bool {
	index := -1
	for i, tp := range pm.players {
		if tp.Id() == p.Id() {
			index = i
			break
		}
	}
	if index == -1 {
		return true
	}
	pm.players = append(pm.players[:index], pm.players[index+1:]...)
	return true
}

func (pm *playerManager) GetPlayerByPos(pos int) Player {
	for _, p := range pm.players {
		if p.Position() == pos {
			return p
		}
	}
	return nil
}

func (pm *playerManager) GetPlayerById(id int64) Player {
	for _, p := range pm.players {
		if p.Id() == id {
			return p
		}
	}
	return nil
}

func NewPlayerManager() PlayerManager {
	pm := &playerManager{}
	return pm
}
