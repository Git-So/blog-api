/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package models

import "github.com/jinzhu/gorm"

// Config 配置
type Config struct {
	Key   string `gorm:"primary_key,not null"`
	Value string `gorm:"not null"`
	State uint
}

// Info 获取配置信息详情
func (cf *Config) Info(isAdmin bool) (err error) {
	if isAdmin {
		cf.State = 1
	}
	return db.Scopes(isAdminConfig(isAdmin)).Where(&cf).FirstOrCreate(&cf).Error
}

// Create 创建配置
func (cf *Config) Create() (err error) {
	return db.Create(&cf).Error
}

// Update 更新配置
func (cf *Config) Update() (err error) {
	return db.Save(&cf).Error
}

// UpdateAll 批量更新
func (cf *Config) UpdateAll(cfs []*Config) (err error) {
	tx := db.Begin()
	for _, val := range cfs {
		if err := tx.Where("key = ? ", val.Key).FirstOrCreate(&val).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&val).Where("key = ? ", val.Key).Updates(&val).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// isAdminConfig 权限管理
func isAdminConfig(isAdmin bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !isAdmin {
			return db.Where("state = ?", 0)
		}
		return db
	}
}

// List 配置列表
func (cf *Config) List(isAdmin bool, condition []interface{}) (configList []*Config, err error) {
	err = db.Scopes(
		where(condition),
	).Scopes(isAdminConfig(isAdmin)).Find(&configList).Error
	return
}
