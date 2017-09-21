package service

import (
	"context"
	gamedb "game/db"
	gameredis "game/redis"
	roommanagemodel "game/roommanage/model"
	usermodel "game/user/model"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

type RoomCardService interface {
	GetTodayUse() (int64, error)
	ChangeCardNum(userId int64, num int64) error
	Free(flag bool) error
	IfFree() (flag bool, err error)
	Check(flag bool) error
	IfCheck() (flag bool, err error)
}

const (
	getTodayUseSql = `
	select sum(changeNum) as sum from t_card_record where createTime>=? and reason!=0
	`
)

const (
	roomCardServiceKey = "gm.roomcard_service"
)

const (
	gameConfigFreeKey = "game.config.free"
)

const (
	gameConfigCheckKey = "game.config.check"
)

func WithRoomCardService(ctx context.Context, us RoomCardService) context.Context {
	return context.WithValue(ctx, roomCardServiceKey, us)
}

func RoomCardServiceInContext(ctx context.Context) RoomCardService {
	us, ok := ctx.Value(roomCardServiceKey).(RoomCardService)
	if !ok {
		return nil
	}
	return us
}

func NewRoomCardService(db gamedb.DBService, rs gameredis.RedisService) RoomCardService {
	us := &roomCardService{
		db: db,
		rs: rs,
	}
	return us
}

type roomCardService struct {
	db gamedb.DBService
	rs gameredis.RedisService
}

func (us *roomCardService) GetTodayUse() (num int64, err error) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	timeNumber := t.UnixNano() / int64(time.Millisecond)
	tempStruct := &SumModel{}
	tdb := us.db.DB().Raw(getTodayUseSql, timeNumber).Scan(tempStruct)
	err = tdb.Error
	if err != nil {
		return
	}
	return int64(tempStruct.Sum), nil
}

func (us *roomCardService) Check(flag bool) (err error) {

	pool := us.rs.Pool()
	conn := pool.Get()
	if conn.Err() != nil {
		return conn.Err()
	}
	defer conn.Close()

	_, err = conn.Do("set", gameConfigCheckKey, flag)
	if err != nil {
		return err
	}
	return nil
}

func (us *roomCardService) IfCheck() (flag bool, err error) {
	pool := us.rs.Pool()
	conn := pool.Get()
	if conn.Err() != nil {
		return false, conn.Err()
	}
	defer conn.Close()

	flag, err = redis.Bool(conn.Do("get", gameConfigCheckKey))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (us *roomCardService) Free(flag bool) (err error) {

	pool := us.rs.Pool()
	conn := pool.Get()
	if conn.Err() != nil {
		return conn.Err()
	}
	defer conn.Close()

	_, err = conn.Do("set", gameConfigFreeKey, flag)
	if err != nil {
		return err
	}
	return nil
}

func (us *roomCardService) IfFree() (flag bool, err error) {
	pool := us.rs.Pool()
	conn := pool.Get()
	if conn.Err() != nil {
		return false, conn.Err()
	}
	defer conn.Close()

	flag, err = redis.Bool(conn.Do("get", gameConfigFreeKey))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (us *roomCardService) ChangeCardNum(userId int64, num int64) error {
	//添加redis锁
	user := &usermodel.User{}
	tdb := us.db.DB().First(user, "id=?", userId)
	err := tdb.Error
	if err != nil {
		return err
	}

	newCardNum := user.CardNum + num
	tdb = us.db.DB().Model(user).Update("cardNum", newCardNum)
	err = tdb.Error
	if err != nil {
		return err
	}

	//添加日志
	cardRecordModel := &roommanagemodel.CardRecordModel{}
	cardRecordModel.ChangeNum = num
	cardRecordModel.UserId = userId
	cardRecordModel.Reason = int(roommanagemodel.ReasonTypeGM)
	cardRecordModel.CreateTime = time.Now().UnixNano() / int64(time.Millisecond)
	tdb = us.db.DB().Save(cardRecordModel)
	if tdb.Error != nil {
		log.WithFields(
			log.Fields{
				"userId": userId,
				"num":    num,
			}).Error("保存房卡纪录失败")
	}
	return nil
}

type SumModel struct {
	Sum int `gorm:"column:sum"`
}
