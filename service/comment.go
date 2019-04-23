/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 * 评论不做缓存
 */

package service

import (
	"github.com/Git-So/blog-api/models"
)

// CommentTotal .
func (s *Service) CommentTotal(isAdmin bool, where []interface{}) (count uint, err error) {
	var comment models.Comment

	// 查询数据
	count, err = comment.Total(isAdmin, where)
	if isErrDB(err) {
		return 0, err
	}

	return count, nil
}

// isExistsComment 是否存在评论
func (s *Service) isExistsComment(isAdmin bool, where ...interface{}) (IsExists bool, err error) {
	var count uint
	count, err = s.CommentTotal(isAdmin, where)

	if count > 0 {
		IsExists = true
	}
	return
}

// IsExistsCommentByID 。
func (s *Service) IsExistsCommentByID(isAdmin bool, id uint) (IsExists bool, err error) {
	return s.isExistsComment(isAdmin, "id = ?", id)
}

// CreateComment .
func (s *Service) CreateComment(comment *models.Comment) (err error) {
	return comment.Create()
}

// UpdateComment 。
func (s *Service) UpdateComment(comment *models.Comment) (err error) {
	return comment.Update()
}

// DeleteComment .
func (s *Service) DeleteComment(id uint) (err error) {
	var comment models.Comment

	comment.ID = id
	err = comment.Delete()
	if isErrDB(err) {
		return
	}

	return
}

// GetCommentList .
func (s *Service) GetCommentList(isAdmin bool, pageNum, pageSize uint, where []interface{}) (commentList []*models.Comment, err error) {

	// 查询数据
	var cacheComment models.Comment
	commentList, err = cacheComment.List(isAdmin, pageNum, pageSize, where)
	if isErrDB(err) {
		return nil, err
	}

	return
}
