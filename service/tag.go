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
	"encoding/json"
	"strconv"

	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/utils/cache"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/helper"
	"github.com/wonderivan/logger"
)

// TagTotal .
func (s *Service) TagTotal(where []interface{}) (count uint, err error) {
	var cacheTagInfo models.Tag
	key := cache.GetKey(append(where, `TagTotal`)...)

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
	count, err = cacheTagInfo.Total(where)
	if isErrDB(err) {
		return 0, err
	}

	// 保存缓存
	if count > 0 {
		cache.Get().SetEx(key, conf.Get().Cache.Expired, count)
	}

	return count, nil
}

// isExistsTag 是否存在标签
func (s *Service) isExistsTag(where []interface{}) (IsExists bool, err error) {
	var count uint
	count, err = s.TagTotal(where)

	if count > 0 {
		IsExists = true
	}
	return
}

// IsExistsTagByName .
func (s *Service) IsExistsTagByName(name string) (IsExists bool, err error) {
	return s.isExistsTag([]interface{}{"name = ?", name})
}

// CreateTag .
func (s *Service) CreateTag(tag *models.Tag) (err error) {
	return tag.Create()
}

// DeleteTag 删除标签
func (s *Service) DeleteTag(id uint) (err error) {
	var tag models.Tag

	// 删除标签
	tag.ID = id
	err = tag.Delete()
	if isErrDB(err) {
		return
	}

	return
}

// DeleteTagByName 使用标签名删除标签
func (s *Service) DeleteTagByName(name string) (err error) {
	var tag models.Tag

	// 删除标签
	tag.Name = name
	err = tag.Delete()
	if isErrDB(err) {
		return
	}

	return
}

// GetTagList .
func (s *Service) GetTagList(pageNum, pageSize uint, where []interface{}) (tagList []*models.Tag, err error) {
	key := cache.GetKey(append(where, `GetTagList`, pageNum, pageSize)...)

	// 获取缓存
	if s.IsCache {
		data, stat, err := cache.GetCacheData(key)
		if err == nil && stat {
			// 数据解析
			jsonData, err := helper.Debase64(data)
			if err == nil {
				json.Unmarshal(jsonData, &tagList)
				return tagList, nil
			}
			logger.Warn("缓存数据有误,无法解析：", key, data)
		}
	}

	// 查询数据
	var cacheTag models.Tag
	tagList, err = cacheTag.List(pageNum, pageSize, where)
	if isErrDB(err) {
		return nil, err
	}

	// 保存缓存
	dataString, err := json.Marshal(&tagList)
	if err != nil {
		return nil, err
	}
	cache.Get().SetEx(key, conf.Get().Cache.Expired, helper.Enbase64(dataString))

	return
}
