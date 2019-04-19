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

	"github.com/gomodule/redigo/redis"
	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/utils/conf"
)

var (
	redisConn *redis.Conn
	redisConf *conf.Cache
)

type redisCache struct {
	Cache
}

// Connect redis连接
func (cache *redisCache) Connect() {
	config := conf.Get()
	redisConf = config.Cache

	conn, err := redis.Dial(
		"tcp",
		fmt.Sprintf(
			"%s:%d",
			redisConf.Host,
			redisConf.Port,
		),
	)
	if err != nil {
		logger.Error("Redis 连接失败,缓存功能未能打开: ", err)
	} else {
		logger.Info("Redis 已连接")
	}

	redisConn = &conn
	return
}

// isRedisErr redis错误判断
func (cache *redisCache) isRedisErr(err error) bool {
	if err != nil {
		logger.Error("执行redis出错：", err)
		return true
	}
	return false
}

// Set 缓存设置
func (cache *redisCache) Set(key, value interface{}) (err error) {
	// 设置缓存
	_, err = (*redisConn).Do("SET", key, value)
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（SetKey）：", key, err)
		return
	}
	return
}

// SetNx 缓存设置
func (cache *redisCache) SetNx(key, value interface{}) (err error) {
	// 设置缓存
	_, err = (*redisConn).Do("SETNX", key, value)
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（SetKey）：", key, err)
		return
	}
	return
}

// SetEx 缓存设置
func (cache *redisCache) SetEx(key, seconds, value interface{}) (err error) {
	// 设置缓存
	_, err = (*redisConn).Do("SETEX", key, seconds, value)
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（SetKey）：", key, err)
		return
	}
	return
}

// Get 获取缓存
func (cache *redisCache) Get(key interface{}) (value string, err error) {
	value, err = redis.String((*redisConn).Do("GET", key))
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（GET）", key, err)
		return
	}
	return
}

// Del 删除缓存
func (cache *redisCache) Del(key interface{}) (stat bool, err error) {
	stat, err = redis.Bool((*redisConn).Do("DEL", key))
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（Del）", key, err)
		return false, err
	}
	return stat, nil
}

// Exists 存在缓存
func (cache *redisCache) Exists(key interface{}) (stat bool, err error) {
	stat, err = redis.Bool((*redisConn).Do("EXISTS", key))
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（EXISTS）：", key, err)
		return false, err
	}
	return stat, nil
}

// Expire 设置过期时间
func (cache *redisCache) Expire(key, expire interface{}) (err error) {
	_, err = (*redisConn).Do("EXPIRE", expire)
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（Expire）：", key, err)
		return
	}
	return
}

// Persist 取消过期时间
func (cache *redisCache) Persist(key interface{}) (err error) {
	_, err = (*redisConn).Do("PERSIST", key)
	if cache.isRedisErr(err) {
		logger.Error("执行 redis 错误（PERSIST）：", key, err)
		return
	}
	return
}
