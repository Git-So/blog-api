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
)

// Tag 标签接口参数
type Tag struct {
	ID      uint   `validate:"required|number"`
	Name    string `validate:"required|maxLen:20"`
	PageNum uint   `validate:"required|number"`
	Search  string
}

// ConfigValidation 验证配置
func (tag Tag) ConfigValidation(v *validate.Validation) {
	// 字段名称
	v.AddTranslates(validate.MS{
		"ID":      "标签ID",
		"Name":    "标签名称",
		"PageNum": "页数",
	})

	// 错误信息
	v.AddMessages(validate.MS{
		"required":         "{field}不能为空",
		"number":           "{field}仅能为数字",
		"maxLen":           "{field}最大长度%d",
		"ID.required":      "标签ID错误",
		"PageNum.required": "页数错误",
	})

	// 场景
	v.WithScenes(validate.SValues{
		"CreateTag": []string{
			"Name",
		},
		"TagList": []string{
			"PageNum",
		},
		"DeleteTag": []string{
			"Name",
		},
	})
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	// request
	var request Tag
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "CreateTag")

	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 标签是否存在
	stat, err := service.New(c).IsExistsTagByName(request.Name)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if stat {
		api.Err(e.ErrExistsTag).Output(c)
		return
	}

	// 创建标签
	tag := &models.Tag{
		Name: request.Name,
	}
	err = service.New(c).CreateTag(tag)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("标签创建成功").Output(c)
	return
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	// request
	var request *Tag
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "DeleteTag")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 标签是否已使用
	stat, err := service.New(c).IsExistsTagByTagName(request.Name)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if stat {
		api.Err(e.ErrUsedTag).Output(c)
		return
	}

	// do
	err = service.New(c).DeleteTagByName(request.Name)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("标签删除成功").Output(c)
	return
}

// TagList 标签列表
func TagList(c *gin.Context) {
	// request
	var request *Tag
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
	v := validate.Struct(request, "TagList")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 过滤
	var whereKey string
	var whereVal []interface{}
	if len(request.Search) > 0 { // 关键词
		whereKey += " name LIKE ? "
		whereVal = append(whereVal, "%"+request.Search+"%")
	}
	where := append([]interface{}{whereKey}, whereVal...)
	count, err := service.New(c).TagTotal(where)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	var tagList []*models.Tag
	pageSize := conf.Get().Page.Tag
	if count > 0 {
		tagList, err = service.New(c).GetTagList(request.PageNum, pageSize, where)
	}

	response := map[string]interface{}{
		"count":    count,
		"pageSize": pageSize,
		"pageNum":  request.PageNum,
		"tagList":  &tagList,
		"search":   request.Search,
	}
	api.Succ().SetData(&response).Output(c)
	return
}
