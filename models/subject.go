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
	"time"
)

// Subject 专题
type Subject struct {
	ID                   uint
	Title                string
	State                uint
	ArticleNum           uint
	ArticleLastUpdatedAt time.Time
}

// Info 专题详情
func (sub *Subject) Info(isAdmin bool) (err error) {
	return db.Where(&sub).Last(&sub).Error
}

// Create 创建专题
func (sub *Subject) Create() (err error) {
	return db.Create(&sub).Error
}

// Update 更新专题
func (sub *Subject) Update() (err error) {
	return db.Save(&sub).Error
}

// Delete 删除专题
func (sub *Subject) Delete() (err error) {
	return db.Delete(&sub).Error
}

// List 专题列表
func (sub *Subject) List(pageNum, pageSize uint, order string, condition []interface{}) (sujectList []*Subject, count int, err error) {
	err = db.Scopes(
		where(condition),
	).Model(&sub).Count(&count).Limit(pageSize).Offset(pageNum*pageSize - pageSize).Order("id desc").Find(&sujectList).Error
	return
}

// // isAdminSubject 权限管理
// func isAdminSubject(isAdmin bool) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if !isAdmin {
// 			return db.Where("state = ?", 0)
// 		}
// 		return db
// 	}
// }

// Total 专题统计
func (sub *Subject) Total(condition []interface{}) (count uint, err error) {
	err = db.Scopes(
		where(condition),
	).Model(&sub).Count(&count).Error
	return
}
