package service

import (
	"context"
	gamedb "game/db"
	"game/pkg/timeutils"
	usermodel "game/user/model"
	"math"

	"github.com/jinzhu/gorm"
)

type UserService interface {
	GetTodayNewUsers() (int, error)
	GetOnlinesUsers() (int, error)
	GetTotalUsers() (int, error)
	ForbidUser(id int64) error
	UnforbidUser(id int64) error
	GetUserById(id int64) (*usermodel.User, error)
	GetUserByWxId(wxId string) (*usermodel.User, error)
	GetUsers(page int, pageSize int) (int, []*usermodel.User, error)
}

const (
	getTodayNewUsersSql = `
	select count(*) as count from t_user where createTime>=?
	`
)

const (
	getOnlineUsersSql = `
	select count(*) as count from t_user where state=1
	`
)

const (
	getTotalUserSql = `
	select count(*) as count from t_user
	`
)

const (
	updateForbidSql = `
	UPDATE t_user SET forbid=? where id=?
	`
)

const (
	userServiceKey = "gm.user_service"
)

func WithUserService(ctx context.Context, us UserService) context.Context {
	return context.WithValue(ctx, userServiceKey, us)
}

func UserServiceInContext(ctx context.Context) UserService {
	us, ok := ctx.Value(userServiceKey).(UserService)
	if !ok {
		return nil
	}
	return us
}

func NewUserService(db gamedb.DBService) UserService {
	us := &userService{
		db: db,
	}
	return us
}

type userService struct {
	db gamedb.DBService
}

func (us *userService) GetTodayNewUsers() (num int, err error) {

	timeNumber, err := timeutils.BeginOfNow()
	if err != nil {
		return
	}
	tempStruct := &CountModel{}
	tdb := us.db.DB().Raw(getTodayNewUsersSql, timeNumber).Scan(tempStruct)
	err = tdb.Error
	if err != nil {
		return
	}
	return tempStruct.Count, nil
}

func (us *userService) GetOnlinesUsers() (num int, err error) {

	tempStruct := &CountModel{}
	tdb := us.db.DB().Raw(getOnlineUsersSql).Scan(tempStruct)
	err = tdb.Error
	if err != nil {
		return
	}
	return tempStruct.Count, nil
}

func (us *userService) GetUsers(page int, pageSize int) (totalPage int, users []*usermodel.User, err error) {
	offset := (page - 1) * pageSize

	tdb := us.db.DB().Order("state desc,cardNum desc,createTime desc").Offset(offset).Limit(pageSize).Find(&users)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil, nil
		}
		return
	}
	count := 0
	tdb = us.db.DB().Model(&usermodel.User{}).Count(&count)
	err = tdb.Error
	if err != nil {
		return
	}
	totalPage = int(math.Ceil(float64(count) / float64(pageSize)))
	return totalPage, users, nil
}

func (us *userService) GetTotalUsers() (num int, err error) {
	tempStruct := &CountModel{}
	tdb := us.db.DB().Raw(getTotalUserSql).Scan(tempStruct)
	err = tdb.Error
	if err != nil {
		return
	}
	return tempStruct.Count, nil
}

func (us *userService) ForbidUser(id int64) error {
	user := &usermodel.User{}
	user.Id = id
	tdb := us.db.DB().Model(user).Update("forbid", 1)
	err := tdb.Error
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) UnforbidUser(id int64) error {
	user := &usermodel.User{}
	user.Id = id
	tdb := us.db.DB().Model(user).Update("forbid", 0)
	err := tdb.Error
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) GetUserById(id int64) (u *usermodel.User, err error) {
	u = &usermodel.User{}
	tdb := us.db.DB().First(u, "id=?", id)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (us *userService) GetUserByWxId(wxId string) (u *usermodel.User, err error) {
	u = &usermodel.User{}
	tdb := us.db.DB().First(u, "weixin=?", wxId)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

type CountModel struct {
	Count int `gorm:"column:count"`
}
