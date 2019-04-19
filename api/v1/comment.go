/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package v1

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/service"
	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/captcha"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/e"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

// Comment 评论
type Comment struct {
	ID            uint   `validate:"required|number"`
	Email         string `validate:"required|email"`
	Nickname      string `validate:"required"`
	Content       string `validate:"required"`
	ArticleID     uint   `validate:"required|number"`
	ParentID      uint   `validate:"number"`
	ReplyEmail    string `validate:"email"`
	ReplyNickname string `validate:"minLen:1"`
	State         uint   `validate:"required|number"`
	PageNum       uint   `validate:"required|number"`
	CaptchaKey    string `validate:"required"`
	CaptchaValue  string `validate:"required"`
}

// ConfigValidation 验证配置
func (comment Comment) ConfigValidation(v *validate.Validation) {

	// 字段名称
	v.AddTranslates(validate.MS{
		"ID":            "评论ID",
		"Email":         "邮箱",
		"Nickname":      "用户昵称",
		"Title":         "文章标题",
		"Content":       "回复内容",
		"ArticleID":     "回复文章ID",
		"ParentID":      "回复评论ID",
		"ReplyEmail":    "回复评论邮箱",
		"ReplyNickname": "回复评论昵称",
		"State":         "评论状态",
		"PageNum":       "页数",
		"CaptchaKey":    "验证Key",
		"CaptchaValue":  "验证码",
	})

	// 错误信息
	v.AddMessages(validate.MS{
		"required":           "{field}不能为空",
		"number":             "{field}仅能为数字",
		"minLen":             "{field}最小长度%d",
		"email":              "{field}不是合法邮箱",
		"ID.required":        "评论ID错误",
		"ArticleID.required": "回复文章ID错误",
		"PageNum.required":   "页数错误",
	})

	// 场景
	v.WithScenes(validate.SValues{
		"CommentList": []string{
			"PageNum",
		},
		"CreateComment": []string{
			"ArticleID", "Nickname", "Email", "Content", "CaptchaKey", "CaptchaValue",
		},
		"UpdateComment": []string{
			"ID", "ArticleID", "Nickname", "Content", "State",
		},
		"DeleteComment": []string{
			"ID",
		},
	})
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// request
	var request *Comment
	if err := c.ShouldBind(&request); err != nil {
		logger.Warn(err)
		api.ErrValidate().Output(c)
		return
	}
	logger.Debug(request)

	// 验证
	v := validate.Struct(request, "CreateComment")

	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 验证验证码
	stat := captcha.VerifyCaptcha(request.CaptchaKey, request.CaptchaValue)
	if !stat {
		api.New(e.ErrCaptcha).Output(c)
		return
	}

	// 父级评论ID
	if request.ParentID > 0 {
		stat, err := service.IsExistsCommentByID(isAdmin(c), request.ID)
		if _, isErr := api.IsServiceError(c, err); isErr {
			return
		}
		if !stat {
			api.Err(e.ErrNotFoundParentComment).Output(c)
			return
		}
	}

	// 文章是否存在
	stat, err := service.IsExistsArticleByID(isAdmin(c), request.ArticleID)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if !stat {
		api.Err(e.ErrNotFoundArticle).Output(c)
		return
	}

	// 添加评论
	comment := &models.Comment{
		Email:         request.Email,
		Nickname:      request.Nickname,
		Content:       request.Content,
		ArticleID:     request.ArticleID,
		ParentID:      request.ParentID,
		ReplyEmail:    request.ReplyEmail,
		ReplyNickname: request.ReplyNickname,
	}
	err = service.CreateComment(comment)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("发表评论成功").Output(c)
	return
}

// UpdateComment 更新评论
func UpdateComment(c *gin.Context) {
	// request
	var request Comment
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "CreateComment")

	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 评论是否存在
	stat, err := service.IsExistsCommentByID(isAdmin(c), request.ID)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if !stat {
		api.Err(e.ErrNotFoundComment).Output(c)
		return
	}

	// 父级评论ID
	if request.ParentID > 0 {
		stat, err := service.IsExistsCommentByID(isAdmin(c), request.ID)
		if _, isErr := api.IsServiceError(c, err); isErr {
			return
		}
		if !stat {
			api.Err(e.ErrNotFoundParentComment).Output(c)
			return
		}
	}

	// 文章是否存在
	stat, err = service.IsExistsArticleByID(isAdmin(c), request.ArticleID)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if !stat {
		api.Err(e.ErrNotFoundArticle).Output(c)
		return
	}

	// 更新评论
	comment := &models.Comment{
		ID:            request.ID,
		Email:         request.Email,
		Nickname:      request.Nickname,
		Content:       request.Content,
		ArticleID:     request.ArticleID,
		ParentID:      request.ParentID,
		ReplyEmail:    request.ReplyEmail,
		ReplyNickname: request.ReplyNickname,
		State:         request.State,
	}
	err = service.UpdateComment(comment)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("更新评论成功").Output(c)
	return
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	// request
	var request Comment
	if err := c.ShouldBind(&request); err != nil {
		logger.Warn(err)
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "DeleteComment")

	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 删除评论
	err := service.DeleteComment(request.ID)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("删除评论成功").Output(c)
	return
}

// CommentList 评论列表
func CommentList(c *gin.Context) {
	// request
	var request *Comment
	var data = fmt.Sprintf(`{
		"PageNum":%v,
		"ArticleID":%v,
		"ParentID":%v
		}`,
		c.Param("PageNum"),
		c.DefaultQuery("article_id", "0"),
		c.DefaultQuery("parent_id", "0"),
	)
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		logger.Warn(err)
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "CommentList")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 过滤
	var whereKey string
	var whereVal []interface{}
	// 文章ID
	whereKey += " article_id = ? "
	whereVal = append(whereVal, request.ArticleID)

	// 父级评论ID
	// if request.ParentID > 0 {
	whereKey += "AND parent_id = ? "
	whereVal = append(whereVal, request.ParentID)
	// }

	where := append([]interface{}{whereKey}, whereVal...)
	count, err := service.CommentTotal(isAdmin(c), where)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	var commentList []*models.Comment
	pageSize := conf.Get().Page.Comment
	if count > 0 {
		commentList, err = service.GetCommentList(isAdmin(c), request.PageNum, pageSize, where)
		for key, item := range commentList {
			h := md5.New()
			h.Write([]byte(item.Email))
			commentList[key].Email = hex.EncodeToString(h.Sum(nil))
		}
	}

	response := map[string]interface{}{
		"count":       count,
		"commentList": &commentList,
		"pageSize":    pageSize,
		"pageNum":     request.PageNum,
	}
	api.Succ().SetData(&response).Output(c)
	return
}
