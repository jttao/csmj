package service

import (
	"context"

	gamedb "game/db"
	"game/mahjong/record/model"

	"github.com/jinzhu/gorm"
)

type RecordService interface {
	GetRecordList(playerId int64) ([]*model.RoomRecordModel, error)
	GetRoundList(roomId int64) ([]*model.RoundModel, error)
	GetRound(roundId int64) (m *model.RoundModel, err error)
}

const (
	recordServiceKey = "RecordService"
)

func RecordServiceInContext(ctx context.Context) RecordService {
	c, ok := ctx.Value(recordServiceKey).(RecordService)
	if !ok {
		return nil
	}
	return c
}

func WithRecordService(ctx context.Context, rs RecordService) context.Context {
	return context.WithValue(ctx, recordServiceKey, rs)
}

type recordService struct {
	ds gamedb.DBService
}

func (rs *recordService) GetRecordList(playerId int64) (rl []*model.RoomRecordModel, err error) {
	tdb := rs.ds.DB().Order("createTime desc").Limit(10).Find(&rl, "player1=? or player2=? or player3=? or player4 =?", playerId, playerId, playerId, playerId)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rl, nil
}

func (rs *recordService) GetRoundList(roomId int64) (ml []*model.RoundModel, err error) {

	var rl model.RoomRecordModel
	tdb := rs.ds.DB().Order("createTime desc").First(&rl, "roomId=?", roomId)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return
	}
	startTime := rl.CreateTime
	tdb = rs.ds.DB().Order("round asc").Find(&ml, "roomId=? and createTime >=?", roomId, startTime)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return ml, nil
}

func (rs *recordService) GetRound(roundId int64) (m *model.RoundModel, err error) {
	m = &model.RoundModel{}
	tdb := rs.ds.DB().First(m, "id=?", roundId)
	err = tdb.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return m, nil
}

func NewRecordService(ds gamedb.DBService) RecordService {
	rs := &recordService{}
	rs.ds = ds
	return rs
}
