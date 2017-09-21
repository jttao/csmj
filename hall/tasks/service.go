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
	GetTasks() (tks []*Task,err error) 
	SetTask(t *Task)  error
	AddTask(id int32,reward int32,targetNum int32,content string) (t *Task,err error)
	
	SetUserTask(ut *UserTask) error
	GetUserTasks(userId int64) (us []*UserTask, err error)
	GetUserTask(userId int64,taskId int32) (us *UserTask, err error)
	GetReward(userId int64,taskId int32) (us *UserTask, err error) 

	refresh(us *UserTask, now int64) (err error)

	FinishUserTask(userId int64,taskId int32) (us *UserTask, err error) 

}

type taskService struct {
	db gamedb.DBService
	rs gameredis.RedisService 
}

func (ts *taskService) GetTasks() (tks []*Task,err error) {
	nms := make([]*model.TaskModel, 0, 4)
	tdb := ts.db.DB().Find(&nms)
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tdb.Error
	} 

	tks = make([]*Task, 0, len(nms))
	for _, nm := range nms {
		n := &Task{} 
		n.Content = nm.Content
		n.Reward = nm.Reward
		n.TargetNum = nm.TargetNum
		n.Id = nm.Id
		tks = append(tks, n)
	} 	
	return tks,nil
}

func (ts * taskService) SetTask(t *Task) (err error){

	tm := &model.TaskModel{}
	tdb := ts.db.DB().First(tm,"id=?",t.Id)
	
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	}
	
	tm.Content = t.Content
	tm.Reward = t.Reward
	tm.TargetNum = t.TargetNum  

	tdb = ts.db.DB().Model(tm).Update(tm)

	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	}
	return nil
} 

func (ts * taskService) AddTask(id int32,reward int32,targetNum int32,content string) (tm *model.TaskModel,err error) {
	tm = &model.TaskModel{} 
	tm.Id = id
	tm.Reward = reward
	tm.TargetNum = targetNum
	tm.Content = content
	
	tm.CreateTime = time.Now().UnixNano() / int64(time.Millisecond)
	tm.UpdateTime = tm.CreateTime 
	tdb := ts.db.DB().Save(tm)

	if tdb.Error != nil {
		return nil, tdb.Error
	}

	return
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
		n := &UserTask{}
		n.Id = nm.Id
		n.TaskId = nm.TaskId
		n.UserId = nm.UserId
		n.Reward = nm.Reward
		n.Finish = nm.Finish
		n.TargetNum = nm.TargetNum
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

	return us,nil
}

func (ts *taskService) SetUserTask(ut *UserTask) error{ 
	utm := &model.UserTaskModel{}
	tdb := ts.db.DB().First(utm,"taskId=? and userId=?",ut.TaskId,ut.UserId) 
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	} 	
	utm.Finish = ut.Finish
	utm.Reward = ut.Reward
	utm.TargetNum = ut.TargetNum   
	now := time.Now().UnixNano() / int64(time.Millisecond)
	err := ts.refresh(ut, now) 
	if err != nil {
		return err
	} 
	if ut.Finish != 0 {
		return nil
	} 
	tdb = ts.db.DB().Model(utm).Update(utm) 
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	} 	
	return nil
}

func (ts *taskService) GetReward(userId int64,taskId int32) (us *UserTask, err error) {
	
	us , err = ts.GetUserTask(userId,taskId)

	if err != nil { 
		return nil,err
	}
	
	if us.Finish < us.TargetNum {
		return nil,nil
	}

	if us.Finish==us.TargetNum {
		return nil,nil
	}
	
	return us,nil

}

func (ts *taskService) FinishUserTask(userId int64,taskId int32) (us *UserTask, err error) {
	
	us , err = ts.GetUserTask(userId,taskId)

	if err != nil {
		return nil , err
	}  
	
	if us.Finish==us.TargetNum {
		return us,nil
	}

	us.Finish = us.TargetNum 
	
	err = ts.SetUserTask(us)
	
	if err != nil {
		return nil,err
	}  
	
	return us,nil
}

func (ts *taskService) refresh(us *UserTask, now int64) (err error) {
	//判断是否跨天
	flag, err :=  pkgtimutils.IsSameDay(now, us.UpdateTime)
	if err != nil {
		return
	}
	if flag {
		return
	}
	us.Finish = 0 
	us.UpdateTime = now 
	return
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
