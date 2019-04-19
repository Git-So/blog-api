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

	"github.com/wonderivan/logger"

	"github.com/jinzhu/gorm"
)

// Comment 评论
type Comment struct {
	ID            uint
	Email         string `gorm:"not null"`
	Nickname      string `gorm:"not null"`
	Content       string `gorm:"not null"`
	ArticleID     uint   `gorm:"not null"`
	ParentID      uint   `gorm:"default:0;not null"`
	ReplyEmail    string `gorm:"not null"`
	ReplyNickname string `gorm:"not null"`
	IsSub         bool   // 冗余字段
	State         uint   `gorm:"default:0;not null;index:comment_state" json:"State,omitempty"`
	CreatedAt     time.Time
}

// Info 获取评论详情
func (cm *Comment) Info(isAdmin bool) (err error) {
	return db.Where(&cm).Scopes(isAdminComment(isAdmin)).Last(&cm).Error
}

// Create 创建评论
func (cm *Comment) Create() (err error) {
	tx := db.Begin()
	if err := tx.Create(&cm).Error; err != nil {
		tx.Rollback()
		return err
	}
	logger.Debug(cm.ParentID)
	if cm.ParentID > 0 {
		err := tx.Model(&Comment{}).Where("id = ?", cm.ParentID).Update("is_sub", true).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// Update 更新评论
func (cm *Comment) Update() (err error) {
	return db.Save(&cm).Error
}

// Delete 删除评论
func (cm *Comment) Delete() (err error) {
	tx := db.Begin()
	// 删除子评论
	if err := tx.Where("parent_id = ?", cm.ParentID).Delete(Comment{}).Error; err != nil {
		tx.Callback()
		return err
	}

	if err := tx.Delete(&cm).Error; err != nil {
		tx.Callback()
		return err
	}
	tx.Commit()
	return nil
}

// List 评论列表
func (cm *Comment) List(isAdmin bool, pageNum, pageSize uint, condition []interface{}) (commentList []*Comment, err error) {
	err = db.Scopes(
		where(condition),
	).Scopes(isAdminComment(isAdmin)).Order("id desc").Limit(pageSize).Offset(pageNum*pageSize - pageSize).Find(&commentList).Error
	return
}

// Total 评论数量
func (cm *Comment) Total(isAdmin bool, condition []interface{}) (count uint, err error) {
	err = db.Scopes(
		where(condition),
	).Model(&cm).Scopes(isAdminComment(isAdmin)).Count(&count).Error
	return
}

// isAdminComment 管理权限
func isAdminComment(isAdmin bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !isAdmin {
			return db.Where("state = ?", 0)
		}
		return db
	}
}
