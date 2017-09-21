package handler

import (
	"fmt"

	"game/basic/pb"
	"game/mahjong/client/client"
	changshapb "game/mahjong/pb/changsha"

	"github.com/golang/protobuf/proto"
)

func handleSettle(c *client.Client, msg *pb.Message) error {

	val, err := proto.GetExtension(msg, changshapb.E_GcPlayerHu)
	if err != nil {
		return err
	}
	_, ok := val.(*changshapb.GCPlayerHu)
	if !ok {
		panic("no gc player play card")
	}

	fmt.Println("结算中")

	// huList := gcPlayerHu.GetHuList()
	// if len(huList) == 0 {
	// 	fmt.Println("和牌")
	// 	c.PrepareStart()
	// 	return nil
	// }

	// tc := card.NewCardValue(gcPlayerHu.GetCard())

	// for _, info := range gcPlayerHu.GetSettlePlayerList() {
	// 	fmt.Printf("玩家[%d],手牌[%s],吃碰杠[%s],输赢[%d]", info.GetPlayerId(), cardsFromInts(info.GetCardList()), composeListFromComposeInfoList(info.GetComposeList()), info.GetScore())
	// 	for _, hu := range huList {
	// 		if hu.GetPlayerId() == info.GetPlayerId() {
	// 			fmt.Printf("胡牌[%s],胡牌种类[%s]", tc, changsha.HandCardType(hu.GetHuType()))
	// 		}
	// 	}

	// }
	c.PrepareStart()
	return nil
}
