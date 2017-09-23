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

const (
	taskShareId 	= 1
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
	
	RewardUserTask(userId int64,taskId int32) error
	GetUserTasks(userId int64) (us []*UserTask, err error)
	
	refresh(us *UserTask, now int64) (flag bool,err error)
	
}

type taskService struct {
	db gamedb.DBService
	rs gameredis.RedisService 
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
	now := time.Now().UnixNano() / int64(time.Millisecond)
	uss = make([]*UserTask, 0, len(nms)) 
	for _, nm := range nms {
		
		n := &UserTask{}
		n.Id = nm.Id
		n.TaskId = nm.TaskId
		n.UserId = nm.UserId
		n.Reward = nm.Reward
		n.Finish = nm.Finish
		n.TargetNum = nm.TargetNum
		n.UpdateTime = nm.UpdateTime
		n.CreateTime = nm.CreateTime 
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
	
	us = &UserTask{}
	us.Id = utm.Id
	us.TaskId = utm.TaskId
	us.UserId = utm.UserId
	us.Reward = utm.Reward
	us.Finish = utm.Finish 
	us.TargetNum = utm.TargetNum  
	us.UpdateTime = us.UpdateTime
	us.CreateTime = us.CreateTime

	return us,nil
}

func (ts *taskService) RewardUserTask(userId int64,taskId int32) error{ 
	utm := &model.UserTaskModel{}
	tdb := ts.db.DB().First(utm,"taskId=? and userId=?",taskId, userId ) 
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	} 	
	//是否已经领取奖励
	if ut.Reward != 0 {
		return nil
	}
	//是否完成
	if ut.Finish < ut.TargetNum {
		return nil
	}   

	t , err := ts.GetTask(taskId) 
	if err != nil {
		return nil
	} 
	//设置奖励数据
	utm.Reward = t.Reward
	
	tdb = ts.db.DB().Model(utm).Update(utm) 
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	} 	
	return nil
}

//完成一次任务
func (ts *taskService) FinishUserTask(userId int64,taskId int32,state bool) (reward bool, us *UserTask, err error) {

	us , err = ts.GetUserTask(userId,taskId)

	if err != nil {
		return false,nil , err
	}  
	
	isnew := false

	now := time.Now().UnixNano() / int64(time.Millisecond)
	
	//没有记录,添加一条记录
	if us ==nil {
		t , err := ts.GetTask(taskId) 
		if err != nil {
			return  false,nil,err
		}   
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

		//是否跨天
		flag := ts.refresh(us, now)
		if flag {
			t , err := ts.GetTask(taskId) 
			if err != nil {
				return  false,nil,err
			} 
			us.TargetNum = t.TargetNum
		}

	}

	if us.Finish==us.TargetNum {
		return  false,us,nil
	}

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
	utm := convertUserTaskToModel(us)  
	
	//保存数据
	if isnew {
		tdb = ts.db.DB().Save(utm) 
	}else{
		tdb = ts.db.DB().Model(utm).Update(utm)  
	}   
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return  false,nil,nil
		}
		return  false,nil,tdb.Error
	}  
	reward = false  
	if us.Finish==us.TargetNum {
		reward = true 
	} 	
	return reward,us,nil
}

func (ts *taskService) refresh(us *UserTask, now int64) (flag bool,err error) {
	//判断是否跨天
	flag, err :=  pkgtimutils.IsSameDay(now, us.UpdateTime)
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
