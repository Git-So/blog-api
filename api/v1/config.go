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
	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/service"
	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/e"
	"github.com/gin-gonic/gin"
)

// MeInfo .
func MeInfo(c *gin.Context) {
	var configList []*models.Config
	configList, err := service.GetConfigList(isAdmin(c), []interface{}{})
	isNotFound, isErr := api.IsServiceError(c, err)
	if isErr {
		return
	}
	if isNotFound {
		api.New(e.ErrNotFoundData).Output(c)
		return
	}

	response := map[string]interface{}{
		"configList": &configList,
	}
	api.Succ().SetData(&response).Output(c)
	return
}

// UpdateConfig .
func UpdateConfig(c *gin.Context) {
	// request
	var request []*models.Config
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// validate
	for _, val := range request {
		if len(val.Key) < 1 {
			api.ErrValidate("Key 不能为空").Output(c)
			return
		}
	}

	err := service.ConfigUpdateAll(request)
	if _, isErr := api.IsServiceError(c, err); isErr {
		return
	}
	api.Succ().SetMsg("更新成功").Output(c)
	return
}
