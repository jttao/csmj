package model

type RoomRecordModel struct {
	Id         int64  `gorm:"primary_key,column:id"`
	RoomType   int64  `gorm:"column:roomType"`
	RoomId     int64  `gorm:"column:roomId"`
	OwnerId    int64  `gorm:"column:ownerId"`
	Player1    int64  `gorm:"column:player1"`
	Player2    int64  `gorm:"column:player2"`
	Player3    int64  `gorm:"column:player3"`
	Player4    int64  `gorm:"column:player4"`
	Settle     string `gorm:"column:settle"`
	UpdateTime int64  `gorm:"column:updateTime"`
	CreateTime int64  `gorm:"column:createTime"`
	DeleteTime int64  `gorm:"column:deleteTime"`
}

func (rr *RoomRecordModel) TableName() string {
	return "t_room_record"
}

type RoundModel struct {
	Id           int64  `gorm:"primary_key,column:id"`
	RoomRecordId int64  `gorm:"column:roomRecordId"`
	RoomType     int64  `gorm:"column:roomType"`
	RoomId       int64  `gorm:"column:roomId"`
	Round        int32  `gorm:"column:round"`
	TotoalRound  int32  `gorm:"column:totalRound"`
	Logs         string `gorm:"column:logs"`
	Settle       string `gorm:"column:settle"`
	Config       string `gorm:"column:config"`
	UpdateTime   int64  `gorm:"column:updateTime"`
	CreateTime   int64  `gorm:"column:createTime"`
	DeleteTime   int64  `gorm:"column:deleteTime"`
}

func (rr *RoundModel) TableName() string {
	return "t_round"
}
