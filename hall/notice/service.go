package news

import (
	"context"

	gamedb "game/db"
	"game/hall/model"
	gameredis "game/redis"

	"github.com/jinzhu/gorm"
)

type Notice struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
}

type NoticeService interface {
	GetNotices() ([]*Notice, error)
}

type noticeService struct {
	db gamedb.DBService
	rs gameredis.RedisService
}

func (ns *noticeService) GetNotices() (nl []*Notice, err error) {
	nms := make([]*model.NoticeModel, 0, 4)
	tdb := ns.db.DB().Find(&nms)
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tdb.Error
	}

	nl = make([]*Notice, 0, len(nms))
	for _, nm := range nms {
		n := &Notice{}
		n.Content = nm.Content
		n.Id = nm.Id
		nl = append(nl, n)
	}
	return
}

func NewNoticeService(db gamedb.DBService, rs gameredis.RedisService) NoticeService {
	ns := &noticeService{}
	ns.db = db
	ns.rs = rs
	return ns
}

const (
	noticeServiceKey = "NoticeService"
)

func NoticeServiceInContext(ctx context.Context) NoticeService {
	c, ok := ctx.Value(noticeServiceKey).(NoticeService)
	if !ok {
		return nil
	}
	return c
}

func WithNoticeService(ctx context.Context, ns NoticeService) context.Context {
	return context.WithValue(ctx, noticeServiceKey, ns)
}
