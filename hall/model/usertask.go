package model

type UserTaskModel struct {
	Id         int64 `gorm:"primary_key;column:id"`
	TaskId     int32 `gorm:"column:taskId"`
	UserId     int64 `gorm:"column:userId"`
	Reward     int32 `gorm:"column:reward"`
	Finish     int32 `gorm:"column:finish"`
	TargetNum  int32 `gorm:"column:targetNum"`
	DeleteTime int64 `gorm:"column:deleteTime"`
	UpdateTime int64 `gorm:"column:updateTime"`
	CreateTime int64 `gorm:"column:createTime"`
}

func (usm *UserTaskModel) TableName() string {
	return "t_user_task"
}
