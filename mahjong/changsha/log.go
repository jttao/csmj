package changsha

import (
	"encoding/json"

	"game/mahjong/card"
)

type ChangShaLogType int

const (
	ChangShaLogTypeInit ChangShaLogType = iota
	ChangShaLogTypeXiaoHu
	ChangShaLogTypePlayerPlay
	ChangShaLogTypePlayerMo
	ChangShaLogTypePlayerWaitOperation
	ChangShaLogTypePlayerOperation
	ChangShaLogTypePlayerChi
	ChangShaLogTypePlayerPeng
	ChangShaLogTypePlayerBu
	ChangShaLogTypePlayerGang
	ChangShaLogTypePlayerGangMo
	ChangShaLogTypePlayerHu
)

var (
	changShaLogTypeMap = map[ChangShaLogType]string{}
)

func (cslt ChangShaLogType) String() string {
	return changShaLogTypeMap[cslt]
}

//日志
type ChangShaLog struct {
	ChangShaLogType int32  `json:"changShaLogType"`
	Time            int64  `json:"time"`
	Content         string `json:"content"`
}

//长沙房间初始化
type ChangShaInitLog struct {
	//玩家信息
	Players []*ChangShaLogPlayer `json:"players"`
	//庄家位置
	BankerPos int32 `json:"bankerPos"`
}

func buildChangShaInitLog(r *Room, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypeInit)
	initLog := &ChangShaInitLog{}
	initLog.Players = make([]*ChangShaLogPlayer, 0, len(r.RoomPlayerManager().Players()))
	for _, pl := range r.RoomPlayerManager().Players() {
		logPlayer := buildChangShaLogPlayer(pl)
		initLog.Players = append(initLog.Players, logPlayer)
	}
	initLog.BankerPos = int32(r.bankerPos)
	content, err := json.Marshal(initLog)
	if err != nil {
		panic(err)
	}
	csl.Time = now
	csl.Content = string(content)
	return csl
}

func buildChangShaLogPlayer(pl Player) *ChangShaLogPlayer {
	logPlayer := &ChangShaLogPlayer{}
	logPlayer.PlayerId = pl.Id()
	logPlayer.Name = pl.Name()
	logPlayer.Img = pl.Image()
	logPlayer.Score = pl.Score()
	logPlayer.Cards = card.Values(pl.Cards())
	logPlayer.Position = int32(pl.Position())
	return logPlayer
}

//玩家信息
type ChangShaLogPlayer struct {
	//玩家id
	PlayerId int64 `json:"playerId"`
	//玩家名字
	Name string `json:"name"`
	//玩家图片
	Img string `json:"img"`
	//玩家分数
	Score int64 `json:"score"`
	//卡牌
	Cards []int32 `json:"cards"`
	//玩家位置
	Position int32 `json:"position"`
}

//长沙小胡
type ChangShaXiaoHuLog struct {
	PlayerId   int64   `json:"playerId"`
	XiaoHuType int32   `json:"xiaoHuType"`
	Cards      []int32 `json:"cards"`
}

func buildChangShaXiaoHuLog(pl Player, xho *XiaoHuOperation, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypeXiaoHu)
	xiaoHuLog := &ChangShaXiaoHuLog{}
	xiaoHuLog.PlayerId = pl.Id()
	xiaoHuLog.XiaoHuType = int32(xho.XiaoHuType)
	xiaoHuLog.Cards = card.Values(xho.Cards)
	content, err := json.Marshal(xiaoHuLog)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	csl.Time = now
	return csl
}

type ChangShaPlayerPlayLog struct {
	PlayerId int64 `json:"playerId"`
	Card     int32 `json:"card"`
}

//打牌
func buildChangShaPlayerPlayLog(pl Player, c *card.Card, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerPlay)
	csl.Time = now
	playerPlayLog := &ChangShaPlayerPlayLog{}
	playerPlayLog.Card = card.Value(c)
	playerPlayLog.PlayerId = pl.Id()
	content, err := json.Marshal(playerPlayLog)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

//摸牌
type ChangShaPlayerMoLog struct {
	PlayerId int64 `json:"playerId"`
	Card     int32 `json:"card"`
}

//摸牌
func buildChangShaPlayerMoLog(pl Player, c *card.Card, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerMo)
	csl.Time = now
	playerPlayMo := &ChangShaPlayerMoLog{}
	playerPlayMo.Card = card.Value(c)
	playerPlayMo.PlayerId = pl.Id()
	content, err := json.Marshal(playerPlayMo)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

type ChangShaPlayerChiLog struct {
	PlayerId int64   `json:"playerId"`
	Cards    []int32 `json:"cards"`
}

//吃
func buildChangShaPlayerChiLog(pl Player, cs []int32, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerChi)
	csl.Time = now
	playerChi := &ChangShaPlayerChiLog{}
	playerChi.Cards = cs
	playerChi.PlayerId = pl.Id()
	content, err := json.Marshal(playerChi)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

type ChangShaPlayerPengLog struct {
	PlayerId int64 `json:"playerId"`
	Card     int32 `json:"card"`
}

//碰
func buildChangShaPlayerPengLog(pl Player, cs int32, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerPeng)
	csl.Time = now
	playerPeng := &ChangShaPlayerPengLog{}
	playerPeng.Card = cs
	playerPeng.PlayerId = pl.Id()
	content, err := json.Marshal(playerPeng)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

type ChangShaPlayerBuLog struct {
	PlayerId int64 `json:"playerId"`
	Card     int32 `json:"card"`
	GangType int32 `json:"gangType"`
}

//碰
func buildChangShaPlayerBuLog(pl Player, cs int32, gt int32, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerBu)
	csl.Time = now
	playerBu := &ChangShaPlayerBuLog{}
	playerBu.Card = cs
	playerBu.PlayerId = pl.Id()
	playerBu.GangType = gt
	content, err := json.Marshal(playerBu)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

type ChangShaPlayerGangLog struct {
	PlayerId int64 `json:"playerId"`
	Card     int32 `json:"card"`
	GangType int32 `json:"gangType"`
}

//杠
func buildChangShaPlayerGangLog(pl Player, cs int32, gt int32, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerGang)
	csl.Time = now
	playerGang := &ChangShaPlayerGangLog{}
	playerGang.Card = cs
	playerGang.PlayerId = pl.Id()
	playerGang.GangType = gt
	content, err := json.Marshal(playerGang)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

//杠摸
type ChangShaPlayerGangMoLog struct {
	PlayerId int64   `json:"playerId"`
	Cards    []int32 `json:"cards"`
}

//摸牌
func buildChangShaPlayerGangMoLog(pl Player, cs []*card.Card, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerGangMo)
	csl.Time = now
	playerPlayGangMo := &ChangShaPlayerGangMoLog{}
	playerPlayGangMo.Cards = card.Values(cs)
	playerPlayGangMo.PlayerId = pl.Id()
	content, err := json.Marshal(playerPlayGangMo)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

type ChangShaPlayerHuLog struct {
	PlayerId int64 `json:"playerId"`
	Card     int32 `json:"card"`
}

//杠
func buildChangShaPlayerHuLog(pl Player, cs int32, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerHu)
	csl.Time = now
	playerHu := &ChangShaPlayerHuLog{}
	playerHu.Card = cs
	playerHu.PlayerId = pl.Id()

	content, err := json.Marshal(playerHu)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

type ChangShaPlayerOperation struct {
	OperationType int32   `json:operationType`
	Cards         []int32 `json:"cards"`
	TargetCard    int32   `json:"targetCard"`
}

type ChangShaPlayerWaitOperationLog struct {
	PlayerId   int64                      `json:"playerId"`
	Operations []*ChangShaPlayerOperation `json:"operations"`
}

func buildOperations(pl Player) []*ChangShaPlayerOperation {
	ops := make([]*ChangShaPlayerOperation, 0, len(pl.PossibleOperations()))
	for _, op := range pl.PossibleOperations() {
		top := buildOperation(op)
		ops = append(ops, top)
	}
	return ops
}

func buildOperation(op *Operation) *ChangShaPlayerOperation {
	top := &ChangShaPlayerOperation{}
	opType := int32(op.OperationType)
	top.OperationType = opType
	top.Cards = card.Values(op.Cards)
	top.TargetCard = card.Value(op.TargetCard)
	return top
}

//等候玩家操作
func buildChangShaPlayerWaitOpereationLog(r *Room, pl Player, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerWaitOperation)
	csl.Time = now
	playerWaitOperationLog := &ChangShaPlayerWaitOperationLog{}

	playerWaitOperationLog.PlayerId = pl.Id()
	playerWaitOperationLog.Operations = buildOperations(pl)

	content, err := json.Marshal(playerWaitOperationLog)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}

type ChangShaPlayerOperationLog struct {
	PlayerId      int64   `json:"playerId"`
	OperationType int32   `json:"operationType"`
	TargetCard    int32   `json:"targetCard"`
	Cards         []int32 `json:"cards"`
	TargetIndex   int32   `json:"targetIndex"`
}

//等候玩家操作
func buildChangShaPlayerOpereationLog(r *Room, pl Player, targetIndex int32, operationType int32, targetCard int32, cardValues []int32, now int64) *ChangShaLog {
	csl := &ChangShaLog{}
	csl.ChangShaLogType = int32(ChangShaLogTypePlayerOperation)
	csl.Time = now
	playerOperationLog := &ChangShaPlayerOperationLog{}

	playerOperationLog.PlayerId = pl.Id()
	playerOperationLog.OperationType = operationType
	playerOperationLog.TargetCard = targetCard
	playerOperationLog.Cards = cardValues
	playerOperationLog.TargetIndex = targetIndex
	content, err := json.Marshal(playerOperationLog)
	if err != nil {
		panic(err)
	}
	csl.Content = string(content)
	return csl
}
