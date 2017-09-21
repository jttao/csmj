package redis

import (
	"math/rand"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	defaulRetryTimes = 10   //重试次数
	defaultTimeout   = 1000 //毫秒
	sleepMaxInterval = 100  //毫秒
	sleepMinInterval = 10   //毫秒
)
const (
	joinkeySeparate    = "."
	combineKeySeparate = ":"
)

func Join(keys ...string) string {
	return strings.Join(keys, joinkeySeparate)
}

func Combine(keys ...string) string {
	return strings.Join(keys, combineKeySeparate)
}

const (
	lockKey = "lock"
)

//redis 分布式锁 阻塞

func LockDefault(conn redis.Conn, key string) (bool, error) {
	return Lock(conn, key, defaultTimeout, defaulRetryTimes)
}

//加锁 毫毛为单位
func Lock(conn redis.Conn, key string, timeout int64, retryTimes int) (bool, error) {
	//不允许不设置过期时间，否则将会死锁
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	//拼接lock key
	lockKeyStr := Join(lockKey, key)
	//循环拿锁
	for i := 0; i < retryTimes; i++ {

		now := time.Now().UnixNano() / int64(time.Millisecond)
		timeout = now + timeout + 1
		val, err := redis.Int(conn.Do("setnx", lockKeyStr, timeout))

		//redis错误
		if err != nil {
			return false, err
		}

		//获得锁了
		if val == 1 {
			return true, nil
		}

		//获取过期时间戳
		getValue, err := redis.Int64(conn.Do("get", lockKeyStr))
		//redis 错误
		if err != nil {
			//重试
			if err == redis.ErrNil {
				continue
			}
			return false, err
		}

		//过期了
		if getValue < now {

			//原子操作取值设置新值
			getOldValue, err := redis.Int64(conn.Do("getset", lockKeyStr, now))
			//redis 错误
			if err != nil {
				//重试
				if err == redis.ErrNil {
					continue
				}

				return false, err
			}

			// 拿到锁了
			if getOldValue == getValue {
				return true, nil
			}
			//被人抢先了
		}

		//睡眠
		rand.Seed(time.Now().UnixNano())
		sleepInterval := rand.Intn(sleepMaxInterval-sleepMinInterval) + sleepMinInterval
		time.Sleep(time.Duration(sleepInterval) * time.Millisecond)
	}
	return false, nil
}

//解锁
func Unlock(conn redis.Conn, key string) (bool, error) {
	lockKeyStr := Join(lockKey, key)

	now := time.Now().Unix() / int64(time.Millisecond)
	getValue, err := redis.Int64(conn.Do("get", lockKeyStr))
	if err != nil {
		return false, err
	}

	//已经不是自己的锁了
	if now > getValue {
		return true, nil
	}

	//删除key
	_, err = redis.Int(conn.Do("del", lockKeyStr))
	if err != nil {
		return false, err
	}

	return true, nil
}
