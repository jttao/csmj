package handler

import (
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"
)

func handleDealCards(c *client.Client, msg *pb.Message) error {
	log.Println("deal cards")

	val, err := proto.GetExtension(msg, changshapb.E_GcDealCards)
	if err != nil {
		return err
	}
	gcDealCards, ok := val.(*changshapb.GCDealCards)
	if !ok {
		panic("gc deal cards")
	}

	cards := cardsFromInts(gcDealCards.Cards)
	r := c.Room()
	r.DealCards(c.Id(), cards)

	return nil
}
