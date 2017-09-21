package mahjong

import (
	"fmt"
	"game/mahjong/server/player"
	"game/session"
	"time"

	"game/basic/pb"
	mahjongutil "game/mahjong/server/util"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

//对话启动
type SessionOpener struct {
	m *Mahjong
}

func (so *SessionOpener) Handle(s session.Session) error {
	nctx := WithMahjong(s.Context(), so.m)
	s.SetContext(nctx)

	return nil
}

func NewSessionOpener(m *Mahjong) *SessionOpener {
	sp := &SessionOpener{}
	sp.m = m
	return sp
}

//验证过期时间 和ping的
func AuthTimeoutMiddleware(sh session.SessionHandler) session.SessionHandler {
	return session.SessionHandlerFunc(func(s session.Session) error {
		go func() {
			<-time.After(time.Second * 10)
			pl := player.PlayerInContext(s.Context())
			if pl == nil {
				log.WithFields(
					log.Fields{
						"sessionId": s.Id()}).Warn("验证超时")
				err := mahjongutil.CloseWithError(s, int32(ErrorCodeAuthTimeout))
				if err != nil {
					log.WithFields(
						log.Fields{
							"sessionId": s.Id(),
							"error":     err,
						}).Warn("验证超时,关闭错误")
				}
			}
		}()
		return sh.Handle(s)
	})
}

//对话关闭
type SessionCloser struct {
}

func (sc *SessionCloser) Handle(s session.Session) error {
	pl := player.PlayerInContext(s.Context())
	if pl == nil {
		return nil
	}
	pl.Stop()
	mahjongContext := MahjongInContext(s.Context())
	mahjongContext.PlayerManager.RemovePlayer(pl)

	r := mahjongContext.RoomManager.GetRoomById(pl.RoomId())
	if r == nil {
		return nil
	}

	roomPlayer := r.RoomPlayerManager().GetPlayerById(pl.Id())
	if roomPlayer == nil {
		return nil
	}

	r.PlayerDisconnect(roomPlayer)
	return nil
}

func NewSessionCloser() *SessionCloser {
	sc := &SessionCloser{}
	return sc
}

//对话处理器
type SessionProcessor struct {
}

func (sp *SessionProcessor) Handle(s session.Session, msg []byte) error {
	log.WithFields(
		log.Fields{
			"sessionId": s.Id(),
			"msg":       msg,
		}).Debug("对话处理器,接收消息")
	m := &pb.Message{}
	err := proto.Unmarshal(msg, m)
	if err != nil {
		return err
	}
	mah := MahjongInContext(s.Context())
	m2 := NewMessage(s, m)

	pl := player.PlayerInContext(s.Context())
	if pl == nil {
		//放不同的消息进去不同的处理器
		mah.GlobalProcessor.Receive(m2)
	} else {
		rId := pl.RoomId()
		r := mah.RoomManager.GetRoomById(rId)
		if r == nil {
			log.WithFields(
				log.Fields{
					"玩家id": pl.Id(),
				}).Warn("对话处理器,玩家不在房间")
			return nil
		}
		rp := RoomProcessorInContext(r.Context())
		rp.Receive(m2)
	}
	return nil
}

func NewSessionProcessor() *SessionProcessor {
	sp := &SessionProcessor{}
	return sp
}

type Message struct {
	s   session.Session
	msg *pb.Message
}

func (m *Message) Session() session.Session {
	return m.s
}

func (m *Message) Msg() *pb.Message {
	return m.msg
}

func NewMessage(s session.Session, msg *pb.Message) *Message {
	m := &Message{
		s:   s,
		msg: msg,
	}
	return m
}

type MessageHandler interface {
	Handle(s session.Session, msg *pb.Message) error
}

type MessageHandlerFunc func(s session.Session, msg *pb.Message) error

func (mhf MessageHandlerFunc) Handle(s session.Session, msg *pb.Message) error {
	return mhf(s, msg)
}

type Dispatcher struct {
	handlerMap map[int32]MessageHandler
}

func (d *Dispatcher) Register(messageType int32, h MessageHandler) error {
	d.handlerMap[messageType] = h
	return nil
}

func (d *Dispatcher) Handle(s session.Session, msg *pb.Message) error {
	h, exist := d.handlerMap[msg.GetMessageType()]
	if !exist {
		return fmt.Errorf("no exist handler for message type %d", msg.GetMessageType())
	}

	return h.Handle(s, msg)
}

func NewDispatch() *Dispatcher {
	d := &Dispatcher{}
	d.handlerMap = make(map[int32]MessageHandler)
	return d
}
