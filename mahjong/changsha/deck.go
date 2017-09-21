package changsha

import (
	"game/mahjong/card"
	"math/rand"
	"time"
)

const (
	totalCardNum = 108
)

//牌组
type Deck interface {
	Shuffle()
	GetFirst() *card.Card
	GetLast() *card.Card
	Remains() int
}

type deck struct {
	//未使用的牌
	source []*card.Card
	//已经发的牌
	dest []*card.Card
}

func (d *deck) init() {
	d.source = make([]*card.Card, 0, totalCardNum)
	d.dest = make([]*card.Card, 0, totalCardNum)

	for i := card.CardTypeWang; i <= card.CardTypeTiao; i++ {
		for j := card.CardValueOne; j <= card.CardValueNine; j++ {
			for k := 0; k < 4; k++ {
				c := card.NewCard(i, j)
				d.source = append(d.source, c)
			}
		}
	}
}

func (d *deck) Shuffle() {
	d.source = append(d.source, d.dest...)
	d.dest = d.dest[0:0]
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(d.source))
	for i, v := range perm {
		d.source[i], d.source[v] = d.source[v], d.source[i]
	}
}

func (d *deck) GetFirst() *card.Card {
	c := d.source[0]
	d.source = d.source[1:]
	d.dest = append(d.dest, c)
	return c
}

func (d *deck) GetLast() *card.Card {
	c := d.source[len(d.source)-1]
	d.source = d.source[:len(d.source)-1]
	d.dest = append(d.dest, c)
	return c
}

func (d *deck) Remains() int {
	return len(d.source)
}

func NewDeck() Deck {
	d := &deck{}
	d.init()
	return d
}
