package model

type TaskModel struct {
	Id         int32  `gorm:"primary_key;column:id"`  
	Reward     int32  `gorm:"column:reward"`
	TargetNum  int32  `gorm:"column:targetNum"`
	Content    string `gorm:"column:content"`
	DeleteTime int64  `gorm:"column:deleteTime"`
	UpdateTime int64  `gorm:"column:updateTime"`
	CreateTime int64  `gorm:"column:createTime"`
}

func (usm *TaskModel) TableName() string {
	return "t_tasks"
}
