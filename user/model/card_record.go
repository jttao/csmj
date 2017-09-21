package model

type CardRecordModel struct {
	Id         int64 `gorm:"primary_key;column:id"`
	UserId     int64 `gorm:"column:userId"`
	ChangeNum  int64 `gorm:"column:changeNum"`
	Reason     int   `gorm:"column:reason"`
	CreateTime int64 `gorm:"column:createTime"`
}

func (rm *CardRecordModel) TableName() string {
	return "t_card_record"
}

type ReasonType int

const (
	ReasonTypeGM ReasonType = iota
	ReasonTypeCost
	ReasonTypeRefund
)
