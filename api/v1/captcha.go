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
	"github.com/Git-So/blog-api/utils/captcha"
	"github.com/Git-So/blog-api/utils/e"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

// Captcha 接口提交数据
type Captcha struct {
	Key     string `validate:"required"`
	Captcha string `validate:"required"`
}

// ConfigValidation 验证配置
func (cap Captcha) ConfigValidation(v *validate.Validation) {

	// 字段名称
	v.AddTranslates(validate.MS{
		"Key":     "验证Key",
		"Captcha": "验证码",
	})

	// 错误信息
	v.AddMessages(validate.MS{
		"required": "{field}不能为空",
	})
}

// GetCaptcha 获取验证码
func GetCaptcha(c *gin.Context) {
	captchaType, idKey, base64string, stat := captcha.Captcha(c)
	if !stat {
		api.Err(e.ErrCaptcha, "获取验证码失败,可能太过于频繁").Output(c)
		return
	}

	// data
	response := map[string]string{
		"type":   captchaType,
		"key":    idKey,
		"base64": base64string,
	}
	api.Succ().SetData(response).Output(c)
	return
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(c *gin.Context) {
	// request
	var request Captcha
	if err := c.ShouldBind(&request); err != nil {
		api.ErrValidate().Output(c)
		return
	}

	// 验证
	v := validate.Struct(request, "verify")
	if !v.Validate() {
		api.ErrValidate(v.Errors.One()).Output(c)
		return
	}

	// 验证验证码
	stat := captcha.VerifyCaptcha(request.Key, request.Captcha)
	if !stat {
		api.New(e.ErrCaptcha).Output(c)
		return
	}
	api.Succ().Output(c)
	return
}
