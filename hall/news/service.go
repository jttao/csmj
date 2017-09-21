package news

import (
	"context"

	gamedb "game/db"
	"game/hall/model"
	gameredis "game/redis"

	"github.com/jinzhu/gorm"
)

type NewsService interface {
	GetNews() (string, error)
	SetNews(news string) error
}

type newsService struct {
	db gamedb.DBService
	rs gameredis.RedisService
}

func (ns *newsService) GetNews() (string, error) {
	news := &model.NewsModel{}
	tdb := ns.db.DB().First(news)
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", tdb.Error
	}

	return news.Content, nil
}

func (ns *newsService) SetNews(newsStr string) error {
	news := &model.NewsModel{}
	tdb := ns.db.DB().First(news)
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	}

	news.Content = newsStr
	tdb = ns.db.DB().Model(news).Update("content", newsStr)
	if tdb.Error != nil {
		if tdb.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return tdb.Error
	}
	return nil
}

func NewNewsService(db gamedb.DBService, rs gameredis.RedisService) NewsService {
	ns := &newsService{}
	ns.db = db
	ns.rs = rs
	return ns
}

const (
	newsServiceKey = "NewsService"
)

func NewsServiceInContext(ctx context.Context) NewsService {
	c, ok := ctx.Value(newsServiceKey).(NewsService)
	if !ok {
		return nil
	}
	return c
}

func WithNewsService(ctx context.Context, ns NewsService) context.Context {
	return context.WithValue(ctx, newsServiceKey, ns)
}
