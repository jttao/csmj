package model

type NoticeModel struct {
	Id         int64  `gorm:"primary_key,column:id"`
	Content    string `gorm:"column:content"`
	UpdateTime int64  `gorm:"column:updateTime"`
	CreateTime int64  `gorm:"column:createTime"`
	DeleteTime int64  `gorm:"column:deleteTime"`
}

func (nm *NoticeModel) TableName() string {
	return "t_notice"
}
