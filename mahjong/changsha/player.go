package changsha

import (
	"fmt"
	"sort"

	"game/mahjong/card"
	"game/mahjong/server/player"
)

const (
	initialScore = 0
)

//操作类型
type OperationType int

const (
	//吃
	OperationTypeEat OperationType = iota
	//碰
	OperationTypePeng
	//杠
	OperationTypeGang
	//补
	OperationTypeBu
	//胡
	OperationTypeHu
	//过
	OperationTypePass
)

var operationTypeMap = map[OperationType]string{
	OperationTypeEat:  "吃",
	OperationTypePeng: "碰",
	OperationTypeGang: "杠",
	OperationTypeBu:   "补",
	OperationTypeHu:   "胡",
	OperationTypePass: "过",
}

var operationTypePriorityMap = map[OperationType]int{
	OperationTypeEat:  1,
	OperationTypePeng: 2,
	OperationTypeGang: 2,
	OperationTypeBu:   2,
	OperationTypeHu:   3,
	OperationTypePass: 0,
}

func (ot OperationType) String() string {
	return operationTypeMap[ot]
}

func (ot OperationType) Priority() int {
	return operationTypePriorityMap[ot]
}

type Operation struct {
	OperationType OperationType
	TargetCard    *card.Card
	Cards         []*card.Card
	TargetIndex   int
}

type OperationList []*Operation

func (ol OperationList) Len() int {
	return len(ol)
}

func (ol OperationList) Less(i, j int) bool {
	if ol[i].OperationType.Priority() < ol[j].OperationType.Priority() {
		return true
	}
	if ol[i].OperationType.Priority() == ol[j].OperationType.Priority() {
		if ol[i].TargetIndex > ol[j].TargetIndex {
			return true
		}
	}
	return false
}

func (ol OperationList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
}

func (op *Operation) Equal(targetIndex int, ot OperationType, targetCardValue int32, cardsValue []int32) bool {
	if op.TargetIndex != targetIndex {
		return false
	}
	if op.OperationType != ot {
		return false
	}
	if card.Value(op.TargetCard) != targetCardValue {
		return false
	}
	if len(op.Cards) != len(cardsValue) {
		return false
	}

	for i := 0; i < len(op.Cards); i++ {
		if card.Value(op.Cards[i]) != cardsValue[i] {
			return false
		}
	}
	return true
}

func (op *Operation) String() string {
	return fmt.Sprintf("操作[%s],卡牌[%s]", op.OperationType, op.Cards)
}

type XiaoHuOperation struct {
	XiaoHuType XiaoHuType
	Cards      []*card.Card
}

func (xho *XiaoHuOperation) String() string {
	return fmt.Sprintf("小胡[%s],卡牌[%s]", xho.XiaoHuType, xho.Cards)
}

type PlayerState int

const (
	PlayerStateOnline = iota
	PlayerStateOffline
	PlayerStateLeave
)

var (
	playerStateMap = map[PlayerState]string{
		PlayerStateOnline:  "线上",
		PlayerStateOffline: "离线",
		PlayerStateLeave:   "离开",
	}
)

func (ps PlayerState) String() string {
	return playerStateMap[ps]
}

type PlayerActualState int

const (
	PlayerActualStateInit PlayerActualState = iota
	PlayerActualStateStart
	PlayerActualStatePlay
	PlayerActualStateSettle
	PlayerActualStateEnd
)

var (
	playerActualStateMap = map[PlayerActualState]string{
		PlayerActualStateInit:   "初始化",
		PlayerActualStateStart:  "准备开始",
		PlayerActualStatePlay:   "正在玩",
		PlayerActualStateSettle: "结算",
	}
)

func (ps PlayerActualState) String() string {
	return playerActualStateMap[ps]
}

type Player interface {
	Id() int64
	Name() string
	Ip() string
	Image() string
	Sex() int
	Player() player.Player
	SetPosition(pos int)
	Position() int
	SetState(s PlayerState)
	State() PlayerState
	SetActualState(s PlayerActualState)
	ActualState() PlayerActualState
	Cards() []*card.Card
	PlayedCards() []*card.Card
	CurrentXiaoHus() []*XiaoHuOperation
	XiaoHus() []*XiaoHuOperation
	RemainXiaoHus() []*XiaoHuOperation
	IfFinishXiaoHu() bool
	Settles() map[SettleType]int32
	ComplementXiaoHu() int
	PlayCard(cardValue int32) *card.Card
	AddPlayedCard(c *card.Card)
	RemovePlayedCard(c *card.Card)
	XiaoHuO(xht XiaoHuType) *XiaoHuOperation
	PassXiaoHu()
	Eat(cardValues []int32, c *card.Card) []*card.Card
	Peng(c *card.Card) []*card.Card
	Bu(c *card.Card, cpo CardPlayObject, cpm CardPlayMethod) (bool, GangType)
	Gang(c *card.Card, cpo CardPlayObject, cpm CardPlayMethod) (bool, GangType)
	Mo(c *card.Card)
	GangMo(cs []*card.Card)
	Hu(c *card.Card, cardPlayObject CardPlayObject, cardPlayMethod CardPlayMethod, cardPlayOrder CardPlayOrder) int32
	ComplementOperations(cs []*card.Card, cardPlayObject CardPlayObject, cardPlayMethod CardPlayMethod, cardPlayOrder CardPlayOrder, next bool, hasCards bool)
	Operate(targetIndex int32, operationType OperationType, targetCardValue int32, cardValues []int32) bool
	CurrentOperation(targetIndex int32, c *card.Card) *Operation
	CurrentOperations() []*Operation
	PossibleOperations() []*Operation
	HasOperations(targetIndex int32, c *card.Card) bool
	ClearOperations()
	Clear()
	//	IsHu() bool
	IsListen() bool
	SetListen(l bool)
	IsOperate(targetIndex int32, c *card.Card) bool
	Result() int32
	ComposeList() []*Compose
	CurrentScore() int64
	Score() int64
	Disconnect()
	Reconnect(p player.Player)
	LevaeTime(state bool)  
	IfPrepare() bool
	IfCanPrepare() bool
	IfCanXiaoHu() bool
	Prepare()
	Start()
	Settle()
	End()
	ZhuaNiao()
	ZhongNiao() int32
	AddScore(s int)
	SubScore(s int)
	AddSettleData(st SettleType)
	IfOperationsFinish(cs []*card.Card) bool
	HasRemainOperations(cs []*card.Card, ot OperationType) bool
	MaxPriorityOpreations(cs []*card.Card) []*Operation
	SetLocation(location string)
	Location() string
}

type roomPlayer struct {
	id                 int64
	ip                 string
	name               string
	image              string
	sex                int
	player             player.Player
	position           int
	state              PlayerState
	actualState        PlayerActualState
	playedCards        []*card.Card
	cards              []*card.Card
	composes           []*Compose
	possibleOperations []*Operation
	currentOpertaions  []*Operation
	currentXiaoHus     []*XiaoHuOperation
	xiaoHus            []*XiaoHuOperation
	xiaoHuPass         bool
	result             int32
	currentScore       int64
	score              int64
	listen             bool
	settles            map[SettleType]int32
	
	//中鸟
	zhongNiao int32
	location string 
}

func NewPlayer(id int64, ip string, name string, image string, sex int,location string, player player.Player) Player {
	p := &roomPlayer{
		id:     id,
		ip:     ip,
		player: player,
		name:   name,
		image:  image,
		sex:    sex,
		location:location, 
	}
	p.init()

	return p
}

//初始化
func (rp *roomPlayer) init() {
	rp.score = initialScore
	rp.settles = make(map[SettleType]int32)
	for i := SettleTypeXiaoHuZiMo; i <= SettleTypeXiaoHuJiePao; i++ {
		rp.settles[i] = 0
	}
}

func (rp *roomPlayer) Id() int64 {
	return rp.id
}

func (rp *roomPlayer) Ip() string {
	return rp.ip
}

func (rp *roomPlayer) Name() string {
	return rp.name
}

func (rp *roomPlayer) Sex() int {
	return rp.sex
}
func (rp *roomPlayer) Image() string {
	return rp.image
}

func (rp *roomPlayer) Player() player.Player {
	return rp.player
}

func (rp *roomPlayer) SetPosition(pos int) {
	rp.position = pos
}

func (rp *roomPlayer) Position() int {
	return rp.position
}

func (rp *roomPlayer) SetLocation(location string) {
	rp.location = location
}

func (rp *roomPlayer) Location() string {
	return rp.location
}

func (rp *roomPlayer) Cards() []*card.Card {
	return rp.cards
}

func (rp *roomPlayer) PlayedCards() []*card.Card {
	return rp.playedCards
}

func (rp *roomPlayer) SetState(s PlayerState) {
	rp.state = s
}

func (rp *roomPlayer) SetActualState(s PlayerActualState) {
	rp.actualState = s
}

func (rp *roomPlayer) Disconnect() {
	rp.player = nil
	rp.state = PlayerStateOffline
}

func (rp *roomPlayer) LevaeTime(state bool) { 
	if state {
		rp.state = PlayerStateLeave
	}else{
		rp.state = PlayerStateOnline
	} 
}

func (rp *roomPlayer) Reconnect(p player.Player) {
	rp.player = p
	rp.state = PlayerStateOnline
}

func (rp *roomPlayer) State() PlayerState {
	return rp.state
}

func (rp *roomPlayer) ActualState() PlayerActualState {
	return rp.actualState
}

func (rp *roomPlayer) ComplementXiaoHu() int {
	rp.xiaoHus = XiaoHu(rp.cards)
	return len(rp.xiaoHus)
}

func (rp *roomPlayer) SetListen(l bool) {
	rp.listen = l
}

func (rp *roomPlayer) IsListen() bool {
	return rp.listen
}

func (rp *roomPlayer) CurrentOperation(targetIndex int32, tempCard *card.Card) *Operation {
	for _, op := range rp.currentOpertaions {
		if op.TargetIndex == int(targetIndex) && op.TargetCard.Equal(tempCard) {
			return op
		}
	}
	return nil
}

func (rp *roomPlayer) CurrentOperations() []*Operation {
	return rp.currentOpertaions
}

func (rp *roomPlayer) PossibleOperations() []*Operation {
	return rp.possibleOperations
}

func (rp *roomPlayer) CurrentXiaoHus() []*XiaoHuOperation {
	return rp.currentXiaoHus
}

func (rp *roomPlayer) RemainXiaoHus() []*XiaoHuOperation {
	if !rp.IfCanXiaoHu() {
		return nil
	}
	remianXiaoHus := make([]*XiaoHuOperation, 0, 4)
	for _, xho := range rp.xiaoHus {
		found := false
		for _, txho := range rp.currentXiaoHus {
			if txho.XiaoHuType == xho.XiaoHuType {
				found = true
				break
			}
		}
		if found {
			continue
		}
		remianXiaoHus = append(remianXiaoHus, xho)
	}
	return remianXiaoHus
}

func (rp *roomPlayer) XiaoHus() []*XiaoHuOperation {
	return rp.xiaoHus
}

func (rp *roomPlayer) Settles() map[SettleType]int32 {
	return rp.settles
}

func (rp *roomPlayer) Result() int32 {
	return rp.result
}

func (rp *roomPlayer) ComposeList() []*Compose {
	return rp.composes
}
func (rp *roomPlayer) CurrentScore() int64 {
	return rp.currentScore
}
func (rp *roomPlayer) Score() int64 {
	return rp.score
}

func (rp *roomPlayer) PlayCard(cardValue int32) *card.Card {
	index := -1
	for i, c := range rp.cards {
		if card.Value(c) == cardValue {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	c := rp.cards[index]
	//rp.playedCards = append(rp.playedCards, c)
	rp.cards = append(rp.cards[:index], rp.cards[index+1:]...)
	return c
}

func (rp *roomPlayer) AddPlayedCard(c *card.Card) {
	rp.playedCards = append(rp.playedCards, c)
}

func (rp *roomPlayer) RemovePlayedCard(c *card.Card) {
	index := -1
	for i := len(rp.playedCards) - 1; i >= 0; i-- {
		if rp.playedCards[i].Equal(c) {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}

	rp.playedCards = append(rp.playedCards[:index], rp.playedCards[index+1:]...)
}

//小胡
func (rp *roomPlayer) XiaoHuO(xht XiaoHuType) *XiaoHuOperation {
	if !rp.IfCanXiaoHu() {
		return nil
	}
	//已经小胡过
	for _, xho := range rp.currentXiaoHus {
		if xho.XiaoHuType == xht {
			return nil
		}
	}
	for _, xho := range rp.xiaoHus {
		if xho.XiaoHuType == xht {
			rp.currentXiaoHus = append(rp.currentXiaoHus, xho)
			return xho
		}
	}
	return nil
}

func (rp *roomPlayer) PassXiaoHu() {
	if !rp.IfCanXiaoHu() {
		return
	}
	rp.xiaoHuPass = true
}

//摸牌
func (rp *roomPlayer) Mo(c *card.Card) {
	rp.cards = append(rp.cards, c)
	sort.Sort(card.CardList(rp.cards))
}

//杠摸
func (rp *roomPlayer) GangMo(cs []*card.Card) {
	rp.playedCards = append(rp.playedCards, cs...)
}

//吃牌
func (rp *roomPlayer) Eat(cardValues []int32, c *card.Card) []*card.Card {
	if len(cardValues) != 2 {
		return nil
	}

	//swap
	if cardValues[0] > cardValues[1] {
		cardValues[0], cardValues[1] = cardValues[1], cardValues[0]
	}

	var (
		tempCardList  []*card.Card
		tempIndexList []int
	)

	index := 0
	for i, tempCard := range rp.cards {
		if card.Value(tempCard) == cardValues[index] {
			tempCardList = append(tempCardList, tempCard)
			tempIndexList = append(tempIndexList, i)
			index++
			if index >= len(cardValues) {
				break
			}
		}
	}

	tempCardList = append(tempCardList, c)
	if len(tempCardList) != 3 {
		return nil
	}

	sort.Sort(card.CardList(tempCardList))

	rp.cards = append(rp.cards[0:tempIndexList[0]], rp.cards[tempIndexList[0]+1:]...)
	rp.cards = append(rp.cards[0:tempIndexList[1]-1], rp.cards[tempIndexList[1]:]...)

	compose := NewCompose(ComposeTypeChi, card.CardList(tempCardList))
	rp.composes = append(rp.composes, compose)
	return tempCardList
}

//碰牌
func (rp *roomPlayer) Peng(c *card.Card) []*card.Card {
	num := 2
	find, startIndex := rp.findEqual(num, c)
	if !find {
		return nil
	}
	tempCardList := make([]*card.Card, 0, 3)
	tempCardList = append(tempCardList, c)
	tempCardList = append(tempCardList, rp.cards[startIndex:startIndex+num]...)
	rp.cards = append(rp.cards[:startIndex], rp.cards[startIndex+num:]...)

	compose := NewCompose(ComposeTypePeng, card.CardList(tempCardList))
	rp.composes = append(rp.composes, compose)
	return tempCardList
}

//补牌
func (rp *roomPlayer) Bu(c *card.Card, cpo CardPlayObject, cpm CardPlayMethod) (bool, GangType) {
	return rp.gang(c, cpo, cpm)
}

//杠牌
func (rp *roomPlayer) Gang(c *card.Card, cpo CardPlayObject, cpm CardPlayMethod) (bool, GangType) {
	return rp.gang(c, cpo, cpm)
}

//杠
func (rp *roomPlayer) gang(c *card.Card, cpo CardPlayObject, cpm CardPlayMethod) (flag bool, gt GangType) {
	//判断是不是杠后
	if cpm == CardPlayMethodAfterGang {
		goto FindCompose
	}
	//查找手里四张牌
	if cpo == CardPlayObjectSelf {
		num := 4
		flag, startIndex := rp.findEqual(num, c)
		if !flag {
			goto FindCompose
		}
		tempCardList := make([]*card.Card, 0, 4)
		tempCardList = append(tempCardList, rp.cards[startIndex:startIndex+num]...)
		rp.cards = append(rp.cards[:startIndex], rp.cards[startIndex+num:]...)
		compose := NewCompose(ComposeTypeAnGang, card.CardList(tempCardList))
		rp.composes = append(rp.composes, compose)
		return true, GangTypeAnGang
	} else {
		//查找手里三张牌
		num := 3
		flag, startIndex := rp.findEqual(num, c)
		if !flag {
			return false, GangTypeGang
		}
		tempCardList := make([]*card.Card, 0, 4)
		tempCardList = append(tempCardList, rp.cards[startIndex:startIndex+num]...)
		rp.cards = append(rp.cards[:startIndex], rp.cards[startIndex+num:]...)
		tempCardList = append(tempCardList, c)
		compose := NewCompose(ComposeTypeGang, card.CardList(tempCardList))
		rp.composes = append(rp.composes, compose)
		return true, GangTypeGang
	}
FindCompose:
	for _, compose := range rp.composes {
		if compose.ComposeType == ComposeTypePeng {
			if compose.CardList[0].Equal(c) {
				compose.ComposeType = ComposeTypeGang
				compose.CardList = append(compose.CardList, c)
				if cpm == CardPlayMethodAfterGang {
					return true, GangTypeGangBu
				}
				flag, startIndex := rp.findEqual(1, c)
				if !flag {
					return false, GangTypeGang
				}
				rp.cards = append(rp.cards[:startIndex], rp.cards[startIndex+1:]...)
				return true, GangTypeDiPaiGang
			}
		}
	}
	return true, GangTypeGang
}

//是否胡牌
// func (rp *roomPlayer) IsHu() bool {
// 	op := rp.CurrentOperation()
// 	if op == nil {
// 		return false
// 	}
// 	return op.OperationType == OperationTypeHu
// }

//判断是否操作完成
func (rp *roomPlayer) IsOperate(targetIndex int32, c *card.Card) bool {
	//判断是否有卡牌操作
	if !rp.HasOperations(targetIndex, c) {
		return true
	}
	//判断 是否操作过
	for _, top := range rp.currentOpertaions {
		if top.TargetCard.Equal(c) && top.TargetIndex == int(targetIndex) {
			return true
		}
	}
	return false
}

//判断是否所有操作完成
func (rp *roomPlayer) IfOperationsFinish(cs []*card.Card) bool {
	for i := 0; i < len(cs); i++ {
		c := cs[i]
		if !rp.HasOperations(int32(i), c) {
			continue
		}
		op := rp.CurrentOperation(int32(i), c)
		if op == nil {
			return false
		}
	}
	return true
}

//判断是否所有操作完成
func (rp *roomPlayer) HasRemainOperations(cs []*card.Card, ot OperationType) bool {
	for i := 0; i < len(cs); i++ {
		c := cs[i]
		//没有操作
		if !rp.HasOperations(int32(i), c) {
			continue
		}
		//有操作过了
		op := rp.CurrentOperation(int32(i), c)
		if op != nil {
			continue
		}
		//判断所有可能操作
		for _, po := range rp.possibleOperations {
			if po.TargetIndex == i && po.TargetCard.Equal(c) && po.OperationType == ot {
				return true
			}
		}

	}
	return false
}

//判断是否所有操作完成
func (rp *roomPlayer) MaxPriorityOpreations(cs []*card.Card) []*Operation {
	ops := make([]*Operation, 0, 2)
	for i := 0; i < len(cs); i++ {
		c := cs[i]
		//没有操作
		if !rp.HasOperations(int32(i), c) {
			continue
		}
		//有操作过了
		op := rp.CurrentOperation(int32(i), c)
		if op != nil {
			goto AfterOp
		}

		//判断所有可能操作
		for _, po := range rp.possibleOperations {
			if po.TargetIndex == i && po.TargetCard.Equal(c) {
				if op == nil {
					op = po
				} else {
					if op.OperationType.Priority() < po.OperationType.Priority() {
						op = po
					}
				}
			}
		}
		//没有操作
		if op == nil {
			continue
		}
	AfterOp:
		//过的话 忽略
		if op.OperationType == OperationTypePass {
			continue
		}
		//没有操作
		if len(ops) == 0 {
			ops = append(ops, op)
			continue
		}

		//当前优先级低
		if ops[0].OperationType.Priority() > op.OperationType.Priority() {
			continue
		}
		//当前优先级高,清空
		if ops[0].OperationType.Priority() < op.OperationType.Priority() {
			ops = nil
			ops = append(ops, op)
		}

		//胡的化
		if ops[0].OperationType == OperationTypeHu {
			ops = append(ops, op)
		}
		//其它的忽略

	}
	return ops
}

//判断是否有操作
func (rp *roomPlayer) HasOperations(targetIndex int32, c *card.Card) bool {
	for _, op := range rp.possibleOperations {
		if op.TargetCard.Equal(c) && op.TargetIndex == int(targetIndex) {
			return true
		}
	}
	return false
}

//胡牌
func (rp *roomPlayer) Hu(c *card.Card, cardPlayObject CardPlayObject, cardPlayMethod CardPlayMethod, cardPlayOrder CardPlayOrder) int32 {
	rp.result = int32(Hu(card.CardList(rp.cards), rp.composes, c, cardPlayObject, cardPlayMethod, cardPlayOrder))
	fmt.Printf("player result %d\n", rp.result)
	return rp.result
}

//试着胡牌
func (rp *roomPlayer) ifListen(cards card.CardList) bool {
	for ct := card.CardTypeWang; ct <= card.CardTypeTiao; ct++ {
		for cv := card.CardValueOne; cv <= card.CardValueNine; cv++ {
			c := card.NewCard(ct, cv)
			hasFour, _ := findEqual(cards, 4, c)
			if hasFour {
				continue
			}
			result := int32(Hu(cards, rp.composes, c, CardPlayObjectSelf, CardPlayMethodAfterGang, CardPlayOrderNormal))
			if result != 0 {
				return true
			}
		}
	}
	return false
}

//操作
func (rp *roomPlayer) Operate(targetIndex int32, operateType OperationType, targetCardValue int32, cardValues []int32) bool {

	if rp.IsOperate(targetIndex, card.NewCardValue(targetCardValue)) {
		return false
	}

	for _, op := range rp.possibleOperations {
		if op.Equal(int(targetIndex), operateType, targetCardValue, cardValues) {
			rp.currentOpertaions = append(rp.currentOpertaions, op)
			sort.Reverse(OperationList(rp.currentOpertaions))
			return true
		}

	}
	return false
}

//抓鸟
func (rp *roomPlayer) ZhuaNiao() {
	rp.zhongNiao += 1
}

func (rp *roomPlayer) ZhongNiao() int32 {
	return rp.zhongNiao
}

//清除当前操作
func (rp *roomPlayer) ClearOperations() {
	rp.possibleOperations = nil
	rp.currentOpertaions = nil
}

//清除
func (rp *roomPlayer) Clear() {
	rp.playedCards = nil
	rp.cards = nil
	rp.composes = nil
	rp.possibleOperations = nil
	rp.currentOpertaions = nil
	rp.xiaoHus = nil
	rp.currentXiaoHus = nil
	rp.result = 0
	rp.listen = false
	rp.xiaoHuPass = false
	rp.zhongNiao = 0
	rp.currentScore = 0
	rp.zhongNiao = 0
}

//是否可以准备
func (rp *roomPlayer) IfCanPrepare() bool {
	switch rp.actualState {
	case PlayerActualStateInit:
		return true
	case PlayerActualStateSettle:
		return true
	}
	return false
}

//是否可以小胡
func (rp *roomPlayer) IfCanXiaoHu() bool {
	if len(rp.xiaoHus) == 0 {
		return false
	}
	if rp.xiaoHuPass {
		return false
	}
	if len(rp.xiaoHus) == len(rp.currentXiaoHus) {
		return false
	}
	return true
}

//是否完成小胡
func (rp *roomPlayer) IfFinishXiaoHu() bool {
	if len(rp.xiaoHus) == 0 {
		return true
	}
	if rp.xiaoHuPass {
		return true
	}
	if len(rp.xiaoHus) == len(rp.currentXiaoHus) {
		return true
	}
	return false
}

//是否准备
func (rp *roomPlayer) IfPrepare() bool {
	return rp.actualState == PlayerActualStateStart
}

//准备
func (rp *roomPlayer) Prepare() {
	rp.actualState = PlayerActualStateStart
}

//开始
func (rp *roomPlayer) Start() {
	rp.actualState = PlayerActualStatePlay
}

//结算
func (rp *roomPlayer) Settle() {
	rp.actualState = PlayerActualStateSettle
}

//结束
func (rp *roomPlayer) End() {
	rp.actualState = PlayerActualStateEnd
}

//加分
func (rp *roomPlayer) AddScore(s int) {
	rp.score += int64(s)
	rp.currentScore += int64(s)
}

//扣分
func (rp *roomPlayer) SubScore(s int) {
	rp.score -= int64(s)
	rp.currentScore -= int64(s)
}

//加结算数据
func (rp *roomPlayer) AddSettleData(st SettleType) {
	rp.settles[st] = rp.settles[st] + 1
}

func (ro *roomPlayer) ComplementOperations(cs []*card.Card, cardPlayObject CardPlayObject, cardPlayMethod CardPlayMethod, cardPlayOrder CardPlayOrder, next bool, hasCards bool) {

	for j := 0; j < len(cs); j++ {
		c := cs[j]
		//检查是否胡了
		hasOperation := false

		//杠胡后就不能胡了 吃碰后不能胡
		if cardPlayMethod != CardPlayMethoAfterGangHu && cardPlayMethod != CardPlayMethodAfterChiPeng {
			result := Hu(card.CardList(ro.cards), ro.composes, c, cardPlayObject, cardPlayMethod, cardPlayOrder)
			if result != 0 {
				tempHuOperation := &Operation{OperationType: OperationTypeHu}
				tempHuOperation.TargetCard = c
				tempHuOperation.TargetIndex = j
				ro.possibleOperations = append(ro.possibleOperations, tempHuOperation)
				hasOperation = true
			}
		}

		//抢杠或海底或杠后只能胡
		if cardPlayMethod == CardPlayMethodGang || cardPlayOrder == CardPlayOrderLast || cardPlayMethod == CardPlayMethodAfterGang {
			if hasOperation {
				tempPass := &Operation{OperationType: OperationTypePass}
				tempPass.TargetCard = c
				tempPass.TargetIndex = j
				ro.possibleOperations = append(ro.possibleOperations, tempPass)
			}
			continue
		}

		if cardPlayObject == CardPlayObjectOther {
			//检查是否有杠
			if hasCards {
				flag, startIndex := ro.findEqual(3, c)
				if flag {
					op := &Operation{}
					op.OperationType = OperationTypeBu
					op.Cards = []*card.Card{c}
					op.TargetCard = c
					op.TargetIndex = j
					ro.possibleOperations = append(ro.possibleOperations, op)

					//检查是否听牌 有的话添加杠操作
					remainCards := make([]*card.Card, 0, len(ro.cards)-3)
					remainCards = append(remainCards, ro.cards[0:startIndex]...)
					remainCards = append(remainCards, ro.cards[startIndex+3:]...)
					if ro.ifListen(remainCards) {
						opg := &Operation{}
						opg.OperationType = OperationTypeGang
						opg.Cards = []*card.Card{c}
						opg.TargetCard = c
						opg.TargetIndex = j
						ro.possibleOperations = append(ro.possibleOperations, opg)
					}

					// opGang := &Operation{}
					// opGang.OperationType = OperationTypeGang
					// opGang.Cards = []*card.Card{c}
					// opGang.TargetCard = c
					// opGang.TargetIndex = j
					// ro.possibleOperations = append(ro.possibleOperations, opGang)
					hasOperation = true
				}
			}

			if !ro.listen {

				flag, _ := ro.findEqual(2, c)
				if flag {
					hasOperation = true
					tempPeng := &Operation{}
					tempPeng.OperationType = OperationTypePeng
					tempPeng.TargetCard = c
					tempPeng.TargetIndex = j
					ro.possibleOperations = append(ro.possibleOperations, tempPeng)
				}

				//相邻
				//检查吃的
				if next {

					minValue := c.CardValue - 2
					if minValue < card.CardValueOne {
						minValue = card.CardValueOne
					}
					maxValue := c.CardValue + 2
					if maxValue > card.CardValueNine {
						maxValue = card.CardValueNine
					}

					for i := minValue; i <= maxValue-2; i++ {
						tempCardList := make([]*card.Card, 0, 2)
						for j := i; j < i+3; j++ {
							if j == c.CardValue {
								continue
							}
							found, startIndex := ro.findEqualValue(1, c.CardType, j)
							if found {
								tempCardList = append(tempCardList, ro.cards[startIndex])
							}
						}
						if len(tempCardList) == 2 {
							hasOperation = true
							op := &Operation{}
							op.OperationType = OperationTypeEat
							op.Cards = tempCardList
							op.TargetCard = c
							op.TargetIndex = j
							ro.possibleOperations = append(ro.possibleOperations, op)
						}
					}
				}

			}
			if hasOperation {
				tempPass := &Operation{OperationType: OperationTypePass}
				tempPass.TargetCard = c
				tempPass.TargetIndex = j
				ro.possibleOperations = append(ro.possibleOperations, tempPass)
			}
			return
		} else {

			if hasCards {
				//数个数
				if ro.listen {
					//只能杠摸的
					for _, com := range ro.composes {
						if com.ComposeType == ComposeTypePeng {
							if com.CardList[0].Equal(c) {
								opBu := &Operation{}
								opBu.OperationType = OperationTypeBu
								opBu.Cards = []*card.Card{com.CardList[0]}
								opBu.TargetCard = c
								opBu.TargetIndex = j
								ro.possibleOperations = append(ro.possibleOperations, opBu)

								opGang := &Operation{}
								opGang.OperationType = OperationTypeGang
								opGang.Cards = []*card.Card{com.CardList[0]}
								opGang.TargetCard = c
								opGang.TargetIndex = j
								ro.possibleOperations = append(ro.possibleOperations, opGang)
								hasOperation = true
							}
						}
					}
				} else {
					if cardPlayMethod != CardPlayMethoAfterGangHu {
						gangList, _, _, _ := card.Count(ro.cards)
						for _, tc := range gangList {
							op := &Operation{}
							op.OperationType = OperationTypeBu
							op.Cards = []*card.Card{tc}
							op.TargetCard = c
							op.TargetIndex = j
							ro.possibleOperations = append(ro.possibleOperations, op)

							_, startIndex := ro.findEqual(4, tc)
							//检查是否听牌 有的话添加杠操作
							remainCards := make([]*card.Card, 0, len(ro.cards)-4)
							remainCards = append(remainCards, ro.cards[0:startIndex]...)
							remainCards = append(remainCards, ro.cards[startIndex+4:]...)

							if ro.ifListen(remainCards) {
								opg := &Operation{}
								opg.OperationType = OperationTypeGang
								opg.Cards = []*card.Card{tc}
								opg.TargetCard = c
								opg.TargetIndex = j
								ro.possibleOperations = append(ro.possibleOperations, opg)
							}
							// opGang := &Operation{}
							// opGang.OperationType = OperationTypeGang
							// opGang.Cards = []*card.Card{tc}
							// opGang.TargetCard = c
							// opGang.TargetIndex = j
							// ro.possibleOperations = append(ro.possibleOperations, opGang)
							hasOperation = true

						}
					}

					//判断碰的牌是否可以补或杠
					for _, com := range ro.composes {
						if com.ComposeType == ComposeTypePeng {
							for _, tc := range ro.cards {
								if com.CardList[0].Equal(tc) {
									opBu := &Operation{}
									opBu.OperationType = OperationTypeBu
									opBu.Cards = []*card.Card{com.CardList[0]}
									opBu.TargetCard = c
									opBu.TargetIndex = j
									ro.possibleOperations = append(ro.possibleOperations, opBu)

									_, startIndex := ro.findEqual(1, tc)
									//检查是否听牌 有的话添加杠操作
									remainCards := make([]*card.Card, 0, len(ro.cards)-1)
									remainCards = append(remainCards, ro.cards[0:startIndex]...)
									remainCards = append(remainCards, ro.cards[startIndex+1:]...)

									if ro.ifListen(remainCards) {
										opg := &Operation{}
										opg.OperationType = OperationTypeGang
										opg.Cards = []*card.Card{com.CardList[0]}
										opg.TargetCard = c
										opg.TargetIndex = j
										ro.possibleOperations = append(ro.possibleOperations, opg)
									}
									// opGang := &Operation{}
									// opGang.OperationType = OperationTypeGang
									// opGang.Cards = []*card.Card{com.CardList[0]}
									// opGang.TargetCard = c
									// opGang.TargetIndex = j
									// ro.possibleOperations = append(ro.possibleOperations, opGang)
									hasOperation = true
								}
							}
						}
					}
				}
			}

			if hasOperation {
				tempPass := &Operation{OperationType: OperationTypePass}
				tempPass.TargetCard = c
				tempPass.TargetIndex = j
				ro.possibleOperations = append(ro.possibleOperations, tempPass)
			}
		}
		return
	}

	// for j := 0; j < len(cs); j++ {
	// 	c := cs[j]
	// 	//检查是否胡了
	// 	hasOperation := false
	// 	result := Hu(card.CardList(ro.cards), ro.composes, c, cardPlayObject, cardPlayMethod, cardPlayOrder)

	// 	if result != 0 {
	// 		tempHuOperation := &Operation{OperationType: OperationTypeHu}
	// 		tempHuOperation.TargetCard = c
	// 		tempHuOperation.TargetIndex = j
	// 		ro.possibleOperations = append(ro.possibleOperations, tempHuOperation)
	// 		hasOperation = true
	// 	}

	// 	//抢杠或海底 只能胡
	// 	if cardPlayMethod == CardPlayMethodGang || cardPlayOrder == CardPlayOrderLast {
	// 		if hasOperation {
	// 			tempPass := &Operation{OperationType: OperationTypePass}
	// 			tempPass.TargetCard = c
	// 			tempPass.TargetIndex = j
	// 			ro.possibleOperations = append(ro.possibleOperations, tempPass)
	// 		}
	// 		continue
	// 	}

	// 	if cardPlayObject == CardPlayObjectOther || cardPlayMethod == CardPlayMethodAfterGang {
	// 		//检查是否有杠
	// 		if hasCards {
	// 			flag, _ := ro.findEqual(3, c)
	// 			if flag {
	// 				log.Printf("发现补[%s]\n", c)
	// 				op := &Operation{}
	// 				op.OperationType = OperationTypeBu
	// 				op.Cards = []*card.Card{c}
	// 				op.TargetCard = c
	// 				op.TargetIndex = j
	// 				ro.possibleOperations = append(ro.possibleOperations, op)

	// 				opGang := &Operation{}
	// 				opGang.OperationType = OperationTypeGang
	// 				opGang.Cards = []*card.Card{c}
	// 				opGang.TargetCard = c
	// 				opGang.TargetIndex = j
	// 				ro.possibleOperations = append(ro.possibleOperations, opGang)
	// 				hasOperation = true
	// 			}
	// 		}
	// 		//检查是否听牌 有的话添加杠操作
	// 		// remainCards := make([]*card.Card, 0, len(ro.cards)-3)
	// 		// remainCards = append(remainCards, ro.cards[0:startIndex]...)
	// 		// remainCards = append(remainCards, ro.cards[startIndex+3:]...)
	// 		// if ro.ifListen(remainCards) {
	// 		// 	log.Printf("发现杠[%s]\n", c)
	// 		// 	opg := &Operation{}
	// 		// 	opg.OperationType = OperationTypeGang
	// 		// 	opg.Cards = []*card.Card{c}
	// 		// 	opg.TargetCard = c
	// 		// 	ro.possibleOperations = append(ro.possibleOperations, opg)
	// 		// }
	// 		//}
	// 		//	}

	// 		if !ro.listen {
	// 			//不胡的话
	// 			//	if result == 0 {
	// 			//检查是否可以碰
	// 			flag, _ := ro.findEqual(2, c)
	// 			if flag {
	// 				hasOperation = true
	// 				tempPeng := &Operation{}
	// 				tempPeng.OperationType = OperationTypePeng
	// 				tempPeng.TargetCard = c
	// 				tempPeng.TargetIndex = j
	// 				ro.possibleOperations = append(ro.possibleOperations, tempPeng)
	// 			}

	// 			//相邻
	// 			//检查吃的
	// 			if next {

	// 				minValue := c.CardValue - 2
	// 				if minValue < card.CardValueOne {
	// 					minValue = card.CardValueOne
	// 				}
	// 				maxValue := c.CardValue + 2
	// 				if maxValue > card.CardValueNine {
	// 					maxValue = card.CardValueNine
	// 				}

	// 				for i := minValue; i <= maxValue-2; i++ {
	// 					tempCardList := make([]*card.Card, 0, 2)
	// 					for j := i; j < i+3; j++ {
	// 						if j == c.CardValue {
	// 							continue
	// 						}
	// 						found, startIndex := ro.findEqualValue(1, c.CardType, j)
	// 						if found {
	// 							tempCardList = append(tempCardList, ro.cards[startIndex])
	// 						}
	// 					}
	// 					if len(tempCardList) == 2 {
	// 						hasOperation = true
	// 						op := &Operation{}
	// 						op.OperationType = OperationTypeEat
	// 						op.Cards = tempCardList
	// 						op.TargetCard = c
	// 						op.TargetIndex = j
	// 						ro.possibleOperations = append(ro.possibleOperations, op)
	// 					}
	// 				}
	// 			}

	// 			//}

	// 		}
	// 		if hasOperation {
	// 			tempPass := &Operation{OperationType: OperationTypePass}
	// 			tempPass.TargetCard = c
	// 			tempPass.TargetIndex = j
	// 			ro.possibleOperations = append(ro.possibleOperations, tempPass)
	// 		}
	// 		return
	// 	} else {

	// 		if hasCards {
	// 			//数个数
	// 			gangList, _, _, _ := card.Count(ro.cards)
	// 			for _, tc := range gangList {
	// 				//flag, startIndex := ro.findEqual(4, tc)
	// 				// if flag {

	// 				op := &Operation{}
	// 				op.OperationType = OperationTypeBu
	// 				op.Cards = []*card.Card{tc}
	// 				op.TargetCard = c
	// 				op.TargetIndex = j
	// 				ro.possibleOperations = append(ro.possibleOperations, op)

	// 				opGang := &Operation{}
	// 				opGang.OperationType = OperationTypeGang
	// 				opGang.Cards = []*card.Card{tc}
	// 				opGang.TargetCard = c
	// 				opGang.TargetIndex = j
	// 				ro.possibleOperations = append(ro.possibleOperations, opGang)
	// 				hasOperation = true

	// 				//检查是否听牌 有的话添加杠操作
	// 				// remainCards := make([]*card.Card, 0, len(ro.cards)-4)
	// 				// remainCards = append(remainCards, ro.cards[0:startIndex]...)
	// 				// remainCards = append(remainCards, ro.cards[startIndex+4:]...)
	// 				// log.Printf("remain cards [%s]\n", remainCards)
	// 				// if ro.ifListen(remainCards) {
	// 				// 	log.Printf("发现杠[%s]\n", tc)
	// 				// 	opg := &Operation{}
	// 				// 	opg.OperationType = OperationTypeGang
	// 				// 	opg.Cards = []*card.Card{tc}
	// 				// 	opg.TargetCard = c
	// 				// 	ro.possibleOperations = append(ro.possibleOperations, opg)
	// 				// }
	// 			}

	// 			//判断碰的牌是否可以补或杠
	// 			for _, com := range ro.composes {
	// 				if com.ComposeType == ComposeTypePeng {
	// 					for _, tc := range ro.cards {
	// 						if com.CardList[0].Equal(tc) {
	// 							opBu := &Operation{}
	// 							opBu.OperationType = OperationTypeBu
	// 							opBu.Cards = []*card.Card{com.CardList[0]}
	// 							opBu.TargetCard = c
	// 							opBu.TargetIndex = j
	// 							ro.possibleOperations = append(ro.possibleOperations, opBu)

	// 							opGang := &Operation{}
	// 							opGang.OperationType = OperationTypeGang
	// 							opGang.Cards = []*card.Card{com.CardList[0]}
	// 							opGang.TargetCard = c
	// 							opGang.TargetIndex = j
	// 							ro.possibleOperations = append(ro.possibleOperations, opGang)
	// 							hasOperation = true
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}

	// 		if hasOperation {
	// 			tempPass := &Operation{OperationType: OperationTypePass}
	// 			tempPass.TargetCard = c
	// 			tempPass.TargetIndex = j
	// 			ro.possibleOperations = append(ro.possibleOperations, tempPass)
	// 		}
	// 	}
	// 	return
	// }

}

func (rp *roomPlayer) findEqualValue(num int, ct card.CardType, cv card.CardValue) (find bool, startIndex int) {
	startIndex = -1
	for i, tempCard := range rp.cards {
		if tempCard.CardValue == cv && tempCard.CardType == ct {
			if startIndex == -1 {
				startIndex = i
			}
			num--
			if num == 0 {
				find = true
				return
			}
		}
	}
	return
}

func (rp *roomPlayer) findEqual(num int, c *card.Card) (find bool, startIndex int) {
	return findEqual(rp.cards, num, c)
}

func findEqual(cs []*card.Card, num int, c *card.Card) (find bool, startIndex int) {
	startIndex = -1
	cardValue := card.Value(c)
	for i, tempCard := range cs {
		if card.Value(tempCard) == cardValue {
			if startIndex == -1 {
				startIndex = i
			}
			num--
			if num == 0 {
				find = true
				return
			}
		}
	}
	return
}

func (rp *roomPlayer) String() string {
	return fmt.Sprintf("玩家[%d],位置[%d],分数[%d],状态[%s],实际状态[%s],小胡列表[%s],手牌[%s],下的牌[%s],吃碰杠[%s],可能的组合[%s],当前的操作[%s]", rp.id, rp.position, rp.score, rp.state, rp.actualState, rp.xiaoHus, rp.cards, rp.playedCards, rp.composes, rp.possibleOperations, rp.currentOpertaions)
}

type RoomPlayerManager interface {
	Players() []Player
	AddPlayer(player Player) bool
	RemovePlayer(player Player) bool
	DisconnectPlayer(player Player) bool 
	LeaveTimePlayer(player Player,state bool) bool
	ReconnectPlayer(p player.Player) bool
	GetPlayerByPos(pos int) Player
	GetPlayerById(id int64) Player
	NumPlayers() int
	CurrentNumPlayers() int
	Next(p Player) Player
}

type roomPlayerManager struct {
	positions  []int
	numPlayers int
	players    []Player
}

//初始化
func (rpm *roomPlayerManager) init() {
	rpm.players = make([]Player, 0, rpm.numPlayers)
	rpm.positions = make([]int, 0, rpm.numPlayers)

	for i := 0; i < rpm.numPlayers; i++ {
		rpm.positions = append(rpm.positions, i)
	}
}

func (rpm *roomPlayerManager) Players() []Player {
	return rpm.players
}

func (rpm *roomPlayerManager) CurrentNumPlayers() int {
	return len(rpm.players)
}

func (rpm *roomPlayerManager) GetPlayerByPos(pos int) Player {
	for _, pl := range rpm.players {
		if pl.Position() == pos {
			return pl
		}
	}
	return nil
}

func (rpm *roomPlayerManager) GetPlayerById(playerId int64) Player {
	for _, pl := range rpm.players {
		if pl.Id() == playerId {
			return pl
		}
	}
	return nil
}

func (rpm *roomPlayerManager) NumPlayers() int {
	return rpm.numPlayers
}

func (rpm *roomPlayerManager) Next(p Player) Player {
	pos := p.Position()
	pos++
	if pos >= len(rpm.players) {
		pos = 0
	}
	return rpm.GetPlayerByPos(pos)
}

//判断是否满人
func (rpm *roomPlayerManager) isFull() bool {
	if len(rpm.players) >= rpm.numPlayers {
		return true
	}
	return false
}

func (rpm *roomPlayerManager) AddPlayer(player Player) bool {
	if rpm.isFull() {
		return false
	}
	//检查是否已经加入了
	p := rpm.GetPlayerById(player.Id())
	if p != nil {
		return false
	}

	pos := rpm.positions[0]
	rpm.positions = rpm.positions[1:]
	player.SetPosition(pos)
	rpm.players = append(rpm.players, player)
	return true
}

func (rpm *roomPlayerManager) ReconnectPlayer(pl player.Player) bool {
	p := rpm.GetPlayerById(pl.Id())
	if p.Player() != nil {
		return false
	}
	p.Reconnect(pl)
	return true
}

func (rpm *roomPlayerManager) DisconnectPlayer(player Player) bool {
	p := rpm.GetPlayerByPos(player.Position())
	if p != player {
		return false
	}
	player.Disconnect()
	return true
}

func (rpm *roomPlayerManager) LeaveTimePlayer(player Player,state bool) bool {
	p := rpm.GetPlayerByPos(player.Position())
	if p != player {
		return false
	}   
	player.LevaeTime(state)
	return true
}

func (rpm *roomPlayerManager) RemovePlayer(player Player) bool {
	index := -1
	for i, p := range rpm.players {
		if p == player {
			index = i
			break
		}
	}
	if index == -1 {
		return false
	}
	rpm.positions = append(rpm.positions, player.Position())
	rpm.players = append(rpm.players[:index], rpm.players[index+1:]...)
	return true
}

func NewRoomPlayerManager(nplayers int) RoomPlayerManager {
	pm := &roomPlayerManager{
		numPlayers: nplayers,
	}
	pm.init()
	return pm
}
