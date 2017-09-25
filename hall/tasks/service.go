package tasks

import (
	"context"
	
	gamedb "game/db"
	"game/hall/model"
	gameredis "game/redis"
	
	"github.com/jinzhu/gorm"
	
	pkgtimutils "game/pkg/timeutils" 
	
	"time"
)

type contextKey string

const (
	taskServiceKey = contextKey("TaskService")
)

type TaskId int

const (
	taskShareId TaskId = 1
	taskWinId 		= 2
	taskPlayId 		= 3
)

type Task struct {
	Id      	int32  `json:"id"` 
	Reward      int32  `json:"reward"`
	TargetNum  	int32  `json:"targetNum"`
	Content 	string `json:"content"`
}

type UserTask struct {
	Id         int64 `json:"id"`
	TaskId     int32 `json:"taskId"`
	UserId     int64 `json:"userId"`
	Reward     int32 `json:"reward"`
	Finish     int32 `json:"finish"`
	TargetNum  int32 `json:"targetNum"` 
	UpdateTime int64 `json:"updateTime"`
	CreateTime int64 `json:"createTime"` 
}

type TaskService interface {  
	
	GetTask(taskId int32) (t *Task, err error) 
	GetUserTask(userId int64,taskId int32) (us *UserTask, err error) 
	FinishUserTask(userId int64,taskId int32,state bool) (reward bool,us *UserTask, err error) 
	
	RewardUserTask(userId int64,taskId int32)(reward bool,us *UserTask,err error)
	GetUserTasks(userId int64) (us []*UserTask, err error)
	
	refresh(us *UserTask, now int64) (flag bool,err error)
	
}

type taskService struct {
	db gamedb.DBService
	rs gameredis.RedisService 
}

func (ts *taskService)  GetTask(taskId int32) (t *Task, err error){
	utm := &model.TaskModel{}
	tdb := ts.db.DB().First(utm,"id=?",taskId)
	
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil,nil
		}
		return nil,tdb.Error
	}   
	t = &Task{}
	t.Id = utm.Id 
	t.Reward = utm.Reward 
	t.TargetNum = utm.TargetNum   
	t.Content = utm.Content 
	return t,nil
}  
		
func (ts *taskService) GetUserTasks(userId int64) (uss []*UserTask, err error) {

	nms := make([]*model.UserTaskModel, 0, 4)
	tdb := ts.db.DB().Find(&nms,"userId=?",userId)
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tdb.Error
	}    
	
	uss = make([]*UserTask, 0, len(nms)) 
	for _, nm := range nms {
		n := convertUserTaskFromModel(nm) 
		uss = append(uss, n)  
	}  
	return uss,nil
}

func (ts *taskService) GetUserTask(userId int64,taskId int32) (us *UserTask, err error)  {
	
	utm := &model.UserTaskModel{}
	tdb := ts.db.DB().First(utm,"userId=? and taskId=?",userId,taskId)
	
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil,nil
		}
		return nil,tdb.Error
	}
	
	us = convertUserTaskFromModel(utm)

	return us,nil
}

func (ts *taskService) RewardUserTask(userId int64,taskId int32) (reward bool,us *UserTask,err error) { 
	
	us , err = ts.GetUserTask(userId,taskId)
	
	if err != nil {
		return false,nil , err
	}  
	
	t , err := ts.GetTask(taskId) 
	if err != nil {
		return false,nil,nil
	} 

	now := time.Now().UnixNano() / int64(time.Millisecond)
	
	//是否跨天
	flag , _ := ts.refresh(us, now) 
	if flag { 
		us.TargetNum = t.TargetNum 
	}

	//是否已经领取奖励
	if us.Reward != 0 {
		return false,us,nil
	}
	//是否完成
	if us.Finish < us.TargetNum {
		return false,us,nil
	}   
	
	//设置奖励数据
	us.Reward = t.Reward
	
	utm := convertUserTaskToModel(us)  
	
	tdb := ts.db.DB().Model(utm).Update(utm) 
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return false,us,nil
		}
		return false,us,tdb.Error
	} 	
	
	return true,us,nil
}

//完成一次任务
func (ts *taskService) FinishUserTask(userId int64,taskId int32,state bool) (reward bool, us *UserTask, err error) {

	us , err = ts.GetUserTask(userId,taskId)

	if err != nil {
		return false,nil , err
	}  
	
	isnew := false

	now := time.Now().UnixNano() / int64(time.Millisecond)
	
	t , err := ts.GetTask(taskId) 
	if err != nil {
		return  false,nil,err
	}   
	//没有记录,添加一条记录
	if us ==nil {
		
		us = &UserTask{} 
		us.TaskId = taskId
		us.UserId = userId
		us.Reward = 0
		us.Finish = 0
		us.TargetNum = t.TargetNum
		us.UpdateTime = now
		us.CreateTime = now 
		isnew = true 
	}else{ 
		//是否跨天，刷新任务
		flag , _ := ts.refresh(us, now) 
		if flag { 
			us.TargetNum = t.TargetNum
		}  
	}

	//已经完成，领取奖励了
	if us.Reward != 0 { 
		return  false,us,nil
	}

	//添加完成记录
	count := us.Finish 
	if state {
		count += 1
		if count>us.TargetNum {
			count = us.TargetNum
		}
	}else{
		count = 0
	}   
	us.Finish = count 
	
	reward = false  
	
	//设置完成，获取奖励
	if us.Finish==us.TargetNum {
		us.Reward = t.Reward
		reward = true 
	} 	
	
	utm := convertUserTaskToModel(us)  
	
	//保存数据
	if isnew {  
		tdb := ts.db.DB().Save(utm) 
		if tdb.Error != nil {
			if tdb.Error == gorm.ErrRecordNotFound {
				return  false,nil,nil
			}
			return  false,nil,tdb.Error
		}  
	}else{ 
		tdb := ts.db.DB().Model(utm).Update(utm)  
		if tdb.Error != nil {
			if tdb.Error == gorm.ErrRecordNotFound {
				return  false,nil,nil
			}
			return  false,nil,tdb.Error
		}  
	}   
	
	return reward,us,nil
}

func (ts *taskService) refresh(us *UserTask, now int64) (flag bool,err error) {
	//判断是否跨天
	flag, err =  pkgtimutils.IsSameDay(now, us.UpdateTime)
	if err != nil {
		return false,err
	}
	if flag {
		return false,err
	} 
	us.Finish = 0 
	us.Reward = 0  
	us.UpdateTime = now 
	return true,err
} 	

func TaskServiceInContext(ctx context.Context) TaskService {
	c, ok := ctx.Value(taskServiceKey).(TaskService)
	if !ok {
		return nil
	}
	return c
}

func WithTaskService(ctx context.Context, ss TaskService) context.Context {
	return context.WithValue(ctx, taskServiceKey, ss)
} 

func NewTaskService(db gamedb.DBService, rs gameredis.RedisService) TaskService {
	ns := &taskService{}
	ns.db = db
	ns.rs = rs
	return ns
}

const (
	newsServiceKey = "NewsService"
)
 

func convertUserTaskFromModel(utm *model.UserTaskModel) (us *UserTask) {
	us = &UserTask{
		Id:         utm.Id,
		UserId:     utm.UserId,
		TaskId:     utm.TaskId,
		Reward:  	utm.Reward,
		Finish:     utm.Finish,
		TargetNum:  utm.TargetNum,
		UpdateTime: utm.UpdateTime,
		CreateTime: utm.CreateTime,
	}	
	return us
}

func convertUserTaskToModel(mm *UserTask) *model.UserTaskModel {
	m := &model.UserTaskModel{
		Id:         mm.Id,
		UserId:     mm.UserId,
		TaskId:     mm.TaskId,
		Reward:  	mm.Reward,
		Finish:     mm.Finish,
		TargetNum:  mm.TargetNum,
		UpdateTime: mm.UpdateTime,
		CreateTime: mm.CreateTime,
	}
	return m
}
