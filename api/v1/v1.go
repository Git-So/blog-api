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
	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func isAdmin(c *gin.Context) (isAdmin bool) {
	// 验证身份
	if _, stat := c.Get("claims"); stat {
		// claims := value.(*jwt.CustomClaims)
		return true
	}

	// Token 验证身份
	if token := c.GetHeader("Authorization"); token != "" {
		// 验证Token
		if _, err := jwt.ParseToken(token); err == nil {
			return true
		}
	}
	return false
}

func request(c *gin.Context, obj interface{}) (isErr bool) {
	// request
	if err := c.ShouldBind(&obj); err != nil {
		api.ErrValidate().Output(c)
		return true
	}

	// 验证
	v := validate.Struct(obj, "CreateArticle")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return true
	}
	return
}
