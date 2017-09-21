package changsha

import (
	"sort"

	"game/mahjong/card"
)

type HandCardType int

func (hct HandCardType) String() string {
	r := ""
	for _, i := range handCardTypeSlice {
		if hct&i == i {
			r += handCardTypeMap[i]
			r += ","
		}
	}
	if len(r) > 0 {
		return r[:len(r)-1]
	}
	return r
}

var (
	handCardTypeSlice = []HandCardType{
		HandCardTypePingHu,
		HandCardTypePengPengHu,
		HandCardTypeJiangJiangHu,
		HandCardTypeQingYiSe,
		HandCardTypeHaiDiLaoYue,
		HandCardTypeHaiDiPao,
		HandCardTypeQiXiaoDui,
		HandCardTypeHaoHuaQiXiaoDui,
		HandCardTypeGangShangKaiHua,
		HandCardTypeQiangGangHu,
		HandCardTypeGangShangPao,
		HandCardTypeQuanQiuRen,
		HandCardTypeShuangHaoHuaQiXiaoDui,
		HandCardTypeTianHu,
		HandCardTypeDiHu,
	}
)

var (
	handCardTypeMap = map[HandCardType]string{
		HandCardTypePingHu:                "平胡",
		HandCardTypePengPengHu:            "碰碰胡",
		HandCardTypeJiangJiangHu:          "将将胡",
		HandCardTypeQingYiSe:              "清一色",
		HandCardTypeHaiDiLaoYue:           "海底捞",
		HandCardTypeHaiDiPao:              "海底炮",
		HandCardTypeQiXiaoDui:             "七小对",
		HandCardTypeHaoHuaQiXiaoDui:       "豪华七小对",
		HandCardTypeGangShangKaiHua:       "杠上开花",
		HandCardTypeQiangGangHu:           "抢杠胡",
		HandCardTypeGangShangPao:          "杠上炮",
		HandCardTypeQuanQiuRen:            "全求人",
		HandCardTypeShuangHaoHuaQiXiaoDui: "双豪华七小对",
		HandCardTypeTianHu:                "天胡",
		HandCardTypeDiHu:                  "地胡",
	}
)

const (
	HandCardTypePingHu HandCardType = 1 << iota
	HandCardTypePengPengHu
	HandCardTypeJiangJiangHu
	HandCardTypeQingYiSe
	HandCardTypeHaiDiLaoYue
	HandCardTypeHaiDiPao
	HandCardTypeQiXiaoDui
	HandCardTypeHaoHuaQiXiaoDui
	HandCardTypeGangShangKaiHua
	HandCardTypeQiangGangHu
	HandCardTypeGangShangPao
	HandCardTypeQuanQiuRen
	HandCardTypeShuangHaoHuaQiXiaoDui
	HandCardTypeTianHu
	HandCardTypeDiHu
)

type XiaoHuType int32

var (
	xiaoHuTypeMap = map[XiaoHuType]string{
		XiaoHuTypeDaSiXi:     "大四喜",
		XiaoHuTypeBanBanHu:   "板板胡",
		XiaoHuTypeQueYiSe:    "缺一色",
		XiaoHuTypeLiuLiuShun: "六六顺",
	}
)

func (xht XiaoHuType) String() string {
	return xiaoHuTypeMap[xht]
}

const (
	XiaoHuTypeDaSiXi XiaoHuType = 1 << iota
	XiaoHuTypeBanBanHu
	XiaoHuTypeQueYiSe
	XiaoHuTypeLiuLiuShun
)

func XiaoHu(cl card.CardList) (xiaohus []*XiaoHuOperation) {
	sort.Sort(cl)
	fourConnectXiaoHu := fourConnect(cl)
	if fourConnectXiaoHu != nil {
		xiaohus = append(xiaohus, fourConnectXiaoHu)
	}
	isBbh := isBanBanHu(cl)
	if isBbh {
		bbh := &XiaoHuOperation{}
		bbh.Cards = cl
		bbh.XiaoHuType = XiaoHuTypeBanBanHu
		xiaohus = append(xiaohus, bbh)
	}
	isQueYiSe := isQueYiSe(cl)
	if isQueYiSe {
		qys := &XiaoHuOperation{}
		qys.Cards = cl
		qys.XiaoHuType = XiaoHuTypeQueYiSe
		xiaohus = append(xiaohus, qys)
	}
	lls := liuLiuShun(cl)
	if lls != nil {
		xiaohus = append(xiaohus, lls)
	}
	return
}

//判断是否大四喜
func fourConnect(cl card.CardList) *XiaoHuOperation {
	numConnect := 0
	var tempCard *card.Card
	for _, c := range cl {
		if tempCard != nil && c.Equal(tempCard) {
			numConnect += 1
			if numConnect >= 4 {
				x := &XiaoHuOperation{}
				x.XiaoHuType = XiaoHuTypeDaSiXi
				x.Cards = make([]*card.Card, 0, 4)
				for i := 0; i < 4; i++ {
					x.Cards = append(x.Cards, c)
				}
				return x
			}
			continue
		}
		tempCard = c
		numConnect = 1
	}
	return nil
}

//判断是否板板胡
func isBanBanHu(cl card.CardList) bool {
	for _, c := range cl {
		if c.CardValue == card.CardValueTwo {
			return false
		}
		if c.CardValue == card.CardValueFive {
			return false
		}
		if c.CardValue == card.CardValueEight {
			return false
		}
	}
	return true
}

//判断是否缺一色
func isQueYiSe(cl card.CardList) bool {

	cardTypeMap := make(map[card.CardType]card.CardType)

	for _, c := range cl {
		cardTypeMap[c.CardType] = c.CardType
	}
	if len(cardTypeMap) <= 2 {
		return true
	} else {
		return false
	}
}

//判断是否六六顺
func liuLiuShun(cl card.CardList) *XiaoHuOperation {
	numOfThree := 0
	numConnect := 0
	tempCardValue := int32(-1)
	x := &XiaoHuOperation{}
	x.Cards = make([]*card.Card, 0, 6)
	x.XiaoHuType = XiaoHuTypeLiuLiuShun

	for _, c := range cl {
		if card.Value(c) == tempCardValue {
			numConnect += 1
			continue
		}
		if numConnect >= 3 {
			numOfThree += 1
			for i := 0; i < 3; i++ {
				x.Cards = append(x.Cards, c)
			}
		}
		tempCardValue = card.Value(c)
		numConnect = 1
	}

	if numConnect >= 3 {
		numOfThree += 1
		x.Cards = append(x.Cards, card.NewCardValue(tempCardValue))
	}
	if numOfThree >= 2 {
		return x
	}
	return nil
}

//自己摸 或者别人打
type CardPlayObject int

const (
	//自己摸
	CardPlayObjectSelf = iota
	//别人打
	CardPlayObjectOther
)

var (
	cardPlayerObjectMap = map[CardPlayObject]string{
		CardPlayObjectSelf:  "自己摸",
		CardPlayObjectOther: "别人打",
	}
)

func (cpo CardPlayObject) String() string {
	return cardPlayerObjectMap[cpo]
}

//种类
type CardPlayMethod int

var (
	cardPlayMethodMap = map[CardPlayMethod]string{
		CardPlayMethodNormal:    "普通",
		CardPlayMethodAfterGang: "杠后",
		CardPlayMethodGang:      "踩杠牌",
	}
)

func (cpm CardPlayMethod) String() string {
	return cardPlayMethodMap[cpm]
}

const (
	//普通
	CardPlayMethodNormal = iota
	//杠牌
	CardPlayMethodGang
	//杠后
	CardPlayMethodAfterGang
	//杠后没胡的
	CardPlayMethoAfterGangHu
	//吃碰后
	CardPlayMethodAfterChiPeng
)

//顺序
type CardPlayOrder int

var (
	cardPlayOrderMap = map[CardPlayOrder]string{
		CardPlayOrderNormal: "普通",
		CardPlayOrderFirst:  "起手牌",
		CardPlayOrderLast:   "海底牌",
	}
)

func (cpo CardPlayOrder) String() string {
	return cardPlayOrderMap[cpo]
}

const (

	//顺序方式
	CardPlayOrderNormal = iota
	//起手牌
	CardPlayOrderFirst
	//海底牌
	CardPlayOrderLast
)

//胡牌
func Hu(cards card.CardList, composes ComposeList, c *card.Card, cardPlayObject CardPlayObject, cardPlayMethod CardPlayMethod, cardPlayOrder CardPlayOrder) (result HandCardType) {

	sort.Sort(composes)

	allCards := make([]*card.Card, len(cards), len(cards)+1)
	copy(allCards, cards)
	if cardPlayObject == CardPlayObjectOther || cardPlayMethod == CardPlayMethodAfterGang || cardPlayOrder == CardPlayOrderLast {
		allCards = append(allCards, c)
	}
	sort.Sort(card.CardList(allCards))

	//检查是否将将胡
	for _, tempCard := range allCards {
		if !tempCard.IsJiang() {
			goto AfterJiang
		}
	}

	for _, tempCompose := range composes {
		for _, tempCard := range tempCompose.CardList {
			if !tempCard.IsJiang() {
				goto AfterJiang
			}
		}
	}

	result |= HandCardTypeJiangJiangHu

AfterJiang:
	tGangList, tKeziList, _, tRemains := card.Count(allCards)
	//7对
	if len(tRemains) == 0 && len(tKeziList) == 0 && len(composes) == 0 {
		if len(tGangList) >= 2 {
			result |= HandCardTypeShuangHaoHuaQiXiaoDui
		} else if len(tGangList) == 1 {
			result |= HandCardTypeHaoHuaQiXiaoDui
		} else {
			result |= HandCardTypeQiXiaoDui
		}
	}

	shunziList, gangList, _, duiziList, remains := card.Combine(allCards)
	lenRemains := len(remains)
	if lenRemains != 0 && result == 0 {
		return
	}

	lenShunzi := len(shunziList)
	//lenKezi := len(keziList)
	lenDuizi := len(duiziList)
	lenGang := len(gangList)

	lenComposeShunzi := 0
	lenComposeKezi := 0
	lenComposeGang := 0
	for _, tempCompose := range composes {
		switch tempCompose.ComposeType {
		case ComposeTypeChi:
			lenComposeShunzi += 1
		case ComposeTypeAnGang:
			lenComposeGang += 1
		case ComposeTypeGang:
			lenComposeGang += 1
		case ComposeTypePeng:
			lenComposeKezi += 1
		}
	}

	if (lenShunzi + lenComposeShunzi + lenRemains + lenGang) == 0 {
		//检查是否碰碰胡
		//1对
		if lenDuizi == 1 {
			result |= HandCardTypePengPengHu
		}
	}

	//检查是胡
	//1对或4对 判断是不是组成2个顺子
	if lenRemains == 0 && (lenDuizi == 4 || lenDuizi == 1) {

		//检查是否有将
		hasJiang := false
		// for _, tempCard := range keziList {
		// 	if tempCard.IsJiang() {
		// 		hasJiang = true
		// 		break
		// 	}
		// }
		if !hasJiang {
			for _, tempCard := range duiziList {
				if tempCard.IsJiang() {
					hasJiang = true
					break
				}
			}
		}

		// if !hasJiang {
		// 	for _, tempCompose := range composes {
		// 		switch tempCompose.ComposeType {
		// 		case ComposeTypeAnGang:
		// 			if tempCompose.CardList[0].IsJiang() {
		// 				hasJiang = true
		// 				goto Out
		// 			}
		// 		case ComposeTypeGang:
		// 			if tempCompose.CardList[0].IsJiang() {
		// 				hasJiang = true
		// 				goto Out
		// 			}
		// 		case ComposeTypePeng:
		// 			if tempCompose.CardList[0].IsJiang() {
		// 				hasJiang = true
		// 				goto Out
		// 			}

		// 		}
		// 	Out:
		// 		break
		// 	}

		// }
		// if hasJiang {
		if lenDuizi == 4 {
			if !card.IsStraight(duiziList[:3]) && !card.IsStraight(duiziList[1:]) {
				goto AfterPingHu
			}

		}
		result |= HandCardTypePingHu
		//}
	}
AfterPingHu:
	//检查青一色
	if lenRemains == 0 {
		isSame := true
		var (
			tempCardType card.CardType
		)
		for i, tempCard := range allCards {
			if i == 0 {
				tempCardType = tempCard.CardType
				continue
			}
			if tempCardType != tempCard.CardType {
				isSame = false
				break
			}
		}
		if isSame {
			for _, tempCompose := range composes {
				if tempCompose.CardList[0].CardType != tempCardType {
					isSame = false
					break
				}
			}
			if isSame {
				if result != 0 {
					result |= HandCardTypeQingYiSe
				}
			}
		}

		//检查全求人
		if cardPlayObject == CardPlayObjectOther || cardPlayMethod == CardPlayMethodAfterGang {
			if len(cards) == 1 {
				result |= HandCardTypeQuanQiuRen
			}
		} else {
			if len(cards) == 2 && len(duiziList) == 1 {
				result |= HandCardTypeQuanQiuRen
			}
		}
	}

	//平胡检测是否有将
	if result == HandCardTypePingHu {
		//一个对子
		if len(duiziList) == 1 {
			if duiziList[0].IsJiang() {
				goto Finally
			}
			return 0
		}
		//4个对子
		//检查是否有将
		if (card.IsStraight(duiziList[:3]) && duiziList[3].IsJiang()) || (card.IsStraight(duiziList[1:]) && duiziList[0].IsJiang()) {
			goto Finally
		}

		return 0
	}
	//排除平胡
	result &= ^HandCardTypePingHu
Finally:
	if result != 0 {
		if cardPlayObject == CardPlayObjectOther {
			//检查海底炮
			if cardPlayOrder == CardPlayOrderLast {
				result |= HandCardTypeHaiDiPao
			}

			//检查是否杠上炮
			if cardPlayMethod == CardPlayMethodAfterGang {
				result |= HandCardTypeGangShangPao
			}

			//检查抢杠胡
			if cardPlayMethod == CardPlayMethodGang {
				result |= HandCardTypeGangShangPao
			}

			if cardPlayOrder == CardPlayOrderFirst {
				result |= HandCardTypeDiHu
			}
		} else {
			//检查海底捞月
			if cardPlayOrder == CardPlayOrderLast {
				result |= HandCardTypeHaiDiLaoYue
			}

			//检查杠上开花
			if cardPlayMethod == CardPlayMethodAfterGang {
				result |= HandCardTypeGangShangKaiHua
			}
			if cardPlayOrder == CardPlayOrderFirst {
				result |= HandCardTypeTianHu
			}
		}
	}
	return result

}

type SettleType int32

const (
	SettleTypeDaHuZiMo SettleType = iota
	SettleTypeXiaoHuZiMo
	SettleTypeDaHuDianPao
	SettleTypeXiaoHuDianPao
	SettleTypeDaHuJiePao
	SettleTypeXiaoHuJiePao
)
