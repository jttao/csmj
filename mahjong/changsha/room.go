package changsha

import (
	"context"
	"fmt"
	"math"
	"time"

	"game/mahjong/card"

	"game/mahjong/server/player"

	log "github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

const (
	initialCards = 13
)

const (
	//小胡分数
	xiaoHuScore = 2
	//抓鸟分数
	zhuaNiaoScore = 1
)

const (
	xiaoHuZiMoScore   = 2
	xiaoHuJiePaoScore = 1
	daHuZiMoScore     = 6
	daHuJiePaoScore   = 6
	//庄家额外算分
	zhuanExtraScore = 1
)

var (
	roomStateMap = map[RoomState]string{
		RoomStateInit: "init",
		//等侯开始
		RoomStateWait: "wait",
		//开始发牌中
		RoomStatePrepare: "prepare",
		//等待小胡决定
		RoomStateWaitPlayerXiaoHuAction: "xiaohu action",
		//等待玩家下牌
		RoomStateWaitPlayerPlay: "player play",
		//等待玩家行动
		RoomStateWaitPlayerAction: "player action",
		//结算
		RoomStateSettle: "settle",
		//结束
		RoomStateEnd: "end",
	}
)

type RoomState int

func (rs RoomState) String() string {
	return roomStateMap[rs]
}

const (
	//初始化
	RoomStateInit RoomState = iota
	//等侯开始
	RoomStateWait
	//开始发牌中
	RoomStatePrepare
	//等待小胡决定
	RoomStateWaitPlayerXiaoHuAction
	//等待玩家下牌
	RoomStateWaitPlayerPlay
	//等待玩家行动
	RoomStateWaitPlayerAction
	//结算
	RoomStateSettle
	//结束
	RoomStateEnd
)

type RoomDelegate interface {
	OnRoomStart(room *Room)
	OnRoomPlayerJoin(room *Room, player Player)
	OnRoomPlayerLeave(room *Room, player Player)
	OnRoomSelectBanker(room *Room, pos int)
	OnRoomDealCards(room *Room)
	OnRoomXiaoHu(room *Room)
	OnRoomPlayerXiaoHu(room *Room, player Player, xho *XiaoHuOperation)
	OnRoomPlayerXiaoHuPass(room *Room, player Player)
	OnRoomPlayerMo(room *Room, player Player, c *card.Card)
	OnRoomPlayerGangMo(room *Room, player Player, cs []*card.Card)
	OnRoomWaitPlayerPlay(room *Room, pl Player)
	OnRoomWaitPlayerAction(room *Room)
	OnRoomPlayerPlayCard(room *Room, player Player, c *card.Card)
	OnRoomPlayerChi(room *Room, player Player, c *card.Card, cardValues []int32)
	OnRoomPlayerPeng(room *Room, player Player, c *card.Card)
	OnRoomPlayerBu(room *Room, player Player, c *card.Card, gt GangType)
	OnRoomPlayerGang(room *Room, player Player, c *card.Card, gt GangType)
	OnRoomPlayersSettle(room *Room, players []Player, dianPaoPlayerId int64, c *card.Card)
	OnRoomClear(room *Room)
	OnRoomReconnect(room *Room, player Player)
	OnRoomDisconnect(room *Room, player Player)
	OnLevaeTime(room *Room, player Player)
	OnRoomPlayerStart(room *Room, player Player)
	OnRoomPlayerDissolve(room *Room, player Player)
	OnRoomPlayerDissolveAgree(room *Room, player Player, flag bool)
	OnRoomLiuJu(room *Room)
	OnRoomEnd(room *Room,start bool)
	OnRoomHaiDiAsk(room *Room, pl Player)
	OnRoomHaiDiAnswer(room *Room, pl Player, c *card.Card, flag bool)
	OnRoomSettle(room *Room)
}

type RoomConfig struct {
	WaitTime         int64 `json:"waitTime"`
	PrepareTime      int64 `json:"prepareTime"`
	DealCardsTime    int64 `json:"dealCardsTime"`
	XiaoHuWaitTime   int64 `json:"xiaoHuWaitTime"`
	PlayerPlayTime   int64 `json:"playerPlayTime"`
	PlayerActionTime int64 `json:"playerActionTime"`
	SettleTime       int64 `json:"settleTime"`
	ClearTime        int64 `json:"clearTime"`
	DissolveTime     int64 `json:"DissolveTime"`
}

//自定义配置
type CustomRoomConfig struct {
	ZhuangXian bool `json:"zhuanXian"`
	ZhuaNiao   int  `json:"zhuaNiao"`
	//抓鸟算法 默认乘法 true:2的密次方
	ZhuaNiaoAlg bool `json:"zhuaNiaoAlg"` 
} 

//创建新房间
func NewRoom(config *RoomConfig, customConfig *CustomRoomConfig, numPlayer int, round int, roomId int64,forbidIp int,openRoomType int, createTime int64, forbidJoinTime int64,lastGameTime int64,ownerId int64,rd RoomDelegate) *Room {
	r := &Room{
		config:       config,
		customConfig: customConfig,
	}
	r.playerManager = NewRoomPlayerManager(numPlayer)
	r.deck = NewDeck()
	r.roomId = roomId
	r.delegate = rd
	r.openRoomType = openRoomType
	r.createTime = createTime
	r.forbidJoinTime = forbidJoinTime
	r.lastGameTime = lastGameTime
	r.ownerId = ownerId
	r.init()
	r.totalRound = int32(round)
	r.forbidIp = int32(forbidIp) 
	r.ctx = context.Background()
	return r
}

type Room struct {
	ctx context.Context
	//创建时间
	createTime int64
	//单局开始时间
	roundStartTime int64
	//录像
	logList []*ChangShaLog
	//日志
	logger *log.Logger
	//	RoomProcessor *mahjong.RoomProcessor
	//固定配置
	config *RoomConfig
	//自定义配置
	customConfig *CustomRoomConfig

	//房间id
	roomId int64
	//房间ID
	ownerId int64 
	//玩家管理器
	playerManager RoomPlayerManager
	//牌组
	deck Deck
	//房间委托
	delegate RoomDelegate
	//等候时间
	xiaoHuWait int64
	//总共盘数
	totalRound int32
	//当前盘数
	currentRound int32
	//盘数id
	roundId int
	//庄家位置
	bankerPos int
	//下把庄家位置
	nextBankerPos int
	//当前玩家
	currentPlayer Player
	//当前牌
	currentCards []*card.Card
	//当前牌的来源
	currentCardPlayMethod CardPlayMethod
	//当前牌的顺序
	currentCardPlayOrder CardPlayOrder
	//当前牌的状态
	currentCardPlayObject CardPlayObject
	//是否是杠
	isGang bool

	//最后标记时间
	lastTime int64
	//房间状态
	state RoomState
	//当前最好的行为的玩家
	bestOpereatePlayers []Player
	//申请解散开始时间
	dissolveStartTime int64
	//同意的玩家
	agreeDissolvePlayers map[int64]Player
	//申请的玩家
	applyDissolvePlayer Player

	//鸟牌
	niaoPais []*card.Card
	//鸟牌id
	niaoPaiPlayerIds []int64

	forbidIp int32  
	//房间类型
	openRoomType int
	//禁止加入时间
	forbidJoinTime int64
	//最后游戏时间
	lastGameTime int64

	//海底牌玩家
	haidiPlayer Player
}

func (r *Room) Context() context.Context {
	return r.ctx
}

func (r *Room) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *Room) CreateTime() int64 {
	return r.createTime
}

func (r *Room) LogList() []*ChangShaLog {
	return r.logList
}

func (r *Room) State() RoomState {
	return r.state
}

func (r *Room) RoomId() int64 {
	return r.roomId
}

func (r *Room) OwnerId() int64 {
	return r.ownerId
}

func (r *Room) Deck() Deck {
	return r.deck
}

func (r *Room) RoundStartTime() int64 {
	return r.roundStartTime
}

func (r *Room) CurrentRound() int32 {
	return r.currentRound
}

func (r *Room) TotalRound() int32 {
	return r.totalRound
}

func (r *Room) Name() string {
	return "长沙麻将"
}

func (r *Room) ForbidIp() int32 {
	return r.forbidIp
}

func (r *Room) ForbidJoinTime() int64 {
	return r.forbidJoinTime
}

func (r *Room) LastGameTime() int64 {
	return r.lastGameTime
}

func (r *Room) GetOpenRoomType() int { 
	return r.openRoomType
}

func (r *Room) CustomRoomConfig() *CustomRoomConfig {
	return r.customConfig
}

func (r *Room) RoomPlayerManager() RoomPlayerManager {
	return r.playerManager
}

func (r *Room) CurrentPlayer() Player {
	return r.currentPlayer
}

func (r *Room) CurrentCards() []*card.Card {
	return r.currentCards
}

func (r *Room) CurrentCardPlayMethod() CardPlayMethod {
	return r.currentCardPlayMethod
}

func (r *Room) CurrentCardPlayOrder() CardPlayOrder {
	return r.currentCardPlayOrder
}

func (r *Room) CurrentCardPlayObject() CardPlayObject {
	return r.currentCardPlayObject
}

func (r *Room) DissolveStartTime() int64 {
	return r.dissolveStartTime
}

func (r *Room) AgreeDissolvePlayers() map[int64]Player {
	return r.agreeDissolvePlayers
}

func (r *Room) ApplyDissolvePlayer() Player {
	return r.applyDissolvePlayer
}

func (r *Room) BankerPos() int {
	return r.bankerPos
}

func (r *Room) NiaoPais() []*card.Card {
	return r.niaoPais
}

func (r *Room) NiaoPaiPlayerIds() []int64 {
	return r.niaoPaiPlayerIds
}

//tick
func (r *Room) Tick() {

	now := r.now()

	//判断是否在申请解散中
	if r.ifDissolve() {
		elapse := now - r.dissolveStartTime
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"过了":   elapse / 1000,
		}).Debug("解散倒数中")
		if elapse >= r.config.DissolveTime {
			//流局
			r.liuJu()
		}
		return
	}
	
	//是否超过游戏时间
	if now >= r.lastGameTime {
		r.liuJu()
		return
	}

	elapse := now - r.lastTime
	switch r.state {
	case RoomStateInit:
		{
			if now >= r.forbidJoinTime {
				r.autoLeaveRoom()
			}
		}
	case RoomStateWait:
		{
			if elapse >= r.config.WaitTime {
				r.enterPrepareState(now)
			}
		}
	case RoomStatePrepare:
		{
			if elapse >= r.config.PrepareTime {
				r.firstCheck()
			}
		}
	case RoomStateWaitPlayerXiaoHuAction:
	case RoomStateWaitPlayerPlay:
	case RoomStateWaitPlayerAction:
	case RoomStateSettle:
	case RoomStateEnd:
	}
}

func (r *Room) enterInitState(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("进入初始化状态")
	r.state = RoomStateInit
	r.lastTime = now
}

func (r *Room) enterWaitState(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("进入等候状态")
	r.state = RoomStateWait
	r.lastTime = now
}

func (r *Room) enterPrepareState(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("进入准备状态")
	r.state = RoomStatePrepare
	r.lastTime = now
	r.roundStartTime = now
	r.prepareStart()
}

func (r *Room) enterWaitPlayerXiaoHuAction(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("进入小胡等待")
	r.state = RoomStateWaitPlayerXiaoHuAction
	r.lastTime = now
}

func (r *Room) enterWaitPlayerPlay(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("进入等候玩家打牌状态")
	r.state = RoomStateWaitPlayerPlay
	r.lastTime = now
}

func (r *Room) enterWaitPlayerAction(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("等候玩家行为")
	r.state = RoomStateWaitPlayerAction
	r.lastTime = now
}

func (r *Room) enterSettle(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("进入结算状态")
	r.state = RoomStateSettle
	r.lastTime = now
}

func (r *Room) enterEnd(now int64) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("进入总结算状态")
	r.state = RoomStateEnd
	r.lastTime = now
}

//是否所有都同意了
func (r *Room) ifXiaohuFinish() bool {
	for _, pl := range r.RoomPlayerManager().Players() {
		if !pl.IfFinishXiaoHu() {
			return false
		}
	}
	return true
}

//是否所有都同意了
func (r *Room) ifAllDissolveAgree() bool {
	num := r.RoomPlayerManager().NumPlayers()/2 + 1
	if len(r.agreeDissolvePlayers) >= num {
		return true
	}
	return false
}

//是否正在解散
func (r *Room) ifDissolve() bool {
	return r.dissolveStartTime > 0
}

//判断是否可以进行游戏
func (r *Room) ifCanStart() bool {
	if r.playerManager.CurrentNumPlayers() == r.playerManager.NumPlayers() {
		for _, p := range r.playerManager.Players() {
			if !p.IfPrepare() {
				return false
			}
		}
		return true
	}
	return false
}

//是否可以离开
func (r *Room) ifCanLeave() bool {
	if r.state == RoomStateInit {
		return true
	}
	return false
}

//是否开始了
func (r *Room) ifStart() bool {
	if r.state == RoomStateInit {
		return false
	}
	if r.state == RoomStateEnd {
		return false
	}
	return true
}

//是否结束了
func (r *Room) ifEnd() bool {
	return r.state == RoomStateEnd
}

//离开房间
func (r *Room) LeavePlayer(p Player) bool {
	//只有还没有开始的时候才可以退出
	if r.state != RoomStateInit {
		return false
	}

	tp := r.playerManager.GetPlayerById(p.Id())
	//已经不存在了
	if tp == nil {
		return false
	}

	return r.leaveRoomPlayer(tp)
}

//自动解散
func (r *Room) autoLeaveRoom() {
	if r.openRoomType == 0 {
		ownerPlayer := r.playerManager.GetPlayerById(r.ownerId)
		r.leaveRoomPlayer(ownerPlayer)
	} else {
		r.end(false)
	}
}

//离开房间
func (r *Room) leaveRoomPlayer(tp Player) bool {
	if tp == nil {
		r.end(false)
		return true
	}
	flag := r.PlayerLeave(tp)
	if flag {
		if r.openRoomType == 0 {
			if tp.Id() == r.ownerId {
				//解散房间
				r.end(false)
			}
		}   
	}
	return flag
}

//玩家加入
func (r *Room) PlayerJoin(p Player) bool {
	if r.state != RoomStateInit {
		return false
	} 
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Debug("玩家请求加入房间")

	flag := r.playerManager.AddPlayer(p)
	if !flag {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   p.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Warn("玩家加入房间失败")
		return flag
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家加入房间")

	//设置房间主人
	//if r.playerManager.CurrentNumPlayers() == 1 {
	//	r.ownerPlayer = p
	//}

	r.delegate.OnRoomPlayerJoin(r, p)
	r.PrepareStart(p)
	return flag
}

//玩家离开
func (r *Room) PlayerLeave(p Player) bool {
	//只有还没有开始的时候才可以退出
	if r.state != RoomStateInit {
		return false
	}   
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Debug("玩家请求离开房间")

	tp := r.playerManager.GetPlayerById(p.Id())
	if tp == nil {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   p.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Warn("玩家请求离开房间, 玩家不存在")
		return false
	}

	if !r.ifCanLeave() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   p.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Warn("玩家不能离开房间")
		return false
	}
	
	r.delegate.OnRoomPlayerLeave(r, p)

	flag := r.playerManager.RemovePlayer(p)
	if !flag {
		return flag
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家离开房间")
	
	now := r.now()
	r.enterInitState(now)
	
	return true
}

//申请解散
func (r *Room) PlayerDissolve(p Player) {
	if r.state == RoomStateEnd {
		return
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Debug("玩家请求解散")

	if !r.ifStart() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   p.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Warn("玩家不能请求解散房间")
		return
	}
	if r.ifDissolve() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   p.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Warn("玩家不能请求解散房间，已经有人申请解散了")
		return
	}
	now := r.now()
	r.dissolveStartTime = now
	r.agreeDissolvePlayers = make(map[int64]Player)
	r.agreeDissolvePlayers[p.Id()] = p
	r.applyDissolvePlayer = p
	r.delegate.OnRoomPlayerDissolve(r, p)
	return
}

//同意解散
func (r *Room) PlayerDissolveAgree(p Player, flag bool) {
	if r.state == RoomStateEnd {
		return
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		"同意玩家": fmt.Sprintf("%s", r.agreeDissolvePlayers),
		"同意":   flag,
	}).Debug("玩家请求同意解散")

	if !r.ifDissolve() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   p.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
			"同意玩家": fmt.Sprintf("%s", r.agreeDissolvePlayers),
			"同意":   flag,
		}).Warn("玩家不能请求同意解散房间，房间不是在申请解散中")
		return
	}

	_, exist := r.agreeDissolvePlayers[p.Id()]
	if exist {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   p.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
			"同意玩家": fmt.Sprintf("%s", r.agreeDissolvePlayers),
			"同意":   flag,
		}).Warn("玩家不能请求同意解散房间，玩家已经同意过了")
		return
	}

	r.delegate.OnRoomPlayerDissolveAgree(r, p, flag)
	//同意
	if flag {
		r.agreeDissolvePlayers[p.Id()] = p
		//判断是否全部同意
		if !r.ifAllDissolveAgree() {
			return
		}
		//发送流局
		r.liuJu()

	} else {
		//清空状态
		r.endDissolve()
	}
}

//玩家重连
func (r *Room) PlayerReconnect(p player.Player) bool {
	if r.state == RoomStateEnd {
		return false
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Debug("玩家请求重新连接房间")

	flag := r.playerManager.ReconnectPlayer(p)
	if !flag {
		return flag
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家重新连接房间")

	pl := r.playerManager.GetPlayerById(p.Id())
	r.delegate.OnRoomReconnect(r, pl)

	return flag
}

//玩家掉线
func (r *Room) PlayerDisconnect(p Player) bool {
	if r.state == RoomStateEnd {
		return false
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家请求掉线")
	flag := r.playerManager.DisconnectPlayer(p)
	if !flag {
		return flag
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家掉线")
	r.delegate.OnRoomDisconnect(r, p)
	return flag
}

//玩家自主离开
func (r *Room) PlayerLeaveTime(p Player,state bool) bool {
	if r.state == RoomStateInit {
		//return false
	}

	if r.state == RoomStateEnd {
		return false
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家请求离开")
	flag := r.playerManager.LeaveTimePlayer(p,state)
	if !flag {
		return flag
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   p.Id(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家自主离开")
	r.delegate.OnLevaeTime(r, p)
	return flag
}

//初始化日志
func (r *Room) initLogger() {
	r.logger = log.New()
	r.logger.Formatter = new(log.TextFormatter)
	r.logger.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		log.DebugLevel: fmt.Sprintf("./logs/room_%d.log", r.roomId),
		log.InfoLevel:  fmt.Sprintf("./logs/room_%d.log", r.roomId),
	}))
}

//初始化
func (r *Room) init() {
	r.initLogger()
	now := r.now() 
	r.enterInitState(now)
}

//准备开始
func (r *Room) prepareStart() {
	for _, p := range r.playerManager.Players() {
		p.Start()
	}
	r.currentRound += 1

	r.logger.WithFields(log.Fields{
		"房间状态":  r.state.String(),
		"房间回合数": r.currentRound,
		"玩家列表":  fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("准备开始")
	r.delegate.OnRoomStart(r)
	//选择庄家
	r.selectBanker()
	//发牌
	r.dealCards()
}

//选庄家
func (r *Room) selectBanker() {
	if r.currentRound == 1 {
		r.bankerPos = 0
		goto AfterBanker
	}

	if r.nextBankerPos == -1 {
		r.bankerPos += 1
		if r.bankerPos >= r.playerManager.NumPlayers() {
			r.bankerPos = 0
		}
	} else {
		r.bankerPos = r.nextBankerPos
	}

AfterBanker:
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"庄家位置": r.bankerPos,
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("选择庄家")

	//发送庄家信息 纪录日志
	r.delegate.OnRoomSelectBanker(r, r.bankerPos)
}

//发牌
func (r *Room) dealCards() {
	//洗牌
	r.deck.Shuffle()

	//发牌
	for _, p := range r.playerManager.Players() {
		for i := 0; i < initialCards; i++ {
			c := r.deck.GetFirst()
			p.Mo(c)
		}
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"剩余牌":  r.deck.Remains(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("发牌完成")

	bankerPlayer := r.playerManager.GetPlayerByPos(r.bankerPos)

	r.currentPlayer = bankerPlayer
	c := r.deck.GetFirst()

	//清除当前牌
	r.currentCards = nil
	r.currentCards = append(r.currentCards, c)
	r.currentCardPlayMethod = CardPlayMethodNormal
	r.currentCardPlayObject = CardPlayObjectSelf
	r.currentCardPlayOrder = CardPlayOrderFirst
	r.currentPlayer.Mo(c)

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"庄家牌":  c.String(),
		"剩余牌":  r.deck.Remains(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("发庄家最后一张牌")
	//发送发牌信息 纪录日志
	r.delegate.OnRoomDealCards(r)

	//纪录录像
	initLog := buildChangShaInitLog(r, r.now())
	r.logList = append(r.logList, initLog)

}

//小胡检查
func (r *Room) firstCheck() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("小胡检查")
	//小胡检查
	num := r.complementXiaoHuOperations()
	if num == 0 {
		if !r.complementOperations() {
			r.waitPlayerPlay()
			return
		}
		r.waitPlayerAction()
		return
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"小胡数":  num,
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("小胡")

	r.waitPlayerXiaoHuAction()
}

//玩家小胡
func (r *Room) PlayerXiaoHu(p Player, xht int32) {
	if r.state == RoomStateEnd {
		return
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家id": p.Id(),
	}).Info("玩家小胡")

	//判断房间状态
	if r.state != RoomStateWaitPlayerXiaoHuAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", p),
		}).Warn("玩家不能小胡")
		return
	}

	//判断是否小胡过
	if !p.IfCanXiaoHu() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", p),
		}).Warn("玩家不能小胡")
		return
	}

	xhtt := XiaoHuType(xht)
	//检查小胡是否存在
	xho := p.XiaoHuO(xhtt)
	if xho == nil {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", p),
		}).Warn("玩家不能小胡")
		return
	}
	r.delegate.OnRoomPlayerXiaoHu(r, p, xho)

	//录像
	xiaoHuLog := buildChangShaXiaoHuLog(p, xho, r.now())
	r.logList = append(r.logList, xiaoHuLog)

	//判断是否所有小胡完成了
	if r.ifXiaohuFinish() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
		}).Info("玩家小胡完成")
		if !r.complementOperations() {
			r.waitPlayerPlay()
			return
		}
		r.waitPlayerAction()
	}
}

//玩家小胡
func (r *Room) PlayerXiaoHuPass(p Player) {
	if r.state == RoomStateEnd {
		return
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家id": p.Id(),
	}).Info("玩家小胡过")

	//判断房间状态
	if r.state != RoomStateWaitPlayerXiaoHuAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", p),
		}).Warn("玩家不能小胡")
		return
	}

	//判断是否小胡过
	if !p.IfCanXiaoHu() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", p),
		}).Warn("玩家不能小胡过")
		return
	}
	p.PassXiaoHu()
	r.delegate.OnRoomPlayerXiaoHuPass(r, p)

	//判断是否所有小胡完成了
	if r.ifXiaohuFinish() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
		}).Info("玩家小胡完成")
		if !r.complementOperations() {
			r.waitPlayerPlay()
			return
		}
		r.waitPlayerAction()
	}
}

//等候玩家小胡操作
func (r *Room) waitPlayerXiaoHuAction() {

	now := r.now()
	//发送消息 等待当前玩家小胡信息
	r.enterWaitPlayerXiaoHuAction(now)
	r.delegate.OnRoomXiaoHu(r)
}

//等候玩家打牌
func (r *Room) waitPlayerPlay() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
	}).Info("等候玩家打牌")
	now := r.now()
	r.enterWaitPlayerPlay(now)
	//发送消息 等待当前玩家打牌
	r.delegate.OnRoomWaitPlayerPlay(r, r.currentPlayer)
}

//等候玩家操作
func (r *Room) waitPlayerAction() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
	}).Info("等候玩家操作")
	now := r.now()
	r.enterWaitPlayerAction(now)
	//发送消息 等待当前玩家行动
	r.delegate.OnRoomWaitPlayerAction(r)

	//纪录录像
	for _, pl := range r.RoomPlayerManager().Players() {
		if len(pl.PossibleOperations()) == 0 {
			continue
		}
		csl := buildChangShaPlayerWaitOpereationLog(r, pl, r.now())
		r.logList = append(r.logList, csl)
	}
}

//玩家操作后
func (r *Room) afterPlayerAction() {
	r.clearOperations()
}

// 摸牌 正常摸牌或补上摸牌
func (r *Room) dealCard(c *card.Card, cardPlayMethod CardPlayMethod, cardPlayOrder CardPlayOrder) {

	r.currentPlayer.Mo(c)
	r.currentCards = nil
	r.currentCards = append(r.currentCards, c)
	r.currentCardPlayMethod = cardPlayMethod
	r.currentCardPlayOrder = cardPlayOrder
	r.currentCardPlayObject = CardPlayObjectSelf

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家摸牌")

	r.delegate.OnRoomPlayerMo(r, r.currentPlayer, c)

	//纪录录像
	playerMoLog := buildChangShaPlayerMoLog(r.currentPlayer, c, r.now())
	r.logList = append(r.logList, playerMoLog)
	//填充操作
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}

	//等待玩家打牌
	r.waitPlayerPlay()

	//正在听牌了
	if r.currentPlayer.IsListen() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"当前玩家": r.currentPlayer.Id(),
			"牌的对象": r.currentCardPlayObject.String(),
			"打牌方式": r.currentCardPlayMethod.String(),
			"牌的索引": r.currentCardPlayOrder.String(),
			"当前牌":  fmt.Sprintf("%s", r.currentCards),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Info("玩家听牌")
		r.Play(r.currentPlayer, card.Value(c))
		return
	}
}

// 杠后摸牌
func (r *Room) dealCardsAfterGang(cs []*card.Card, cpo CardPlayOrder) {

	r.currentCards = cs
	for _, c := range cs {
		r.currentPlayer.AddPlayedCard(c)
	}
	r.currentPlayer.SetListen(true)
	r.currentCardPlayMethod = CardPlayMethodAfterGang
	r.currentCardPlayOrder = cpo
	r.currentCardPlayObject = CardPlayObjectSelf

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("玩家杠后摸牌")

	r.delegate.OnRoomPlayerGangMo(r, r.currentPlayer, cs)

	//纪录录像
	playerGangMoLog := buildChangShaPlayerGangMoLog(r.currentPlayer, cs, r.now())
	r.logList = append(r.logList, playerGangMoLog)

	r.checkSelfGangHu()
}

//检测杠胡
func (r *Room) checkSelfGangHu() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("检查自己杠胡")
	//填充操作
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}
	//下一个玩家
	r.checkOtherGangHu()
}

//检查其它人杠胡
func (r *Room) checkOtherGangHu() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("检查别人杠胡")
	//填充操作
	r.currentCardPlayObject = CardPlayObjectOther
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}
	//下一个玩家
	r.checkSelfGangAfterHu()
}

//检查杠不胡自己操作
func (r *Room) checkSelfGangAfterHu() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("检查自己杠不胡")
	//填充操作
	r.currentCardPlayMethod = CardPlayMethoAfterGangHu
	r.currentCardPlayObject = CardPlayObjectSelf
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}
	//下一个玩家
	r.checkOtherGangAfterHu()
}

//检查杠不胡别人操作
func (r *Room) checkOtherGangAfterHu() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": r.currentPlayer.Id(),
		"牌的对象": r.currentCardPlayObject.String(),
		"打牌方式": r.currentCardPlayMethod.String(),
		"牌的索引": r.currentCardPlayOrder.String(),
		"当前牌":  fmt.Sprintf("%s", r.currentCards),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("检查其它人杠不胡")
	//填充操作
	r.currentCardPlayMethod = CardPlayMethoAfterGangHu
	r.currentCardPlayObject = CardPlayObjectOther
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}
	//下一个玩家
	r.next()
}

// 打牌
func (r *Room) Play(player Player, cardValue int32) bool {
	if r.state == RoomStateEnd {
		return false
	}
	// 判断是不是当前玩家
	if player != r.currentPlayer {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"当前玩家": r.currentPlayer.Id(),
			"请求玩家": player.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Warn("玩家不是当前玩家")
		return false
	}
	if r.state != RoomStateWaitPlayerPlay {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"当前玩家": r.currentPlayer.Id(),
			"请求玩家": player.Id(),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Warn("房间状态不是等候玩家打牌")
		return false
	}

	c := player.PlayCard(cardValue)
	if c == nil {
		r.logger.WithFields(log.Fields{
			"房间状态":  r.state.String(),
			"当前玩家":  r.currentPlayer.Id(),
			"请求打的牌": card.NewCardValue(cardValue).String(),
			"玩家当前牌": fmt.Sprintf("%s", r.currentPlayer.Cards()),
		}).Warn("玩家没有这张牌")
		return false
	}
	//添加到下的牌中
	r.currentPlayer.AddPlayedCard(c)

	r.logger.WithFields(log.Fields{
		"房间状态":  r.state.String(),
		"当前玩家":  r.currentPlayer.Id(),
		"打的牌":   card.NewCardValue(cardValue).String(),
		"玩家当前牌": fmt.Sprintf("%s", r.currentPlayer.Cards()),
	}).Info("玩家打牌")

	//更新牌的类型
	r.currentCards = nil
	r.currentCards = append(r.currentCards, c)
	//TODO 更新牌的方式
	r.currentCardPlayMethod = CardPlayMethodNormal
	r.currentCardPlayObject = CardPlayObjectOther

	//广播打牌信息
	r.delegate.OnRoomPlayerPlayCard(r, player, c)

	//纪录录像
	playerPlayLog := buildChangShaPlayerPlayLog(player, c, r.now())
	r.logList = append(r.logList, playerPlayLog)

	//填充其他人操作
	if r.complementOperations() {
		//等候操作
		r.waitPlayerAction()
		return true
	}

	//添加到下的牌中
	//	r.currentPlayer.AddPlayedCard(c)
	//下一个玩家摸牌
	r.next()
	return true
}

//下一家 摸牌
func (r *Room) next() {
	//清除操作
	r.afterPlayerAction()

	//假如没牌了
	if r.deck.Remains() == 0 {
		r.he()
		return
	}

	if r.deck.Remains() == 1 {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"当前玩家": fmt.Sprintf("%s", r.currentPlayer),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Debug("海底牌问询")
		r.currentPlayer = r.playerManager.Next(r.currentPlayer)

		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"当前玩家": fmt.Sprintf("%s", r.currentPlayer),
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Info("海底牌问询 下一个玩家")

		//判断是不是一圈完了
		if r.currentPlayer == r.haidiPlayer {
			//和牌
			r.logger.WithFields(log.Fields{
				"房间状态": r.state.String(),
				"当前玩家": fmt.Sprintf("%s", r.currentPlayer),
				"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
			}).Info("海底牌无人要")
			r.he()
			return
		}

		if r.haidiPlayer == nil {
			r.haidiPlayer = r.currentPlayer
		}

		r.delegate.OnRoomHaiDiAsk(r, r.currentPlayer)
		return
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": fmt.Sprintf("%s", r.currentPlayer),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Debug("准备下一个玩家摸牌")

	r.currentPlayer = r.playerManager.Next(r.currentPlayer)

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": fmt.Sprintf("%s", r.currentPlayer),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("下一个玩家")

	//发牌给下一家
	c := r.deck.GetFirst()
	//普通牌
	r.dealCard(c, CardPlayMethodNormal, CardPlayOrderNormal)
}

func (r *Room) HaidiAnswer(player Player, flag bool) {
	if r.state == RoomStateEnd {
		return
	}
	//判断是不是在海底牌轮询
	if r.haidiPlayer == nil {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"请求玩家": player.Id(),
		}).Warn("当前房间没有海底轮询")
		return
	}
	if player != r.currentPlayer {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"请求玩家": player.Id(),
			"当前玩家": r.currentPlayer.Id(),
		}).Warn("海底轮询 回答不是当前玩家")
		return
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"请求玩家": fmt.Sprintf("%s", player),
		"答案":   flag,
	}).Info("海底轮询 回答")

	//不要牌
	if !flag {
		r.next()
		r.delegate.OnRoomHaiDiAnswer(r, player, nil, flag)
		return
	}

	r.haidiPlayer = nil
	//海底牌
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前玩家": fmt.Sprintf("%s", r.currentPlayer),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("海底牌")
	//发牌给下一家
	c := r.deck.GetFirst()

	r.delegate.OnRoomHaiDiAnswer(r, player, c, flag)

	r.currentCards = nil
	r.currentCards = append(r.currentCards, c)
	r.currentCardPlayMethod = CardPlayMethodNormal
	r.currentCardPlayOrder = CardPlayOrderLast
	r.currentCardPlayObject = CardPlayObjectSelf

	//填充操作
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}
	//其他人海底
	r.otherHaidi()
	//	r.dealCard(c, CardPlayMethodNormal, CardPlayOrderLast)
}

//其它人海底
func (r *Room) otherHaidi() {
	r.currentCardPlayObject = CardPlayObjectOther
	if r.complementOperations() {
		r.waitPlayerAction()
		return
	}
	r.he()
}

// 玩家操作
func (r *Room) Operate(player Player, targetIndex int32, operationType OperationType, targetCard int32, cardValues []int32) {
	if r.state == RoomStateEnd {
		return
	}
	//不是在等玩家操作
	if r.state != RoomStateWaitPlayerAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"请求玩家": fmt.Sprintf("%s", player.Id()),
		}).Warn("房间状态不是等候玩家操作")
		return
	}

	tempCard := card.NewCardValue(targetCard)
	//判断是否有这个操作
	if !player.HasOperations(targetIndex, tempCard) {
		r.logger.WithFields(log.Fields{
			"房间状态":   r.state.String(),
			"请求玩家":   fmt.Sprintf("%s", player.Id()),
			"目标卡牌索引": targetIndex,
			"目标卡牌":   tempCard,
		}).Warn("玩家没有这个操作")
		return
	}

	currentOp := player.CurrentOperation(targetIndex, tempCard)
	//已经操作过了
	if currentOp != nil {
		r.logger.WithFields(log.Fields{
			"房间状态":   r.state.String(),
			"请求玩家":   fmt.Sprintf("%s", player.Id()),
			"当前操作":   currentOp,
			"目标卡牌索引": targetIndex,
			"目标卡牌":   tempCard,
		}).Warn("玩家已经操作过了")
		return
	}

	//操作
	flag := player.Operate(targetIndex, operationType, targetCard, cardValues)
	if !flag {
		r.logger.WithFields(log.Fields{
			"房间状态":      r.state.String(),
			"请求玩家":      fmt.Sprintf("%s", player.Id()),
			"玩家可以操作的组合": fmt.Sprintf("%s", player.PossibleOperations()),
			"操作类型":      operationType.String(),
			"目标卡牌索引":    targetIndex,
			"目标卡牌":      tempCard,
			"卡牌":        fmt.Sprintf("%s", card.NewCardValues(cardValues)),
		}).Warn("玩家操作不对")
		return
	}

	//纪录录像
	playerOperationLog := buildChangShaPlayerOpereationLog(r, player, targetIndex, int32(operationType), targetCard, cardValues, r.now())
	r.logList = append(r.logList, playerOperationLog)

	r.logger.WithFields(log.Fields{
		"房间状态":   r.state.String(),
		"请求玩家":   fmt.Sprintf("%s", player),
		"目标卡牌索引": targetIndex,
		"目标卡牌":   tempCard,
		"卡牌":     fmt.Sprintf("%s", card.NewCardValues(cardValues)),
	}).Info("玩家操作")

	switch operationType {
	//过的话 判断是否所有操作完成
	case OperationTypePass:
		if !r.ifOperationFinish() {
			return
		}
		//自己摸的话,直接胡,判断是否还有剩余胡,别人胡的话判断是否有剩余胡
	case OperationTypeHu:
		if r.hasRemainHuOpreation() {
			return
		}
		//胡

		//其它操作,判断是不是最高优先级
	default:
		if !r.ifMaxPriority(operationType) {
			return
		}
	}

	var curOperationType OperationType

	ops := r.currentPlayer.CurrentOperations()
	if len(ops) != 0 {
		r.bestOpereatePlayers = append(r.bestOpereatePlayers, r.currentPlayer)
		curOperationType = ops[0].OperationType
	}

	//查找最高优先级的
	for pl := r.RoomPlayerManager().Next(r.currentPlayer); pl != r.currentPlayer; pl = r.RoomPlayerManager().Next(pl) {
		ops := pl.CurrentOperations()

		if len(ops) == 0 {
			continue
		}

		//没有最优玩家
		if len(r.bestOpereatePlayers) == 0 {

			r.bestOpereatePlayers = append(r.bestOpereatePlayers, pl)
			curOperationType = ops[0].OperationType
			continue
		}
		//更优
		if ops[0].OperationType.Priority() > curOperationType.Priority() {

			r.bestOpereatePlayers = nil
			curOperationType = ops[0].OperationType
			r.bestOpereatePlayers = append(r.bestOpereatePlayers, pl)
			continue
		}

		//多个胡
		if ops[0].OperationType.Priority() == curOperationType.Priority() && curOperationType == OperationTypeHu {
			r.bestOpereatePlayers = append(r.bestOpereatePlayers, pl)
			continue
		}

	}

	// //不是过的话,判断优先级
	// if operationType != OperationTypePass {
	// 	//没有最优先的玩家
	// 	if len(r.bestOpereatePlayers) == 0 {
	// 		r.bestOpereatePlayers = append(r.bestOpereatePlayers, player)
	// 	} else {
	// 		//判断是否优先级
	// 		if r.bestOpereatePlayers[0].CurrentOperations()[0].OperationType.Priority() < operationType.Priority() {
	// 			r.bestOpereatePlayers = nil
	// 			r.bestOpereatePlayers = append(r.bestOpereatePlayers, player)
	// 		}
	// 	}
	// }

	//操作完成
	//	if r.ifOperationFinish() {
	if len(r.bestOpereatePlayers) == 0 || curOperationType == OperationTypePass {
		r.logger.WithFields(log.Fields{
			"房间状态":   r.state.String(),
			"请求玩家":   fmt.Sprintf("%s", player),
			"目标卡牌索引": targetIndex,
			"目标卡牌":   tempCard,
			"卡牌":     fmt.Sprintf("%s", card.NewCardValues(cardValues)),
		}).Info("都不操作")

		r.afterPlayerAction()
		//海底牌
		if r.currentCardPlayOrder == CardPlayOrderLast {

			//自己 不胡
			if r.currentCardPlayObject == CardPlayObjectSelf {
				r.otherHaidi()
			} else {
				//别人也不胡
				r.he()
			}
			return
		}

		//杠后
		if r.currentCardPlayMethod == CardPlayMethodAfterGang {

			if r.currentCardPlayObject == CardPlayObjectSelf {
				r.checkOtherGangHu()
			} else {
				r.checkSelfGangAfterHu()
			}
			return
		}

		//杠不胡后
		if r.currentCardPlayMethod == CardPlayMethoAfterGangHu {
			if r.currentCardPlayObject == CardPlayObjectSelf {
				r.checkOtherGangAfterHu()
			} else {
				r.next()
			}
			return
		}

		//自己踩杠 没人操作 摸牌
		if r.currentCardPlayMethod == CardPlayMethodGang {
			//踩杠后
			r.currentPlayer = player
			if r.isGang {
				r.moAfterCaiGang(r.currentPlayer)
			} else {
				r.moAfterCaiBu(r.currentPlayer)
			}
			return
		}

		//自己摸牌后
		if r.currentCardPlayObject == CardPlayObjectSelf {
			r.currentPlayer = player
			r.waitPlayerPlay()
		} else {
			//玩家打牌后或杠后摸或自己摸牌后
			//下一个玩家
			r.next()
		}

		return
	}

	op := r.bestOpereatePlayers[0]
	switch curOperationType {
	case OperationTypeEat:
		bestOp := op.CurrentOperations()[0]
		r.eat(op, bestOp.TargetCard, card.Values(bestOp.Cards))
	case OperationTypePeng:
		bestOp := op.CurrentOperations()[0]
		r.peng(op, bestOp.TargetCard)
	case OperationTypeGang:
		bestOp := op.CurrentOperations()[0]
		r.gang(op, bestOp.Cards[0])
	case OperationTypeBu:
		bestOp := op.CurrentOperations()[0]
		r.bu(op, bestOp.Cards[0])
	case OperationTypeHu:
		bestOp := op.CurrentOperations()[0]
		r.hu(r.bestOpereatePlayers, bestOp.TargetCard)
	default:
		panic("never reach here")
	}
	//	}
}

//检查是否最优操作
func (r *Room) ifOperationFinish() bool {
	//TODO 优化是否所有人完成
	// for _, pl := range r.RoomPlayerManager().Players() {
	// 	for i := 0; i < len(r.currentCards); i++ {
	// 		tc := r.currentCards[0]
	// 		if !pl.IsOperateFinish(int32(i), tc) {
	// 			return false
	// 		}
	// 	}
	// }
	for _, pl := range r.RoomPlayerManager().Players() {
		if !pl.IfOperationsFinish(r.currentCards) {
			return false
		}
	}
	return true
}

//是否还有胡的操作
func (r *Room) hasRemainHuOpreation() bool {
	for _, pl := range r.RoomPlayerManager().Players() {
		if pl.HasRemainOperations(r.currentCards, OperationTypeHu) {
			return true
		}
	}
	return false
}

//判断是否最高优先级
func (r *Room) ifMaxPriority(ot OperationType) bool {
	for _, pl := range r.RoomPlayerManager().Players() {
		ops := pl.MaxPriorityOpreations(r.currentCards)
		if len(ops) == 0 {
			continue
		}
		if ops[0].OperationType.Priority() > ot.Priority() {
			return false
		}
	}
	return true
}

// 吃
func (r *Room) eat(player Player, targetCard *card.Card, cardValues []int32) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"卡牌":   fmt.Sprintf("%s", card.NewCardValues(cardValues)),
	}).Info("玩家准备吃")
	//不是在等玩家操作
	if r.state != RoomStateWaitPlayerAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
		}).Panic("房间状态不是等候玩家操作")
	}

	if len(cardValues) != 2 {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
			"卡牌":   fmt.Sprintf("%s", card.NewCardValues(cardValues)),
		}).Panic("吃的牌不能少于2")
	}

	//检查是否可以吃
	cards := player.Eat(cardValues, targetCard)
	if len(cards) == 0 {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
			"卡牌":   fmt.Sprintf("%s", card.NewCardValues(cardValues)),
			"吃完的牌": fmt.Sprintf("%s", cards),
		}).Panic("吃完的牌等于0")

	}

	//TODO 移除杠牌从桌面
	r.currentPlayer.RemovePlayedCard(targetCard)

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"吃完的牌": fmt.Sprintf("%s", cards),
	}).Info("玩家吃")
	//广播消息
	r.delegate.OnRoomPlayerChi(r, player, targetCard, cardValues)

	//纪录录像
	playerChiLog := buildChangShaPlayerChiLog(player, cardValues, r.now())
	r.logList = append(r.logList, playerChiLog)

	r.afterPlayerAction()

	r.currentPlayer = player

	r.currentCardPlayMethod = CardPlayMethodAfterChiPeng
	r.currentCardPlayObject = CardPlayObjectSelf
	//填充操作
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}

	//进入打牌状态
	r.waitPlayerPlay()

}

//碰
func (r *Room) peng(player Player, targetCard *card.Card) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"当前牌":  r.currentCards[0],
	}).Info("玩家准备碰")

	//不是在等玩家操作
	if r.state != RoomStateWaitPlayerAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
		}).Panic("房间状态不是等候玩家操作")
	}

	//检查是否可以吃
	cards := player.Peng(targetCard)
	if len(cards) == 0 {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
		}).Panic("玩家碰完的牌等于0")
	}

	//TODO 移除杠牌从桌面
	r.currentPlayer.RemovePlayedCard(targetCard)

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"碰完的牌": fmt.Sprintf("%s", cards),
	}).Info("玩家碰")

	//广播消息
	r.delegate.OnRoomPlayerPeng(r, player, targetCard)
	//纪录录像
	playerPengLog := buildChangShaPlayerPengLog(player, card.Value(targetCard), r.now())
	r.logList = append(r.logList, playerPengLog)

	r.afterPlayerAction()

	r.currentPlayer = player

	r.currentCardPlayMethod = CardPlayMethodAfterChiPeng
	r.currentCardPlayObject = CardPlayObjectSelf
	//填充操作
	if r.complementOperations() {
		//等待玩家操作
		r.waitPlayerAction()
		return
	}

	//进入打牌状态
	r.waitPlayerPlay()
}

// 补牌
func (r *Room) bu(player Player, targetCard *card.Card) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"当前牌":  r.currentCards[0],
	}).Info("玩家准备补")
	//不是在等玩家操作
	if r.state != RoomStateWaitPlayerAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
		}).Panic("房间状态不是等候玩家操作")
	}

	//检查是否可以补
	flag, gt := player.Bu(targetCard, r.currentCardPlayObject, r.currentCardPlayMethod)
	if !flag {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
		}).Panic("玩家不能补")
	}

	//TODO 移除杠牌从桌面
	if r.currentCardPlayObject == CardPlayObjectOther {
		r.currentPlayer.RemovePlayedCard(targetCard)
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"补的类型": gt.String(),
	}).Info("玩家补")

	//发送消息
	r.delegate.OnRoomPlayerBu(r, player, targetCard, gt)
	//纪录录像
	playerBuLog := buildChangShaPlayerBuLog(player, card.Value(targetCard), int32(gt), r.now())
	r.logList = append(r.logList, playerBuLog)

	r.afterPlayerAction()
	r.currentPlayer = player
	//判断是否别人胡牌
	if gt == GangTypeDiPaiGang {
		r.currentCards = nil
		r.currentCards = append(r.currentCards, targetCard)
		r.currentCardPlayObject = CardPlayObjectOther
		r.currentCardPlayMethod = CardPlayMethodGang
		r.isGang = false
		if r.complementOperations() {
			r.waitPlayerAction()
			return
		}
	}

	r.moAfterCaiBu(r.currentPlayer)
}

//杠牌
func (r *Room) gang(player Player, targetCard *card.Card) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"当前牌":  r.currentCards[0],
	}).Info("玩家准备杠")
	//不是在等玩家操作
	if r.state != RoomStateWaitPlayerAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
		}).Panic("房间状态不是等候玩家操作")
	}

	//检查是否可以杠
	flag, gt := player.Gang(targetCard, r.currentCardPlayObject, r.currentCardPlayMethod)
	if !flag {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", player),
		}).Panic("玩家不能杠")
	}

	//TODO 移除杠牌从桌面
	if r.currentCardPlayObject == CardPlayObjectOther {
		r.currentPlayer.RemovePlayedCard(targetCard)
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家":   fmt.Sprintf("%s", player),
		"杠的类型": gt.String(),
	}).Info("玩家杠")

	r.delegate.OnRoomPlayerGang(r, player, targetCard, gt)
	//纪录录像
	playerGangLog := buildChangShaPlayerGangLog(player, card.Value(targetCard), int32(gt), r.now())
	r.logList = append(r.logList, playerGangLog)

	r.afterPlayerAction()

	r.currentPlayer = player
	//判断是否别人胡牌
	if gt == GangTypeDiPaiGang {
		r.currentCards = nil
		r.currentCards = append(r.currentCards, targetCard)
		r.currentCardPlayObject = CardPlayObjectOther
		r.currentCardPlayMethod = CardPlayMethodGang
		r.isGang = true
		if r.complementOperations() {
			r.waitPlayerAction()
			return
		}
	}
	r.moAfterCaiGang(r.currentPlayer)
}

//踩杠后没人操作
func (r *Room) moAfterCaiGang(player Player) {

	tempCards := make([]*card.Card, 0, 2)
	//杠后摸牌
	for i := 0; i < 2; i++ {
		if r.deck.Remains() == 0 {
			break
		}
		c := r.deck.GetLast()
		tempCards = append(tempCards, c)
	}

	//普通牌
	r.dealCardsAfterGang(tempCards, CardPlayOrderNormal)
}

//踩杠后没人操作
func (r *Room) moAfterCaiBu(player Player) {
	//杠后摸牌
	c := r.deck.GetFirst()

	//普通牌
	r.dealCard(c, CardPlayMethodNormal, CardPlayOrderNormal)
}

//和
func (r *Room) he() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("没牌了")
	r.nextBankerPos = -1
	r.settleXiaoHu()
	r.delegate.OnRoomPlayersSettle(r, nil, 0, nil)
	r.settle()
}

//糊
func (r *Room) hu(players []Player, targetCard *card.Card) {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"胡的玩家": fmt.Sprintf("%s", players),
		"当前牌":  r.currentCards[0],
	}).Info("玩家准备胡")
	//	now := time.Now().UnixNano() / int64(time.Millisecond)
	//不是在等玩家操作
	if r.state != RoomStateWaitPlayerAction {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
		}).Panic("房间状态不是等候玩家操作")
	}

	// //结算小胡
	// for _, p := range r.playerManager.Players() {
	// 	if len(p.CurrentXiaoHus()) != 0 {
	// 		for _, p2 := range r.playerManager.Players() {
	// 			if p != p2 {
	// 				p2.SubScore(xiaoHuScore)
	// 			} else {
	// 				p2.AddScore(xiaoHuScore * (r.playerManager.NumPlayers() - 1))
	// 			}
	// 		}
	// 	}
	// }

	zhuaNiaoPos := 0
	// if r.customConfig.ZhuangXian {
	// 	zhuaNiaoPos = r.bankerPos
	// } else {
	if len(players) >= 2 {
		zhuaNiaoPos = r.currentPlayer.Position()
	} else {
		zhuaNiaoPos = players[0].Position()
	}
	//}
	//抓鸟
	r.zhuaNiao(zhuaNiaoPos)
	//结算小胡
	r.settleXiaoHu()

	//庄家
	//bankerPlayer := r.RoomPlayerManager().GetPlayerByPos(r.bankerPos)

	//检查是否胡牌
	for _, player := range players {
		for _, op := range player.CurrentOperations() {
			if op.OperationType == OperationTypeHu {
				result := player.Hu(op.TargetCard, r.currentCardPlayObject, r.currentCardPlayMethod, r.currentCardPlayOrder)
				if result == 0 {
					r.logger.WithFields(log.Fields{
						"房间状态": r.state.String(),
						"玩家":   fmt.Sprintf("%s", player),
					}).Panic("玩家不能胡")
				}

				tempScore := 0
				_, dh := r.xiaoHuAndDaHuForResult(result)
				//只有小胡
				if dh == 0 {
					//自摸
					if r.currentCardPlayObject == CardPlayObjectSelf {
						tempScore += xiaoHuZiMoScore
					} else {
						tempScore += xiaoHuJiePaoScore
					}
				} else {
					//自摸
					if r.currentCardPlayObject == CardPlayObjectSelf {
						tempScore += daHuZiMoScore * dh
					} else {
						tempScore += daHuZiMoScore * dh
					}
				}
				//庄闲算分
				if r.customConfig.ZhuangXian {
					if dh == 0 {
						tempScore += zhuanExtraScore
					} else {
						tempScore += zhuanExtraScore * dh
					}
				}
				//自摸
				if r.currentCardPlayObject == CardPlayObjectSelf {

					if dh != 0 {
						//大胡自摸
						player.AddSettleData(SettleTypeDaHuZiMo)
					} else {
						//小胡自摸
						player.AddSettleData(SettleTypeXiaoHuZiMo)
					}

					// //庄闲算分
					// if r.customConfig.ZhuangXian {
					// 	if player == bankerPlayer {
					// 		if dh == 0 {
					// 			tempScore += zhuanExtraScore
					// 		} else {
					// 			tempScore += zhuanExtraScore * dh
					// 		}
					// 	}
					// }

					totalScore := 0
					//扣除3家分数
					for _, p2 := range r.RoomPlayerManager().Players() {
						if p2 != player {
							// tempScore2 := tempScore
							// if r.customConfig.ZhuangXian {
							// 	if p2 == bankerPlayer {
							// 		if dh == 0 {
							// 			tempScore2 += zhuanExtraScore
							// 		} else {
							// 			tempScore2 += zhuanExtraScore * dh
							// 		}
							// 	}
							// }

							niaoFan := player.ZhongNiao() + p2.ZhongNiao()
							if !r.customConfig.ZhuangXian {
								niaoFan += 1
							} else {
								niaoFan = int32(math.Exp2(float64(niaoFan)))
							}
							p2.SubScore(tempScore * int(niaoFan))
							totalScore += tempScore * int(niaoFan)
						}
					}

					player.AddScore(totalScore)
				} else {
					if dh != 0 {
						//大胡接炮
						player.AddSettleData(SettleTypeDaHuJiePao)
						r.currentPlayer.AddSettleData(SettleTypeDaHuDianPao)
					} else {
						//小胡接炮
						player.AddSettleData(SettleTypeXiaoHuJiePao)
						r.currentPlayer.AddSettleData(SettleTypeXiaoHuDianPao)
					}
					niaoFan := player.ZhongNiao() + r.currentPlayer.ZhongNiao()
					if !r.customConfig.ZhuangXian {
						niaoFan += 1
					} else {
						niaoFan = int32(math.Exp2(float64(niaoFan)))
					}

					//庄闲算分
					// if r.customConfig.ZhuangXian {
					// 	if player == bankerPlayer || r.currentPlayer == bankerPlayer {
					// 		if dh == 0 {
					// 			tempScore += zhuanExtraScore
					// 		} else {
					// 			tempScore += zhuanExtraScore * dh
					// 		}
					// 	}
					// }

					player.AddScore(tempScore * int(niaoFan))
					r.currentPlayer.SubScore(tempScore * int(niaoFan))
				}
			}
		}
	}

	//设置下局庄家位置
	if len(players) > 1 {
		r.nextBankerPos = r.currentPlayer.Position()
	} else {
		r.nextBankerPos = players[0].Position()
	}

	//纪录录像
	for _, player := range players {
		huLog := buildChangShaPlayerHuLog(player, card.Value(targetCard), r.now())
		r.logList = append(r.logList, huLog)
	}

	r.delegate.OnRoomPlayersSettle(r, players, r.currentPlayer.Id(), targetCard)

	//结算
	r.settle()
}

func (r *Room) endDissolve() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
	}).Info("结束申请解散")
	r.dissolveStartTime = 0
	r.agreeDissolvePlayers = make(map[int64]Player)
	r.applyDissolvePlayer = nil
}

//流局
func (r *Room) liuJu() {
	if r.state == RoomStateEnd {
		return
	}
	r.endDissolve()
	if r.state != RoomStateSettle {

		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
		}).Info("房间流局")
		//发送流局
		r.settleXiaoHu()
		r.delegate.OnRoomLiuJu(r)
	}

	//发送总结算
	r.end(true)
}

//抓鸟
func (r *Room) zhuaNiao(startPos int) {
	if r.customConfig.ZhuaNiao == 0 {
		return
	}

	//获取鸟牌
	if r.deck.Remains() == 0 {
		for i := 0; i < r.customConfig.ZhuaNiao; i++ {
			r.niaoPais = append(r.niaoPais, r.currentCards[0])
		}
	} else {
		for i := 0; i < r.customConfig.ZhuaNiao; i++ {
			if r.deck.Remains() == 0 {
				r.niaoPais = append(r.niaoPais, r.niaoPais[len(r.niaoPais)-1])
			} else {
				r.niaoPais = append(r.niaoPais, r.deck.GetFirst())
			}
		}
	}

	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"鸟牌":   r.niaoPais,
	}).Info("抓鸟中")

	for _, niaoPai := range r.niaoPais {
		pos := ((int(niaoPai.CardValue) - 1) + startPos) % r.playerManager.NumPlayers()
		pl := r.playerManager.GetPlayerByPos(pos)
		pl.ZhuaNiao()
		r.niaoPaiPlayerIds = append(r.niaoPaiPlayerIds, pl.Id())
	}

}

//结算抓鸟加小胡
func (r *Room) settleXiaoHu() {
	//结算小胡
	for _, p := range r.playerManager.Players() {
		tempNumXiaoHus := len(p.CurrentXiaoHus())
		if tempNumXiaoHus != 0 {
			tempTotalScore := 0
			for _, p2 := range r.playerManager.Players() {
				if p != p2 {
					// niaoFan := p.ZhongNiao() + p2.ZhongNiao()
					// if !r.customConfig.ZhuangXian {
					// 	niaoFan += 1
					// } else {
					// 	niaoFan = int32(math.Exp2(float64(niaoFan)))
					// }
					tempScore := xiaoHuScore * tempNumXiaoHus //* int(niaoFan)
					p2.SubScore(tempScore)
					tempTotalScore += tempScore
				}
			}
			p.AddScore(tempTotalScore)
		}
	}
}

//结算
func (r *Room) settle() {
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"当前局数": r.currentRound,
		"总局数":  r.totalRound,
		"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
	}).Info("房间结算中")

	now := time.Now().UnixNano() / int64(time.Millisecond)

	//结算
	for _, p := range r.playerManager.Players() {
		p.Settle()
	}

	r.delegate.OnRoomSettle(r)

	//进入结算页面
	r.enterSettle(now)

	r.afterPlayerAction()

	//判断是否结束 
	if r.totalRound != 0 {
		if r.currentRound == r.totalRound { 
			r.end(true) 
		}
	} else {
		for _, pl := range r.playerManager.Players() {
			if pl.Score() <= 0 {
				r.end(true)
				break
			}
		}
	} 

}

//补充小胡操作
func (r *Room) complementXiaoHuOperations() int {
	totalNum := r.playerManager.NumPlayers()
	num := 0
	//顺序补充
	for i := 0; i < totalNum; i++ {
		pos := r.bankerPos + i
		if pos >= totalNum {
			pos -= totalNum
		}
		pl := r.playerManager.GetPlayerByPos(pos)
		n := pl.ComplementXiaoHu()
		num += n
	}
	return num
}

//清楚玩家操作
func (r *Room) clearOperations() {
	for _, p := range r.playerManager.Players() {
		p.ClearOperations()
	}
	r.bestOpereatePlayers = nil
}

//填充其它人的操作
func (r *Room) complementOperations() bool {

	// if r.currentCardPlayMethod == CardPlayMethodAfterGang {
	// 	found := false
	// 	found = found || r.complementSelfOperations()
	// 	found = found || r.complementOtherOperations()
	// 	return found
	// }
	switch r.currentCardPlayObject {
	//别人打的或者是明杠
	case CardPlayObjectOther:
		{
			//补充除自己以外的操作
			return r.complementOtherOperations()
		}
	//自摸的
	case CardPlayObjectSelf:
		{
			//补充自己操作
			return r.complementSelfOperations()
		}
	default:
		panic("never reach here")
	}

	//return true
}

//填充自己操作
func (r *Room) complementSelfOperations() bool {
	log.Println("填充自己操作")
	hasCards := false
	if r.deck.Remains() > 0 {
		hasCards = true
	}
	r.currentPlayer.ComplementOperations(r.currentCards, r.currentCardPlayObject, r.currentCardPlayMethod, r.currentCardPlayOrder, false, hasCards)
	if len(r.currentPlayer.PossibleOperations()) == 0 {
		return false
	}
	return true
}

//填充别人操作 按顺序
func (r *Room) complementOtherOperations() bool {
	log.Println("填充别人操作")
	tempCurrentPlayer := r.currentPlayer
	found := false
	next := true
	hasCards := false
	if r.deck.Remains() > 0 {
		hasCards = true
	}
	for {
		tempNextPlayer := r.playerManager.Next(tempCurrentPlayer)
		if tempNextPlayer == r.currentPlayer {
			break
		}
		tempNextPlayer.ComplementOperations(r.currentCards, r.currentCardPlayObject, r.currentCardPlayMethod, r.currentCardPlayOrder, next, hasCards)
		if len(tempNextPlayer.PossibleOperations()) != 0 {
			log.Printf("玩家[%d],可以操作[%s]", tempNextPlayer.Id(), tempNextPlayer.PossibleOperations())
			found = true
		}
		tempCurrentPlayer = tempNextPlayer
		next = false
	}
	log.Println("填充别人操作完成")
	return found
}

//清理桌面
func (r *Room) clear() {
	r.logger.Info("清理桌面")
	r.niaoPais = nil
	r.niaoPaiPlayerIds = nil
	r.haidiPlayer = nil
	r.isGang = false
	for _, pl := range r.playerManager.Players() {
		pl.Clear()
	}
	r.delegate.OnRoomClear(r)
}

//玩家准备
func (r *Room) PrepareStart(p Player) {
	if r.state == RoomStateEnd {
		return
	}
	r.logger.WithFields(log.Fields{
		"房间状态": r.state.String(),
		"玩家id": p.Id(),
	}).Info("玩家准备")

	r.logList = nil
	if !p.IfCanPrepare() {
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"玩家":   fmt.Sprintf("%s", p),
		}).Warn("玩家不能准备")
		return
	}
	p.Prepare()
	r.delegate.OnRoomPlayerStart(r, p)
	now := r.now()
	//判断是否全部准备了
	if r.ifCanStart() {
		r.clear()
		r.enterWaitState(now)
	}
}

//结束
func (r *Room) end(start bool) {
	if start { 
		r.logger.WithFields(log.Fields{
			"房间状态": r.state.String(),
			"当前局数": r.currentRound,
			"总局数":  r.totalRound,
			"玩家列表": fmt.Sprintf("%s", r.playerManager.Players()),
		}).Info("房间总结算")

		for _, pl := range r.playerManager.Players() {
			pl.End()
		}

		now := r.now()
		//进入结束状态
		r.enterEnd(now)
	}
	//回调
	r.delegate.OnRoomEnd(r,start)
}

//当前时间
func (r *Room) now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (r *Room) xiaoHuAndDaHuForResult(result int32) (xiaoHu int, daHu int) {
	if result&int32(HandCardTypePingHu) != 0 {
		xiaoHu += 1
	}
	if result&int32(HandCardTypePengPengHu) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeJiangJiangHu) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeQingYiSe) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeHaiDiLaoYue) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeHaiDiPao) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeQiXiaoDui) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeHaoHuaQiXiaoDui) != 0 {
		daHu += 2
	}

	if result&int32(HandCardTypeGangShangKaiHua) != 0 {
		daHu += 1
	}

	if result&int32(HandCardTypeQiangGangHu) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeGangShangPao) != 0 {
		daHu += 1
	}

	if result&int32(HandCardTypeQuanQiuRen) != 0 {
		daHu += 1
	}
	if result&int32(HandCardTypeShuangHaoHuaQiXiaoDui) != 0 {
		daHu += 2
	}

	if result&int32(HandCardTypeTianHu) != 0 {
		daHu += 1
	}

	if result&int32(HandCardTypeDiHu) != 0 {
		daHu += 1
	}
	return
}
