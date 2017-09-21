package model

type TipModel struct {
	Id         int64  `gorm:"primary_key,column:id"`
	Content    string `gorm:"column:content"`
	UpdateTime int64  `gorm:"column:updateTime"`
	CreateTime int64  `gorm:"column:createTime"`
	DeleteTime int64  `gorm:"column:deleteTime"`
}

func (nm *TipModel) TableName() string {
	return "t_tip"
}
