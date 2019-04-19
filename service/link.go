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
func LinkTotal(where []interface{}) (count uint, err error) {
	var cacheLink models.Link
	key := cache.GetKey(append(where, "LinkTotal")...)

	// 获取缓存
	data, stat, err := cache.GetCacheData(key)
	if err == nil && stat {
		// 数据解析
		count, err := strconv.Atoi(data)
		if err == nil {
			return uint(count), nil
		}
		logger.Warn("缓存数据有误,无法解析：", key, data)
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
func isExistsLink(where ...interface{}) (IsExists bool, err error) {
	var count uint
	count, err = LinkTotal(where)

	if count > 0 {
		IsExists = true
	}
	return
}

// IsExistsLinkByURI 。
func IsExistsLinkByURI(uri string) (IsExists bool, err error) {
	return isExistsLink("uri = ?", uri)
}

// IsExistsLinkByID 。
func IsExistsLinkByID(id uint) (IsExists bool, err error) {
	return isExistsLink("id = ?", id)
}

// CreateLink .
func CreateLink(link *models.Link) (err error) {
	return link.Create()
}

// UpdateLink 。
func UpdateLink(link *models.Link) (err error) {
	return link.Update()
}

// DeleteLink .
func DeleteLink(id uint) (err error) {
	var link models.Link

	link.ID = id
	err = link.Delete()
	if isErrDB(err) {
		return
	}

	return
}

// GetLinkList .
func GetLinkList(pageNum, pageSize uint) (linkList []*models.Link, err error) {
	key := cache.GetKey(`GetLinkList`, pageNum, pageSize)

	// 获取缓存
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
