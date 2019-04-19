/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package models

// TagMap 标签文章映射
type TagMap struct {
	ArticleID uint `gorm:"primary_key,not null"`
	TagID     uint `gorm:"primary_key,not null"`
}

// Total 标签文章映射统计
func (tm *TagMap) Total(condition []interface{}) (count uint, err error) {
	err = db.Scopes(
		where(condition),
	).Model(&tm).Count(&count).Error
	return
}
