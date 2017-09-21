package model

type RoomModel struct {
	Id         int64  `gorm:"primary_key,column:id"`
	RoomType   int    `gorm:"column:roomType"`
	RoomConfig string `gorm:"column:roomConfig"`
	Round      int    `gorm:"column:round"`
	Cost       int    `gorm:"column:cost"`
	OwnerId    int64  `gorm:"column:ownerId"`
	UpdateTime int64  `gorm:"column:updateTime"`
	CreateTime int64  `gorm:"column:createTime"`
	DeleteTime int64  `gorm:"column:deleteTime"`
}

func (rm *RoomModel) TableName() string {
	return "t_room"
}
