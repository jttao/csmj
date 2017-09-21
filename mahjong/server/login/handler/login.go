package handler

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"game/mahjong/changsha"
	loginpb "game/mahjong/pb/login"
	messagetypepb "game/mahjong/pb/messagetype"
	"game/mahjong/server/mahjong"
	"game/mahjong/server/player"

	"game/mahjong/server/room"

	"game/session"

	"game/basic/pb"

	"github.com/golang/protobuf/proto"
)

func HandleLogin(s session.Session, msg *pb.Message) error {
	log.Debug("处理登陆消息")
	pl := player.PlayerInContext(s.Context())
	if pl != nil {
		//TODO 发送错误信息
		log.WithFields(
			log.Fields{
				"playerId": pl.Id(),
			}).Warn("玩家重复登陆")
		return nil
	}

	mahjongContext := mahjong.MahjongInContext(s.Context())

	loginMsg, err := proto.GetExtension(msg, loginpb.E_CgLogin)
	if err != nil {
		//TODO 发送异常信息
		log.WithFields(
			log.Fields{
				"sessionId": s.Id(),
				"error":     err.Error(),
			}).Error("玩家登陆消息解析错误")
		return err
	}

	cgLogin := loginMsg.(*loginpb.CGLogin)
	token := cgLogin.GetToken()
	playerId, err := mahjongContext.UserService.Verify(token)
	if err != nil {
		//TODO 发送异常信息
		log.WithFields(
			log.Fields{
				"playerId": playerId,
				"error":    err.Error(),
			}).Error("玩家验证token失败")
		return err
	}

	user, err := mahjongContext.UserService.GetUserById(playerId)
	if err != nil {
		//TODO 发送异常信息
		log.WithFields(
			log.Fields{
				"playerId": playerId,
				"error":    err.Error(),
			}).Error("获取用户基本信息失败")
		return err
	}

	pl = mahjongContext.PlayerManager.GetPlayerById(playerId)
	if pl != nil {
		log.WithFields(
			log.Fields{
				"playerId": playerId,
			}).Error("玩家已经在服务器内")
		pl.Session().Close()
		s.Close()
		return nil
	}

	roomId, _, _, _, maxPlayer, round, roomConfig,forbidIp, err := mahjongContext.RoomManageClient.Query(playerId)
	if err != nil {
		//TODO 发送异常信息
		log.WithFields(
			log.Fields{
				"playerId": playerId,
				"error":    err.Error(),
			}).Error("玩家查询服务器数据失败")
		return err
	}
	//验证房间号
	if roomId == 0 {
		//TODO 发送错误信息
		log.WithFields(
			log.Fields{
				"playerId": playerId,
			}).Warn("玩家房间号为0")
		return fmt.Errorf("join room id [%d]", roomId)
	}
	//获取房间
	r := mahjongContext.RoomManager.GetRoomById(roomId)
	if r == nil {
		log.WithFields(
			log.Fields{
				"playerId": playerId,
				"roomId":   roomId,
			}).Debug("房间不存在")
		cf := &changsha.CustomRoomConfig{}
		err := json.Unmarshal([]byte(roomConfig), cf)
		if err != nil {
			//TODO 发送异常信息
			log.WithFields(
				log.Fields{
					"playerId": playerId,
					"error":    err.Error(),
				}).Error("解析房间配置数据错误")
			return err
		}
		
		r = changsha.NewRoom(mahjongContext.ServerCfg.Room, cf, maxPlayer, round, roomId, forbidIp, room.NewRoomDelegate(mahjongContext.RoomManager, mahjongContext.DB))
		if flag := mahjongContext.RoomManager.AddRoom(r); !flag {
			log.Error("创建房间失败")
			return nil
		}

		nctx := mahjong.WithMahjong(r.Context(), mahjongContext)

		log.WithFields(
			log.Fields{
				"playerId":   playerId,
				"maxPlayers": maxPlayer,
				"round":      round,
				"roomId":     roomId,
				"roomConfig": roomConfig,
			}).Info("创建房间成功")
		//启动房间数据处理
		roomProcessor := mahjong.NewRoomProcessor(r, 10, time.Millisecond*20, mahjongContext.Dispatcher)
		roomProcessor.Start()
		nctx = mahjong.WithRoomProcessor(nctx, roomProcessor)
		r.SetContext(nctx)
	}

	//创建玩家
	pl = player.NewPlayer(playerId, roomId, s)
	tctx := player.WithPlayer(s.Context(), pl)
	s.SetContext(tctx)
	pl.Start()

	err = mahjongContext.PlayerManager.AddPlayer(pl)
	if err != nil {
		log.WithFields(
			log.Fields{
				"playerId": playerId,
			}).Warn("用户管理器添加用户失败")
		pl.Session().Close()
		return nil
	}
	//优化
	gcLogin := &loginpb.GCLogin{}
	gcLogin.PlayerId = &playerId
	gcMsg := &pb.Message{}

	gcMsgType := int32(messagetypepb.MessageType_GCLoginType)
	gcMsg.MessageType = &gcMsgType
	err = proto.SetExtension(gcMsg, loginpb.E_GcLogin, gcLogin)
	if err != nil {
		return err
	}
	gcMsgB, err := proto.Marshal(gcMsg)
	if err != nil {
		log.WithFields(
			log.Fields{
				"playerId": playerId,
				"error":    err.Error(),
			}).Error("发送登陆消息压缩异常")
		return err
	}

	pl.Send(gcMsgB)

	//判断是否在房间内
	rp := r.RoomPlayerManager().GetPlayerById(pl.Id())
	if rp != nil {
		log.WithFields(
			log.Fields{
				"playerId": playerId,
				"roomId":   roomId,
			}).Info("玩家重新连接")
		r.PlayerReconnect(pl)
		return nil
	} 
	ip := pl.Session().Ip() 
	rp = changsha.NewPlayer(playerId, ip, user.Name, user.Image, user.Sex , user.Location , pl)
	log.WithFields(
		log.Fields{
			"playerId": playerId,
			"roomId":   roomId,
		}).Info("玩家加入房间")
	flag := r.PlayerJoin(rp)
	if !flag {
		log.WithFields(
			log.Fields{
				"playerId": playerId,
				"roomId":   roomId,
			}).Error("玩家加入房间失败")
		return nil
	}
	log.WithFields(
		log.Fields{
			"playerId": playerId,
			"roomId":   roomId,
		}).Info("玩家加入房间成功")
	return nil
}
