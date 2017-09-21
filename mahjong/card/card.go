package card

import (
	"fmt"
	"math"
)

var (
	cardValues = [...]string{
		"",
		"一",
		"二",
		"三",
		"四",
		"五",
		"六",
		"七",
		"八",
		"九",
	}
)

type CardValue int

const (
	CardValueNone CardValue = iota
	CardValueOne
	CardValueTwo
	CardValueThree
	CardValueFour
	CardValueFive
	CardValueSix
	CardValueSeven
	CardValueEight
	CardValueNine
)

func (cv CardValue) String() string {
	return cardValues[cv]
}

var (
	cardTypes = [...]string{
		"万",
		"筒",
		"条",
	}
)

type CardType int

const (
	CardTypeWang CardType = iota
	CardTypeTong
	CardTypeTiao
)

func (cv CardType) String() string {
	return cardTypes[cv]
}

type Card struct {
	CardType  CardType
	CardValue CardValue
}

func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.CardValue, c.CardType)
}

func (c *Card) IsJiang() bool {

	if c.CardValue == CardValueTwo || c.CardValue == CardValueFive || c.CardValue == CardValueEight {

		return true
	}

	return false
}

func (c *Card) Equal(c2 *Card) bool {

	if c.CardType != c2.CardType {
		return false
	}
	if c.CardValue != c2.CardValue {
		return false
	}
	return true
}

func (c *Card) Adjacent(c2 *Card) bool {
	if c.CardType != c2.CardType {
		return false
	}
	if math.Abs(float64(c.CardValue-c2.CardValue)) == 1 {
		return true
	}
	return false
}

func NewCard(ct CardType, cv CardValue) *Card {
	return &Card{
		CardType:  ct,
		CardValue: cv,
	}
}

func NewCardValue(val int32) *Card {
	return &Card{
		CardType:  CardType(val >> 4),
		CardValue: CardValue(val & 15),
	}
}

func Value(c *Card) int32 {
	return int32(c.CardValue) | int32(c.CardType)<<4
}

func Values(cs []*Card) (cvs []int32) {
	for _, c := range cs {
		cvs = append(cvs, Value(c))
	}
	return cvs
}

func NewCardValues(vals []int32) (cs []*Card) {
	for _, val := range vals {
		cs = append(cs, NewCardValue(val))
	}
	return
}
