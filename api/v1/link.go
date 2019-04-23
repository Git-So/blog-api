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

// Link 接口提交数据
type Link struct {
	ID        uint   `validate:"required|number"`
	Title     string `validate:"maxLen:90"`
	URI       string `validate:"required|url"`
	Nickname  string `validate:"required|maxLen:90"`
	AvatarURI string
	State     int  `validate:"required|number"`
	PageNum   uint `validate:"required|number"`
}

// ConfigValidation 验证配置
func (link Link) ConfigValidation(v *validate.Validation) {
	// 字段名称
	v.AddTranslates(validate.MS{
		"ID":        "友链ID",
		"Title":     "友链标题",
		"URI":       "友链地址",
		"Nickname":  "站长昵称",
		"AvatarURI": "站长头像",
		"PageNum":   "页数",
	})

	// 错误信息
	v.AddMessages(validate.MS{
		"required":         "{field}不能为空",
		"number":           "{field}仅能为数字",
		"maxLen":           "{field}最大长度%d",
		"url":              "{field}不是合法网址",
		"ID.required":      "友链ID错误",
		"PageNum.required": "页数错误",
	})

	// 场景
	v.WithScenes(validate.SValues{
		"LinkList": []string{
			"PageNum",
		},
		"CreateLink": []string{
			"Title", "URI", "Nickname", "AvatarURI",
		},
		"UpdateLink": []string{
			"ID", "Title", "URI", "Nickname", "AvatarURI",
		},
		"DeleteLink": []string{
			"ID",
		},
	})
}

// CreateLink 创建友链
func CreateLink(c *gin.Context) {
	// request
	var request *Link
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "CreateLink")

	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 友链是否存在
	stat, err := service.New(c).IsExistsLinkByURI(request.URI)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if stat {
		api.Err(e.ErrExistsLink).Output(c)
		return
	}

	// 创建友链
	link := &models.Link{
		Title:     request.Title,
		URI:       request.URI,
		Nickname:  request.Nickname,
		AvatarURI: request.AvatarURI,
	}
	err = service.New(c).CreateLink(link)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("友链创建成功").Output(c)
	return
}

// UpdateLink 更新友链
func UpdateLink(c *gin.Context) {
	// request
	var request *Link
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "UpdateLink")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 友链是否存在
	stat, err := service.New(c).IsExistsLinkByID(request.ID)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	if !stat {
		api.Err(e.ErrNotFoundLink).Output(c)
		return
	}

	// 更新友链
	link := &models.Link{
		ID:        request.ID,
		Title:     request.Title,
		URI:       request.URI,
		Nickname:  request.Nickname,
		AvatarURI: request.AvatarURI,
	}
	err = service.New(c).UpdateLink(link)
	if isErr, _ := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("友链更新成功").Output(c)
	return
}

// DeleteLink 删除友链
func DeleteLink(c *gin.Context) {
	// request
	var request *Link
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "DeleteLink")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// do
	err := service.New(c).DeleteLink(request.ID)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}

	api.Succ().SetMsg("友链删除成功").Output(c)
	return
}

// LinkList 友链列表
func LinkList(c *gin.Context) {
	// request
	var request *Link
	var data = fmt.Sprintf(`{"PageNum":%v}`, c.Param("PageNum"))
	if err := json.Unmarshal([]byte(data), &request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "LinkList")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 过滤
	count, err := service.New(c).LinkTotal([]interface{}{})
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	var linkList []*models.Link
	pageSize := conf.Get().Page.Link
	if count > 0 {
		linkList, err = service.New(c).GetLinkList(request.PageNum, pageSize)
	}

	response := map[string]interface{}{
		"count":    count,
		"pageSize": pageSize,
		"pageNum":  request.PageNum,
		"linkList": &linkList,
	}
	api.Succ().SetData(&response).Output(c)
	return
}
