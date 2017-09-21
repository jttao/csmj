package client

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/card"
	"game/mahjong/changsha"
	"game/mahjong/client/room"
	changshapb "game/mahjong/pb/changsha"
	messagetypepb "game/mahjong/pb/messagetype"
	"golang.org/x/net/websocket"
)

type Client struct {
	id   int64
	conn *websocket.Conn
	r    room.Room
	pl   room.Player
}

func (c *Client) Id() int64 {
	return c.id
}

func (c *Client) SetId(id int64) {
	c.id = id
}

func (c *Client) SetRoom(r room.Room) {
	c.r = r
}

func (c *Client) Room() room.Room {
	return c.r
}

func (c *Client) SetPlayer(p room.Player) {
	c.pl = p
}

func (c *Client) Player() room.Player {
	return c.pl
}

func (c *Client) PrepareStart() {
	time.Sleep(2 * time.Second)
	log.Println("client prepare start")

	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_CGPlayerStartType)
	msg.MessageType = &msgType
	cgPlayerStart := &changshapb.CGPlayerStart{}

	err := proto.SetExtension(msg, changshapb.E_CgPlayerStart, cgPlayerStart)
	if err != nil {
		panic(err)
	}
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	err = websocket.Message.Send(c.conn, msgByte)
	if err != nil {
		panic(err)
	}
}

func (c *Client) PrepareXiaoHu(xht int32) {
	log.Println("client prepare xiao hu")

	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_CGXiaoHuType)
	msg.MessageType = &msgType
	cgXiaoHu := &changshapb.CGXiaoHu{}
	cgXiaoHu.XiaoHu = &xht

	err := proto.SetExtension(msg, changshapb.E_CgXiaoHu, cgXiaoHu)
	if err != nil {
		panic(err)
	}
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	err = websocket.Message.Send(c.conn, msgByte)
	if err != nil {
		panic(err)
	}
}

func (c *Client) PreaprePlay() {
	log.Println("client prepare play")
	time.Sleep(1 * time.Second)
	log.Printf("玩家当前牌[%s]\n", c.Player().Cards())
	tc := c.findBestCard()
	log.Printf("client play card %s\n", tc)
	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_CGPlayerPlayCardType)
	msg.MessageType = &msgType
	cgPlayerPlayCard := &changshapb.CGPlayerPlayCard{}
	cv := card.Value(tc)
	cgPlayerPlayCard.Card = &cv
	err := proto.SetExtension(msg, changshapb.E_CgPlayerPlayCard, cgPlayerPlayCard)
	if err != nil {
		panic(err)
	}
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	err = websocket.Message.Send(c.conn, msgByte)
	if err != nil {
		panic(err)
	}
}

func (c *Client) PrepareDissolve() {
	log.Println("client prepare dissolve")
	time.Sleep(1 * time.Second)

	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_CGPlayerDissolveAgreeType)
	msg.MessageType = &msgType
	cgPlayerDissolveAgree := &changshapb.CGPlayerDissolveAgree{}
	agree := true
	cgPlayerDissolveAgree.Agree = &agree
	err := proto.SetExtension(msg, changshapb.E_CgPlayerDissolveAgree, cgPlayerDissolveAgree)
	if err != nil {
		panic(err)
	}
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	err = websocket.Message.Send(c.conn, msgByte)
	if err != nil {
		panic(err)
	}
}

func (c *Client) PrepareHaidiAnswer() {
	log.Println("client prepare haidi")
	time.Sleep(1 * time.Second)

	msg := &pb.Message{}
	msgType := int32(messagetypepb.MessageType_CGHaiDiAnswerType)
	msg.MessageType = &msgType
	cgHaiDiAnswer := &changshapb.CGHaiDiAnswer{}
	flag := true
	cgHaiDiAnswer.Flag = &flag
	err := proto.SetExtension(msg, changshapb.E_CgHaiDiAnswer, cgHaiDiAnswer)
	if err != nil {
		panic(err)
	}
	msgByte, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	err = websocket.Message.Send(c.conn, msgByte)
	if err != nil {
		panic(err)
	}
}

func (c *Client) PreparePlayerOpereations(ops []*changshapb.Operation) {
	log.Printf("possible opereations %s\n", ops)
	time.Sleep(1 * time.Second)
	bestOps := c.findBestOpereations(ops)
	log.Printf("client operate %s\n", bestOps)

	for _, op := range bestOps {
		msg := &pb.Message{}
		msgType := int32(messagetypepb.MessageType_CGPlayerOperateType)
		msg.MessageType = &msgType
		cgPlayerOpreate := &changshapb.CGPlayerOperate{}

		cgPlayerOpreate.Operation = op

		err := proto.SetExtension(msg, changshapb.E_CgPlayerOpreate, cgPlayerOpreate)
		if err != nil {
			panic(err)
		}
		msgByte, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
		err = websocket.Message.Send(c.conn, msgByte)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Client) findBestOpereations(ops []*changshapb.Operation) []*changshapb.Operation {
	var bestOps []*changshapb.Operation
	var bestOpMap = make(map[int32]*changshapb.Operation)
	for _, op := range ops {
		top, exist := bestOpMap[op.GetTargetIndex()]
		if !exist {
			bestOpMap[op.GetTargetIndex()] = op
			continue
		}

		if changsha.OperationType(top.GetOperationType()).Priority() > changsha.OperationType(op.GetOperationType()).Priority() {
			bestOpMap[op.GetTargetIndex()] = top
		} else if changsha.OperationType(top.GetOperationType()).Priority() == changsha.OperationType(op.GetOperationType()).Priority() {
			if changsha.OperationType(op.GetOperationType()) == changsha.OperationTypeGang {
				bestOpMap[op.GetTargetIndex()] = top
			}
		}
	}
	for _, op := range bestOpMap {
		bestOps = append(bestOps, op)
	}
	return bestOps
}

func (c *Client) findBestCard() *card.Card {
	tc := c.pl.Cards()[0]
	return tc
}

func (c *Client) Close() {
	c.conn.Close()
}

func NewClient(id int64, conn *websocket.Conn) *Client {
	c := &Client{conn: conn}
	c.id = id
	return c
}
