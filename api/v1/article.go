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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/utils/e"

	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/gookit/validate"

	"github.com/Git-So/blog-api/service"

	"github.com/gin-gonic/gin"
)

// Article 接口提交数据
type Article struct {
	ID          uint     `validate:"required|number"`
	Title       string   `validate:"required|maxLen:90"`
	Description string   `validate:"required"`
	Markdown    string   `validate:"required"`
	Content     string   `validate:"required"`
	SubjectID   uint     `validate:"number|min:1"`
	Tags        []string `validate:"-"`
	State       uint     `validate:"required|number"`
	PageNum     uint     `validate:"required|number"`
	Search      string   ``
	TagID       uint     ``
}

// ConfigValidation 验证配置
func (article Article) ConfigValidation(v *validate.Validation) {

	// 字段名称
	v.AddTranslates(validate.MS{
		"ID":          "文章ID",
		"Title":       "文章标题",
		"Description": "文章简介",
		"Markdown":    "Markdown原件",
		"Content":     "文章内容",
		"SubjectID":   "专题ID",
		"State":       "文章状态",
		"PageNum":     "页数",
	})

	// 错误信息
	v.AddMessages(validate.MS{
		"required":         "{field}不能为空",
		"number":           "{field}仅能为数字",
		"min":              "{field}最小为%d",
		"maxLen":           "{field}最大长度%d",
		"ID.required":      "文章ID错误",
		"PageNum.required": "页数错误",
	})

	// 场景
	v.WithScenes(validate.SValues{
		"ArticleList": []string{
			"PageNum",
		},
		"CreateArticle": []string{
			"Title", "Description", "Markdown", "Content", "SubjectID",
		},
		"UpdateArticle": []string{
			"ID", "Title", "Description", "Content", "SubjectID", "State",
		},
		"ArticleInfo": []string{
			"ID",
		},
		"DeleteArticle": []string{
			"ID",
		},
		"HotArticleList": []string{
			"PageNum",
		},
	})
}

// ArticleInfo 文章详情
func ArticleInfo(c *gin.Context) {
	// request
	data := fmt.Sprintf(`{"ID":%s}`, c.Param("ID"))
	var request Article
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		logger.Warn(err)
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "ArticleInfo")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// response
	articleInfo, err := service.New(c).GetArticleInfoByID(request.ID, isAdmin(c))
	isNotFound, isErr := api.IsServiceError(c, err)
	if isErr {
		return
	}
	if isNotFound {
		api.New(e.ErrNotFoundData).Output(c)
		return
	}

	// response
	response := map[string]interface{}{
		"articleInfo": articleInfo,
	}
	api.Succ().SetData(&response).Output(c)
	return
}

// ArticleList 文章列表
func ArticleList(c *gin.Context) {
	// request
	var request *Article
	var data = fmt.Sprintf(`{
		"PageNum":%v,
		"Subject":%v,
		"TagID":%v,
		"Search":"%v"
		}`,
		c.Param("PageNum"),
		c.DefaultQuery("subject_id", "0"),
		c.DefaultQuery("tag_id", "0"),
		c.DefaultQuery("search", ""),
	)
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "ArticleList")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 过滤
	var whereKey = "1 = 1 "
	var whereVal []interface{}
	if request.SubjectID > 0 { // 专题
		whereKey += "AND subject_id = ? "
		whereVal = append(whereVal, request.SubjectID)
	}
	if request.TagID > 0 { // 标签
		whereKey += "AND tag_id LIKE ? "
		whereVal = append(whereVal, "%["+strconv.Itoa(int(request.TagID))+"]%")
	}
	if len(request.Search) > 0 { // 关键词
		whereKey += "AND title LIKE ? "
		whereVal = append(whereVal, "%"+request.Search+"%")
	}
	where := append([]interface{}{whereKey}, whereVal...)
	count, err := service.New(c).ArticleTotal(isAdmin(c), where)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	var articleList []*models.Article
	pageSize := conf.Get().Page.Article
	if count > 0 {
		field := []interface{}{"id,title,description,created_at,view_num,comment_num"}
		order := "id desc"
		articleList, err = service.New(c).GetArticleList(isAdmin(c), request.PageNum, pageSize, field, order, where)
	}

	response := map[string]interface{}{
		"count":       count,
		"pageSize":    pageSize,
		"pageNum":     request.PageNum,
		"articleList": &articleList,
		"search":      request.Search,
	}
	api.Succ().SetData(&response).Output(c)
	return
}

// CreateArticle 创建文章
func CreateArticle(c *gin.Context) {
	// request
	var request Article
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "CreateArticle")

	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 文章是否存在
	stat, err := service.New(c).IsExistsArticleByTitle(isAdmin(c), request.Title)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if stat {
		api.Err(e.ErrExistsArticle).Output(c)
		return
	}

	// 专题是否存在
	logger.Debug(request.SubjectID)
	if request.SubjectID > 0 {
		stat, err = service.New(c).IsExistsSubjectByID(request.SubjectID)
		if _, isErr := api.IsServiceError(c, err); isErr {
			return
		}
		if !stat {
			api.Err(e.ErrNotFoundSubject).Output(c)
			return
		}
	}

	// 创建文章
	article := &models.Article{
		Title:       request.Title,
		Description: request.Description,
		Markdown:    request.Markdown,
		Content:     request.Content,
		SubjectID:   request.SubjectID,
	}
	logger.Debug(request.SubjectID)
	err = service.New(c).CreateArticle(article, request.Tags, request.SubjectID)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("文章创建成功").Output(c)
	return
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context) {
	// request
	var request *Article
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "UpdateArticle")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 文章是否存在
	stat, err := service.New(c).IsExistsArticleByID(isAdmin(c), request.ID)

	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if !stat {
		api.Err(e.ErrNotFoundArticle).Output(c)
		return
	}

	// 专题是否存在
	if request.SubjectID > 0 {
		stat, err = service.New(c).IsExistsSubjectByID(request.ID)
		if _, isErr := api.IsServiceError(c, err); isErr {
			return
		}
		if !stat {
			api.Err(e.ErrNotFoundSubject).Output(c)
			return
		}
	}

	// 更新文章
	article := &models.Article{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
		Content:     request.Content,
		SubjectID:   request.SubjectID,
		State:       request.State,
	}
	err = service.New(c).UpdateArticle(article, request.Tags, request.SubjectID)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("文章更新成功").Output(c)
	return
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	// request
	var request *Article
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "DeleteArticle")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// do
	err := service.New(c).DeleteArticle(request.ID)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("文章删除成功").Output(c)
	return
}

// HotArticleList 获取热门文章
func HotArticleList(c *gin.Context) {
	// request
	data := fmt.Sprintf(`{"PageNum":%s}`, c.Param("PageNum"))
	var request Article
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "HotArticleList")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// response
	articleList, err := service.New(c).GetHotArticleList(isAdmin(c), request.PageNum, conf.Get().Page.HotArticle)
	isNotFound, isErr := api.IsServiceError(c, err)
	if isErr {
		return
	}
	if isNotFound {
		api.New(e.ErrNotFoundData).Output(c)
		return
	}

	response := map[string]interface{}{
		"hotArticleList": &articleList,
	}
	api.Succ().SetData(&response).Output(c)
	return
}
