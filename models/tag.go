/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package models

// Tag 标签
type Tag struct {
	ID   uint
	Name string `gorm:"type:varchar(30);not null;unique_index:tag_name"`
}

// Info 标签详情
func (tag *Tag) Info() (err error) {
	return db.Where(&tag).Last(&tag).Error
}

// Create 新建标签
func (tag *Tag) Create() (err error) {
	return db.Create(&tag).Error
}

// Delete 删除标签
func (tag *Tag) Delete() (err error) {
	return db.Where(&tag).Delete(&tag).Error
}

// List 标签分页
func (tag *Tag) List(pageNum, pageSize uint, condition []interface{}) (tags []*Tag, err error) {
	err = db.Scopes(
		where(condition),
	).Limit(pageSize).Offset(pageNum*pageSize - pageSize).Find(&tags).Error
	return
}

// Total 标签统计
func (tag *Tag) Total(condition []interface{}) (count uint, err error) {
	err = db.Scopes(
		where(condition),
	).Model(&tag).Count(&count).Error
	return
}
