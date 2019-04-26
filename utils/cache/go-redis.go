/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package cache

import (
	"fmt"
	"time"

	"github.com/Git-So/blog-api/utils/conf"
	goredis "github.com/go-redis/redis"
	"github.com/wonderivan/logger"
)

// redis struct
type redis struct {
}

var redisConn *goredis.Client

func (r *redis) Connect() {
	cacheConf := conf.Get().Cache

	conn := goredis.NewClient(
		&goredis.Options{
			Addr:     fmt.Sprintf("%v:%v", cacheConf.Host, cacheConf.Port),
			Password: "",
			DB:       0,
		},
	)

	_, err := conn.Ping().Result()
	if err != nil {
		logger.Error("redis 连接失败:", err)
	}

	redisConn = conn

	return
}

// Set 缓存设置
func (r *redis) Set(key string, value interface{}) (err error) {
	logger.Debug("33333333")
	err = redisConn.Set(key, value, 0).Err()
	if err != nil {
		logger.Debug("55555555555")
		logger.Error("redis Set:", err)
		return
	}
	logger.Debug("444444444444")
	return
}

// SetNx 缓存不存在设置
func (r *redis) SetNx(key string, value interface{}) (err error) {
	err = redisConn.SetNX(key, value, 0).Err()
	if err != nil {
		logger.Error("redis SetNx:", err)
		return
	}
	return
}

// SetEx 缓存设置并且设置过期时间
func (r *redis) SetEx(key string, seconds, value interface{}) (err error) {
	expire := seconds.(int)
	err = redisConn.Set(key, value, time.Duration(expire)*time.Second).Err()
	if err != nil {
		logger.Error("redis SetEx:", err)
		return
	}
	return
}

// Get 获取缓存
func (r *redis) Get(key string) (result string, err error) {
	result, err = redisConn.Get(key).Result()
	if err != nil {
		logger.Error("redis Get:", err)
		return
	}
	return
}

// Del 删除缓存
func (r *redis) Del(key string) (stat bool, err error) {
	err = redisConn.Del(key).Err()
	if err != nil {
		logger.Error("redis Del:", err)
		return
	}
	return
}

// Exists 是否存在
func (r *redis) Exists(key string) (stat bool, err error) {
	var result int64
	result, err = redisConn.Exists(key).Result()
	logger.Debug(result)
	if err != nil {
		logger.Error("redis Exists:", err)
		return
	}
	return
}

// Expire 设置过期时间
func (r *redis) Expire(key string, seconds interface{}) (err error) {
	expire := seconds.(int)
	err = redisConn.Expire(key, time.Duration(expire)*time.Second).Err()
	if err != nil {
		logger.Error("redis Expire:", err)
		return
	}
	return
}

// Persist 取消过期时间
func (r *redis) Persist(key string) (err error) {
	err = redisConn.Persist(key).Err()
	if err != nil {
		logger.Error("redis Persist:", err)
		return
	}
	return
}
