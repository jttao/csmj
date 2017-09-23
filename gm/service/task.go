package service

import (
	"context"
	gamedb "game/db"
	"game/pkg/timeutils"
	model "game/hall/model"
	"math"

	"github.com/jinzhu/gorm"
)

type TaskService interface { 
	GetTasks(page int, pageSize int) (int, []*model.TaskModel, error)
	GetTaskById(id int32) (u *model.TaskModel, err error)
	UpdateTask(id int32,reward int32,targetNum int32,content string) error
} 

const (
	taskServiceKey = "gm.task_service"
)

func WithTaskService(ctx context.Context, ts TaskService) context.Context {
	return context.WithValue(ctx, taskServiceKey, ts)
}

func TaskServiceInContext(ctx context.Context) TaskService {
	ts, ok := ctx.Value(taskServiceKey).(TaskService)
	if !ok {
		return nil
	}
	return ts
}

func NewTaskService(db gamedb.DBService) TaskService {
	ts := &TaskService{
		db: db,
	}
	return ts
}

type taskService struct {
	db gamedb.DBService
}

func (ts *taskService) GetTasks(page int, pageSize int) (totalPage int, tasks []*model.TaskModel, err error) {
	offset := (page - 1) * pageSize  
	tdb := ts.db.DB().Order("id desc").Offset(offset).Limit(pageSize).Find(&tasks)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil, nil
		}
		return
	}	
	count := 0
	tdb = ts.db.DB().Model(&model.TaskModel{}).Count(&count)
	err = tdb.Error
	if err != nil {
		return
	}
	totalPage = int(math.Ceil(float64(count) / float64(pageSize)))
	return totalPage, tasks, nil
}

func (ts *taskService) GetTaskById(id int32) (u *model.TaskModel, err error) {
	u = &model.TaskModel{}
	tdb := ts.db.DB().First(u, "id=?", id)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (ts *taskService) UpdateTask(id int32,reward int32,targetNum int32,content string) error {
	task, err := ts.GetTaskById(id) 
	if err != nil {
		return err
	}
	task.Id = id
	task.Reward = reward
	task.TargetNum = targetNum
	task.Content = content 
	tdb := ts.db.DB().Model(task).Update(task) 
	err := tdb.Error
	if err != nil {
		return err
	}
	return nil 
}

type CountModel struct {
	Count int `gorm:"column:count"`
}
