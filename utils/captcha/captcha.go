/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package captcha

import (
	"strconv"

	"github.com/mojocn/base64Captcha"
	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/utils/cache"
	"github.com/gin-gonic/gin"
)

// Captcha 自动获取验证码
func Captcha(c *gin.Context) (captchaType, idKey, base64string string, stat bool) {
	// 查询该IP 5分钟内获取验证码次数
	var count int
	ip := c.ClientIP()
	key := cache.GetKey("blog_captcha_", ip)
	dataStat, err := cache.Get().Exists(key)
	if err != nil {
		logger.Error("验证码信息状态获取失败")
		return
	}
	if dataStat {
		// 获取当前次数
		value, err := cache.Get().Get(key)
		if err != nil {
			logger.Error("验证码信息状态获取失败")
			return
		}
		count, err = strconv.Atoi(value)
		if err != nil || count > 50 {
			logger.Warn("IP: ", ip, " 获取IP次数已达", value)
			return
		}
	}

	// 更新验证码次数以及时间
	err = cache.Get().SetEx(key, 300, count+1)
	if err != nil {
		logger.Error("验证码信息更新失败")
		return
	}

	// 获取验证码
	stat = true
	switch {
	case count < 8:
		captchaType = "Digit"
		idKey, base64string = ModeDigit()
	case count < 20:
		captchaType = "Arithmetic"
		idKey, base64string = ModeArithmetic()
	default:
		captchaType = "Audio"
		idKey, base64string = ModeAudio()
	}

	return
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(idkey, verifyValue string) (stat bool) {
	verifyResult := base64Captcha.VerifyCaptcha(idkey, verifyValue)
	if verifyResult {
		stat = true
	}
	return
}

// ModeDigit 数字字母混合验证码
func ModeDigit() (idKey, base64string string) {
	var config = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 6,
	}

	var cap base64Captcha.CaptchaInterface
	idKey, cap = base64Captcha.GenerateCaptcha("", config)
	base64string = base64Captcha.CaptchaWriteToBase64Encoding(cap)
	return
}

// ModeArithmetic 算术验证码
func ModeArithmetic() (idKey, base64string string) {
	var config = base64Captcha.ConfigCharacter{
		Height:             60,
		Width:              240,
		Mode:               base64Captcha.CaptchaModeArithmetic,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsUseSimpleFont:    true,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}

	var cap base64Captcha.CaptchaInterface
	idKey, cap = base64Captcha.GenerateCaptcha("", config)
	base64string = base64Captcha.CaptchaWriteToBase64Encoding(cap)
	return
}

// ModeAudio 声音验证码
func ModeAudio() (idKey, base64string string) {
	var config = base64Captcha.ConfigAudio{
		CaptchaLen: 6,
		Language:   "zh",
	}

	var cap base64Captcha.CaptchaInterface
	idKey, cap = base64Captcha.GenerateCaptcha("", config)
	base64string = base64Captcha.CaptchaWriteToBase64Encoding(cap)
	return
}
