/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package models

import (
	"github.com/wonderivan/logger"
)

// Link 友情链接
type Link struct {
	ID        uint
	URI       string
	AvatarURI string
	Nickname  string
	Title     string
}

// Info 友链详情
func (lk *Link) Info() (err error) {
	return db.Where(&lk).Last(&lk).Error
}

// Create 创建友链
func (lk *Link) Create() (err error) {
	logger.Debug("sos")
	return db.Create(&lk).Error
}

// Update 更新友链
func (lk *Link) Update() (err error) {
	return db.Save(&lk).Error
}

// Delete 删除友链
func (lk *Link) Delete() (err error) {
	return db.Delete(&lk).Error
}

// List 友链列表
func (lk *Link) List(pageNum, pageSize uint) (linkLink []*Link, count int, err error) {
	err = db.Model(&lk).Count(&count).Limit(pageSize).Offset(pageNum*pageSize - pageSize).Find(&linkLink).Error
	return
}

// Total 友链统计
func (lk *Link) Total(condition []interface{}) (count uint, err error) {
	err = db.Scopes(
		where(condition),
	).Model(&lk).Count(&count).Error
	return
}
