package client

import ( 
	"context" 

	usermodel "game/user/model"
	userservice "game/user/service"
	
	taskservice "game/hall/tasks"

	//gamepkghttputils "game/pkg/httputils"  
	//"game/hall/api" 
	log "github.com/Sirupsen/logrus"

)

type HallClient interface { 
	GameEnd(playerId int64, flag bool) error
	FinishTask(playerId int64,taskId int32,state bool) error
}

type hallClient struct {
	UserService      userservice.UserService  
	TaskService      taskservice.TaskService 
}

func (hc *hallClient) GameEnd(playerId int64, maxwin bool) error {
	
	log.WithFields(
		log.Fields{
			"playerId": playerId, 
			"maxwin":  maxwin,  
		}).Debug("开始设置玩家任务> ")
	
	//每日游戏任务   
	err := hc.FinishTask(playerId,2,true)
	if err != nil {
		return err
	}  
	if maxwin { 
		//每日连续赢任务 
		err1 := hc.FinishTask(playerId,3,true)
		if err1 != nil {
			return err1
		} 
	}   
	return nil
}

func (hc *hallClient) FinishTask(playerId int64,taskId int32,state bool) error{
	
	log.WithFields(
		log.Fields{
			"userId": playerId,
		}).Debug("请求每日任务完成")
	
	reward,ut,err := hc.TaskService.FinishUserTask(playerId,taskId,state)
	
	if err != nil {  
		log.WithFields(
			log.Fields{
				"userId": playerId,
				"error":  err,
			}).Error("请求每日任务完成,错误")
		return nil
	} 
	
	if reward { 
		reason := usermodel.ReasonType(0)
		if taskId==1 {
			reason = usermodel.ReasonTypeTask1
		}
		if taskId==2 {
			reason = usermodel.ReasonTypeTask2
		}
		if taskId==3 {
			reason = usermodel.ReasonTypeTask3
		} 
		err := hc.UserService.ChangeCardNum( playerId , int64(ut.Reward), reason ) 	
		if err != nil { 
			log.WithFields(log.Fields{
				"userId": 	playerId,
				"Reward":   ut.Reward,
				"error":   	err,
			}).Error("请求每日任务完成,发送奖励失败")
			return nil
		}
	}
	return nil
}

func NewHallClient(ts taskservice.TaskService ,us userservice.UserService,) HallClient {
	rmc := &hallClient{}
	rmc.TaskService = ts 
	rmc.UserService = us
	return rmc
}

const (
	key = "hall_client"
)

func WithHallClient(ctx context.Context, rmc HallClient) context.Context {
	return context.WithValue(ctx, key, rmc)
}

func HallClientInContext(ctx context.Context) HallClient {
	us, ok := ctx.Value(key).(HallClient)
	if !ok {
		return nil
	}
	return us
}
