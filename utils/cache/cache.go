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

	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/helper"
	"github.com/wonderivan/logger"
)

type cacheType int

const (
	// Redis 缓存
	Redis cacheType = 1
)

var (
	cacheInterface *Cache

	// CacheType 缓存类型
	CacheType = Redis
)

// Cache 缓存接口
type Cache interface {
	Connect()
	Set(key string, value interface{}) error
	SetNx(key string, value interface{}) error
	SetEx(key string, seconds, value interface{}) error
	Get(key string) (string, error)
	Del(key string) (bool, error)
	Exists(key string) (bool, error)
	Expire(key string, seconds interface{}) error
	Persist(key string) error
}

// Get 获取缓存实例
func Get() Cache {
	logger.Debug("2222222")
	return *cacheInterface
}

// SetCacheType 设置缓存类型
func SetCacheType(cacheTypeValue cacheType) {
	CacheType = cacheTypeValue
}

// GetKey 获取缓存键值
func GetKey(args ...interface{}) string {
	config := conf.Get()
	data := []byte(config.Cache.Prefix + fmt.Sprint(args...))
	return helper.Enbase64(data)
}

// New 获取缓存实例
func New() Cache {
	var redisCache Cache = &redis{}
	cacheInterface = &redisCache
	logger.Debug("12121")
	return redisCache
}

// GetCacheData 获取缓存数据快捷函数
func GetCacheData(index interface{}) (string, bool, error) {
	key := fmt.Sprintf("%v", index)
	// 是否存在
	stat, err := Get().Exists(key)
	if err != nil {
		logger.Error(" GetCacheData 错误：", err)
		return "", false, err
	}
	if !stat {
		return "", false, nil
	}

	// 获取缓存
	value, err := Get().Get(key)
	if err != nil {
		logger.Warn("GetCacheData 错误：", key)
		return "", true, nil
	}
	logger.Debug("使用缓存")
	return value, true, nil
}
