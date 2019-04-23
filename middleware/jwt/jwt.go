/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package jwt

import (
	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/cache"
	"github.com/Git-So/blog-api/utils/e"
	"github.com/Git-So/blog-api/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

// Auth 检权
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			// 验证Token
			claims, err := jwt.ParseToken(token)
			if err != nil {
				api.New(e.ErrAuth, err.Error()).Output(c)
				c.Abort()
				return
			}
			c.Set("claims", claims)
			c.Next()
			return
		}

		// 检查是否存在登入码
		type loginCode struct {
			LoginCode string `json:"loginCode"`
		}
		var request loginCode
		if err := c.ShouldBind(&request); err == nil {
			// 安全码不存在
			stat, _ := cache.Get().Exists("blog_loginCode")
			if !stat {
				api.ErrValidate("请先获取安全码").Output(c)
				c.Abort()
				return
			}

			// 验证安全码
			if code, _ := cache.Get().Get("blog_loginCode"); request.LoginCode != code {
				logger.Debug(code)
				api.ErrValidate("安全码错误").Output(c)
				c.Abort()
				return

			}

			// 创建Token
			token, err := jwt.CreateToken(request.LoginCode)
			if err != nil {
				api.New(e.Fail).Output(c)
				c.Abort()
				return
			}

			// 返回Token数据
			response := map[string]interface{}{
				"token": token,
			}
			api.Succ().SetData(&response).Output(c)
			c.Abort()
		}
		return
	}
}
