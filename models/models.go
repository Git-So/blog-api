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
	"log"

	"github.com/Git-So/blog-api/utils/conf"
	"github.com/jinzhu/gorm"

	// _ "github.com/jinzhu/gorm/dialects/mysql" // mysql
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite3
)

var db *gorm.DB

// Get 获取DB实例
func Get() *gorm.DB {
	if db == nil {
		new()
	}
	return db
}

func new() {
	config := conf.Get()

	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("DB err: %v", err)
	}

	db.LogMode(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	setTableName(true, config.Database.TablePrefix)

	db.AutoMigrate(&Article{}, &Comment{}, &Config{}, &Link{}, &Subject{}, &TagMap{}, &Tag{})

	return
}

// 表名处理
func setTableName(isSingular bool, prefix string) {
	db := Get()

	// 单数表名
	db.SingularTable(isSingular)

	// 表名规则
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return prefix + defaultTableName
	}
}

// where 条件过滤
func where(condition []interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(condition) < 1 {
			return db
		}
		return db.Where(condition[0], condition[1:]...)
	}
}

// field 查询字段
func field(fields []interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(fields) > 1 {
			db.Select(fields[0], fields[1:]...)
		}
		return db.Select(fields[0])
	}
}
