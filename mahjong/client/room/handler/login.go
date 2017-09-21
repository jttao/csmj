package handler

import (
	"log"

	"github.com/golang/protobuf/proto"
	"game/basic/pb"
	"game/mahjong/client/client"
	loginpb "game/mahjong/pb/login"
)

func handleLogin(c *client.Client, msg *pb.Message) error {
	log.Println("login")
	val, err := proto.GetExtension(msg, loginpb.E_GcLogin)
	if err != nil {
		return err
	}
	gcLogin, ok := val.(*loginpb.GCLogin)
	if !ok {
		panic("no gc login")
	}
	//	c.Player().SetId(gcLogin.GetPlayerId())
	c.SetId(gcLogin.GetPlayerId())
	return nil
}
