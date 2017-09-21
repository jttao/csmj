package model

type User struct {
	Id            int64  `gorm:"primary_key;column:id"`
	Weixin        string `gorm:"column:weixin"`
	DeviceMac     string `gorm:"column:deviceMac"`
	Name          string `gorm:"column:name"`
	Image         string `gorm:"column:image"`
	Sex           int    `gorm:"column:sex"`
	CardNum       int64  `gorm:"column:cardNum"`
	Forbid        int    `gorm:"column:forbid"`
	LastLoginIp   string `gorm:"column:lastLoginIp"`
	LastLoginTime int64  `gorm:"column:lastLoginTime"`
	State         int    `gorm:"column:state"`
	UpdateTime    int64  `gorm:"column:updateTime"`
	CreateTime    int64  `gorm:"column:createTime"`
	DeleteTime    int64  `gorm:"column:deleteTime"`
	Location   	  string `gorm:"column:location"`
}

func (u *User) TableName() string {
	return "t_user"
}
