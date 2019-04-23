/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package service

import (
	"strconv"

	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/utils/cache"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/wonderivan/logger"
)

// TagMapTotal .
func (s *Service) TagMapTotal(where []interface{}) (count uint, err error) {
	var cacheTagMapInfo models.TagMap
	key := cache.GetKey(append(where, `TagMapTotal`)...)

	// 获取缓存
	if s.IsCache {
		data, stat, err := cache.GetCacheData(key)
		if err == nil && stat {
			// 数据解析
			count, err := strconv.Atoi(data)
			if err == nil {
				return uint(count), nil
			}
			logger.Warn("缓存数据有误,无法解析：", key, data)
		}
	}

	// 查询数据
	count, err = cacheTagMapInfo.Total(where)
	if isErrDB(err) {
		return 0, err
	}

	// 保存缓存
	if count > 0 {
		cache.Get().SetEx(key, conf.Get().Cache.Expired, count)
	}

	return count, nil
}

// isExistsTag 是否存在标签文章映射
func (s *Service) isExistsTagMap(where []interface{}) (IsExists bool, err error) {
	var count uint
	count, err = s.TagMapTotal(where)

	if count > 0 {
		IsExists = true
	}
	return
}

// IsExistsTagByTagName .
func (s *Service) IsExistsTagByTagName(name string) (IsExists bool, err error) {
	tag := models.Tag{
		Name: name,
	}
	err = tag.Info()
	if isErrDB(err) {
		return
	}
	return s.isExistsTagMap([]interface{}{"tag_id = ?", tag.ID})
}
