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
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

// Article 文章
type Article struct {
	ID          uint
	Title       string     `gorm:"not null"`
	Description string     `gorm:"not null" json:"Description,omitempty"`
	Markdown    string     `gorm:"type:text;not null" json:"Markdown"`
	Content     string     `gorm:"type:text;not null" json:"Content,omitempty"`
	SubjectID   uint       `gorm:"default:0;index:article2subject" json:"SubjectID,omitempty"`
	ViewNum     uint       `gorm:"default:0;not null"`
	CommentNum  uint       `gorm:"default:0;not null"`
	State       uint       `gorm:"default:0;not null;index:article_state" json:"State,omitempty"`
	TagID       string     `json:"-"` // 冗余字段
	CreatedAt   *time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt   *time.Time `json:"UpdatedAt,omitempty"`
	Tags        []*Tag     `gorm:"many2many:tag_map" `
	Subject     *Subject   `gorm:"ForeignKey:SubjectID"`
}

// Info 文章详情
func (art *Article) Info(isAdmin bool) error {
	return db.Scopes(isAdminArticle(isAdmin)).Where(art).Preload("Tags").Preload("Subject").Last(&art).Error
}

// HotList 热门文章
func (art *Article) HotList(isAdmin bool, pageNum, pageSize uint) (articleList []*Article, err error) {
	err = db.Order("view_num desc").Scopes(isAdminArticle(isAdmin)).Select("id,title").Limit(pageSize).Offset(pageNum*pageSize - pageSize).Find(&articleList).Error
	return
}

// List 文章列表
func (art *Article) List(isAdmin bool, pageNum, pageSize uint, fields []interface{}, order string, condition []interface{}) (articleList []*Article, err error) {
	err = db.Scopes(
		where(condition),
	).Scopes(
		field(fields),
	).Scopes(
		isAdminArticle(isAdmin),
	).Offset(
		pageNum*pageSize - pageSize,
	).Order(order).Limit(pageSize).Preload("Tags").Find(&articleList).Error
	return
}

// isAdminArticle 管理权限
func isAdminArticle(isAdmin bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !isAdmin {
			return db.Where("state = ?", 0)
		}
		return db
	}
}

// Create 添加文章
func (art *Article) Create(tags []string, subjectID uint) (err error) {
	tx := db.Begin()

	if err = tx.Create(&art).Error; err != nil {
		tx.Callback()
		return err
	}

	// 标签
	if err = art.updateTagMap(tx, tags); err != nil {
		tx.Callback()
		return err
	}

	// 专题
	if subjectID > 0 {
		if err = art.updateSubjectInfo(tx, subjectID); err != nil {
			tx.Callback()
			return err
		}
	}

	tx.Commit()

	return nil
}

// Update 更新文章
func (art *Article) Update(tags []string, subjectID uint) (err error) {
	tx := db.Begin()

	// 更新文章
	if err = tx.Save(&art).Error; err != nil {
		tx.Callback()
		return err
	}

	// 获取文章信息
	// if err = tx.First(&art, art.ID).Error; err != nil {
	// 	tx.Callback()
	// 	return err
	// }

	// 标签
	if err = art.updateTagMap(tx, tags); err != nil {
		tx.Callback()
		return err
	}

	// 专题
	if err = art.updateSubjectInfo(tx, subjectID); err != nil {
		tx.Callback()
		return err
	}

	tx.Commit()

	return nil
}

// Delete 删除文章
func (art *Article) Delete() (err error) {
	tx := db.Begin()

	// 获取文章信息
	if err = tx.First(&art, art.ID).Error; err != nil {
		tx.Callback()
		return err
	}

	// 标签
	if err = art.updateTagMap(tx, []string{}); err != nil {
		tx.Callback()
		return err
	}

	// 专题
	if err = art.updateSubjectInfo(tx, 0); err != nil {
		tx.Callback()
		return err
	}

	// 删除文章
	if err = tx.Delete(&art).Error; err != nil {
		tx.Callback()
		return err
	}

	tx.Commit()

	return nil

}

// updateTagMap 更新文章标签关联
func (art *Article) updateTagMap(tx *gorm.DB, tags []string) (err error) {
	// 删除关联标签映射
	err = tx.Where("article_id = ?", art.ID).Delete(TagMap{}).Error
	if err != nil {
		return err
	}

	// 添加标签与标签映射
	if len(tags) < 1 {
		return
	}
	var articleTag string
	for _, val := range tags {
		var tag = Tag{
			Name: val,
		}
		if err = tx.Where(&tag).FirstOrCreate(&tag).Error; err != nil {
			tx.Callback()
			return err
		}
		var tagMap = TagMap{
			ArticleID: art.ID,
			TagID:     tag.ID,
		}
		if err = tx.Where(&tagMap).FirstOrCreate(&tagMap).Error; err != nil {
			tx.Callback()
			return err
		}

		// 冗余字段
		articleTag += "[" + strconv.Itoa(int(tag.ID)) + "]"
	}

	// 更新冗余字段
	art.TagID = articleTag
	if err = tx.Save(&art).Error; err != nil {
		tx.Callback()
		return err
	}

	return
}

// updateSubjectInfo 更新文章专题信息
func (art *Article) updateSubjectInfo(tx *gorm.DB, subjectID uint) (err error) {
	if art.SubjectID == uint(subjectID) {
		return
	}

	// 去除原专题数据
	if art.SubjectID > 0 {
		var subject Subject
		err = tx.Model(&subject).Where("id = ?", art.SubjectID).Updates(
			map[string]interface{}{
				"article_num": gorm.Expr("article_num - ?", 1),
			},
		).Error
		if err != nil {
			tx.Callback()
			return err
		}
	}

	// 添加专题数据
	if subjectID > 0 {
		var subject Subject
		err = tx.Model(&subject).Where("id = ?", subjectID).Updates(
			map[string]interface{}{
				"article_num":             gorm.Expr("article_num + ?", 1),
				"article_last_updated_at": time.Now(),
			},
		).Error
		if err != nil {
			tx.Callback()
			return err
		}
	}

	return nil
}

// Total 文章统计
func (art *Article) Total(isAdmin bool, condition []interface{}) (count uint, err error) {
	err = db.Scopes(
		where(condition),
	).Scopes(
		isAdminArticle(isAdmin),
	).Model(&art).Count(&count).Error
	return
}
