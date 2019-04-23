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

// LinkTotal .
func (s *Service) LinkTotal(where []interface{}) (count uint, err error) {
	var cacheLink models.Link
	key := cache.GetKey(append(where, "LinkTotal")...)

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
	count, err = cacheLink.Total(where)
	if isErrDB(err) {
		return 0, err
	}

	// 保存缓存
	if count > 0 {
		cache.Get().SetEx(key, conf.Get().Cache.Expired, count)
	}

	return count, nil
}

// isExistsLink 是否存在友链
func (s *Service) isExistsLink(where ...interface{}) (IsExists bool, err error) {
	var count uint
	count, err = s.LinkTotal(where)

	if count > 0 {
		IsExists = true
	}
	return
}

// IsExistsLinkByURI 。
func (s *Service) IsExistsLinkByURI(uri string) (IsExists bool, err error) {
	return s.isExistsLink("uri = ?", uri)
}

// IsExistsLinkByID 。
func (s *Service) IsExistsLinkByID(id uint) (IsExists bool, err error) {
	return s.isExistsLink("id = ?", id)
}

// CreateLink .
func (s *Service) CreateLink(link *models.Link) (err error) {
	return link.Create()
}

// UpdateLink 。
func (s *Service) UpdateLink(link *models.Link) (err error) {
	return link.Update()
}

// DeleteLink .
func (s *Service) DeleteLink(id uint) (err error) {
	var link models.Link

	link.ID = id
	err = link.Delete()
	if isErrDB(err) {
		return
	}

	return
}

// GetLinkList .
func (s *Service) GetLinkList(pageNum, pageSize uint) (linkList []*models.Link, err error) {
	key := cache.GetKey(`GetLinkList`, pageNum, pageSize)

	// 获取缓存
	if s.IsCache {
		data, stat, err := cache.GetCacheData(key)
		if err == nil && stat {
			// 数据解析
			jsonData, err := helper.Debase64(data)
			if err == nil {
				json.Unmarshal(jsonData, &linkList)
				return linkList, nil
			}
			logger.Warn("缓存数据有误,无法解析：", key, data)
		}
	}

	// 查询数据
	var cacheLink models.Link
	linkList, _, err = cacheLink.List(pageNum, pageSize)
	if isErrDB(err) {
		return nil, err
	}

	// 保存缓存
	dataString, err := json.Marshal(&linkList)
	if err != nil {
		return nil, err
	}
	cache.Get().SetEx(key, conf.Get().Cache.Expired, helper.Enbase64(dataString))

	return
}
