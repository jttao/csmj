package changsha

import (
	"fmt"

	"game/mahjong/card"
)

//杠类型
type GangType int

var gangTypeMap = map[GangType]string{
	GangTypeGang:      "3张直接杠",
	GangTypeDiPaiGang: "底牌杠",
	GangTypeAnGang:    "暗杠",
	//杠补杠
	GangTypeGangBu: "杠补杠",
}

func (gt GangType) String() string {
	return gangTypeMap[gt]
}

const (
	//3张直接杠
	GangTypeGang GangType = iota
	//底牌杠
	GangTypeDiPaiGang
	//暗杠
	GangTypeAnGang
	//杠补杠
	GangTypeGangBu
)

//组合类型
type ComposeType int

const (
	ComposeTypeChi ComposeType = iota
	ComposeTypePeng
	ComposeTypeGang
	ComposeTypeAnGang
)

var composeTypeMap = map[ComposeType]string{
	ComposeTypeChi:    "吃",
	ComposeTypePeng:   "碰",
	ComposeTypeGang:   "杠",
	ComposeTypeAnGang: "暗杠",
}

func (ct ComposeType) String() string {
	return composeTypeMap[ct]
}

//牌底组合
type Compose struct {
	ComposeType ComposeType
	CardList    card.CardList
}

func (c *Compose) String() string {
	return fmt.Sprintf("%s,牌[%s]", c.ComposeType, c.CardList)
}

//牌底列表
type ComposeList []*Compose

func (cl ComposeList) Len() int {
	return len(cl)
}

func (cl ComposeList) Less(i, j int) bool {
	if cl[i].ComposeType < cl[j].ComposeType {
		return true
	}

	if card.Value(cl[i].CardList[0]) < card.Value(cl[j].CardList[0]) {
		return true
	}

	return false
}

func (cl ComposeList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func NewCompose(composeType ComposeType, cardList card.CardList) *Compose {
	return &Compose{
		ComposeType: composeType,
		CardList:    cardList,
	}
}
