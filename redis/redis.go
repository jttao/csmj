package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	OK = "OK"
)

//redis 配置
type RedisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	Wait        bool   `json:"wait"`
	IdleTimeout int64  `json:"idleTimeout"`
}

type RedisService interface {
	Pool() *redis.Pool
}

type redisService struct {
	config *RedisConfig
	pool   *redis.Pool
}

func (rs *redisService) Pool() *redis.Pool {
	return rs.pool
}

func NewRedisService(rc *RedisConfig) (RedisService, error) {
	pool := &redis.Pool{}
	dataSourceName := fmt.Sprintf("%s:%d", rc.Host, rc.Port)
	idleTimeout := rc.IdleTimeout * int64(time.Second)
	//connect
	pool = &redis.Pool{
		MaxIdle:     rc.MaxIdle,
		MaxActive:   rc.MaxActive,
		Wait:        rc.Wait,
		IdleTimeout: time.Duration(idleTimeout),

		Dial: func() (redis.Conn, error) {
			fmt.Println("redis connect " + dataSourceName)
			c, err := redis.Dial("tcp", dataSourceName)
			if err != nil {
				panic(err.Error())
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &redisService{
		config: rc,
		pool:   pool,
	}, nil
}
