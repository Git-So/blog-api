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

// SubjectTotal .
func SubjectTotal(where []interface{}) (count uint, err error) {
	var cacheSubject models.Subject
	key := cache.GetKey(append(where, `SubjectTotal`)...)

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
	count, err = cacheSubject.Total(where)
	if isErrDB(err) {
		return 0, err
	}

	// 保存缓存
	if count > 0 {
		cache.Get().SetEx(key, conf.Get().Cache.Expired, count)
	}

	return count, nil
}

// isExistsSubject 是否存在专题
func isExistsSubject(where ...interface{}) (IsExists bool, err error) {
	var count uint
	count, err = SubjectTotal(where)

	if count > 0 {
		IsExists = true
	}
	return
}

// IsExistsSubjectByID 。
func IsExistsSubjectByID(id uint) (IsExists bool, err error) {
	return isExistsSubject("id = ?", id)
}

// IsExistsSubjectByTitle 。
func IsExistsSubjectByTitle(title string) (IsExists bool, err error) {
	return isExistsSubject("title = ?", title)
}

// CreateSubject .
func CreateSubject(subject *models.Subject) (err error) {
	return subject.Create()
}

// DeleteSubject 删除专题
func DeleteSubject(id uint) (err error) {
	if id < 1 {
		return
	}
	var subject models.Subject

	subject.ID = id
	err = subject.Delete()
	if isErrDB(err) {
		return
	}

	return
}

// UpdateSubject 。
func UpdateSubject(subject *models.Subject) (err error) {
	return subject.Update()
}

// GetSubjectList .
func GetSubjectList(pageNum, pageSize uint, where []interface{}) (subjectList []*models.Subject, err error) {
	key := cache.GetKey(append(where, `GetSubjectList`, pageNum, pageSize)...)

	// 获取缓存
	data, stat, err := cache.GetCacheData(key)
	if err == nil && stat {
		// 数据解析
		jsonData, err := helper.Debase64(data)
		if err == nil {
			json.Unmarshal(jsonData, &subjectList)
			return subjectList, nil
		}
		logger.Warn("缓存数据有误,无法解析：", key, data)
	}

	// 查询数据
	var cacheSubject models.Subject
	subjectList, _, err = cacheSubject.List(pageNum, pageSize, "id desc", where)
	if isErrDB(err) {
		return nil, err
	}

	// 保存缓存
	dataString, err := json.Marshal(&subjectList)
	if err != nil {
		return nil, err
	}
	cache.Get().SetEx(key, conf.Get().Cache.Expired, helper.Enbase64(dataString))

	return
}
