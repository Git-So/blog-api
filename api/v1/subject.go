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

	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/service"
	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/e"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"github.com/wonderivan/logger"
)

// Subject 接口提交数据
type Subject struct {
	ID      uint   `validate:"required|number"`
	Title   string `validate:"required"`
	State   uint   `validate:"number"`
	PageNum uint   `validate:"required|number"`
	Search  string
}

// ConfigValidation 验证配置
func (subject Subject) ConfigValidation(v *validate.Validation) {
	// 字段名称
	v.AddTranslates(validate.MS{
		"ID":    "专题ID",
		"Title": "专题标题",
		"State": "专题状态",
	})

	// 错误信息
	v.AddMessages(validate.MS{
		"required":         "{field}不能为空",
		"number":           "{field}仅能为数字",
		"maxLen":           "{field}最大长度%d",
		"ID.required":      "专题ID错误",
		"PageNum.required": "页数错误",
	})

	// 场景
	v.WithScenes(validate.SValues{
		"SubjectList": []string{
			"PageNum",
		},
		"SubjectInfo": []string{
			"ID",
		},
		"CreateSubject": []string{
			"Title",
		},
		"UpdateArticle": []string{
			"ID", "Title", "State",
		},
		"DeleteSubject": []string{
			"ID",
		},
	})
}

// CreateSubject 创建专题
func CreateSubject(c *gin.Context) {
	// request
	var request Subject
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "CreateSubject")

	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 专题是否存在
	stat, err := service.New(c).IsExistsSubjectByTitle(request.Title)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if stat {
		api.Err(e.ErrExistsSubject).Output(c)
		return
	}

	// 创建专题
	subject := &models.Subject{
		Title: request.Title,
	}
	err = service.New(c).CreateSubject(subject)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("专题创建成功").Output(c)
	return
}

// UpdateSubject 更新专题
func UpdateSubject(c *gin.Context) {
	// request
	var request *Subject
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

	// 专题是否存在
	if request.ID > 0 {
		stat, err := service.New(c).IsExistsSubjectByID(request.ID)
		if _, isErr := api.IsServiceError(c, err); isErr {
			return
		}
		if !stat {
			api.Err(e.ErrNotFoundSubject).Output(c)
			return
		}
	}

	// 更新专题
	subject := &models.Subject{
		ID:    request.ID,
		Title: request.Title,
		State: request.State,
	}
	err := service.New(c).UpdateSubject(subject)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("专题更新成功").Output(c)
	return
}

// DeleteSubject 删除专题
func DeleteSubject(c *gin.Context) {
	// request
	var request *Subject
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "DeleteSubject")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// do
	err := service.New(c).DeleteSubject(request.ID)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("专题删除成功").Output(c)
	return
}

// SubjectList 专题列表
func SubjectList(c *gin.Context) {
	// request
	var request *Subject
	var data = fmt.Sprintf(`{
		"PageNum":%v,
		"Search":"%v"
		}`,
		c.Param("PageNum"),
		c.DefaultQuery("search", ""),
	)
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "SubjectList")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 过滤
	var whereKey string
	var whereVal []interface{}
	if len(request.Search) > 0 { // 关键词
		whereKey += " title LIKE ? "
		whereVal = append(whereVal, "%"+request.Search+"%")
	}
	where := append([]interface{}{whereKey}, whereVal...)
	count, err := service.New(c).SubjectTotal(where)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	var subjectList []*models.Subject
	pageSize := conf.Get().Page.Subject
	if count > 0 {
		subjectList, err = service.New(c).GetSubjectList(request.PageNum, pageSize, where)
	}

	response := map[string]interface{}{
		"count":       count,
		"pageSize":    pageSize,
		"pageNum":     request.PageNum,
		"subjectList": &subjectList,
		"search":      request.Search,
	}
	api.Succ().SetData(&response).Output(c)
	return
}

// SubjectInfo 专题详情
func SubjectInfo(c *gin.Context) {
	// request
	data := fmt.Sprintf(`{"ID":%s}`, c.Param("ID"))
	var request Article
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		logger.Warn(err)
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "SubjectInfo")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// response
	subjectInfo, err := service.New(c).GetSubjectInfoByID(request.ID, isAdmin(c))
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
		"subjectInfo": subjectInfo,
	}
	api.Succ().SetData(&response).Output(c)
	return
}
