/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 * 啊啊啊,这个文件有点乱,还是用 goto 了,少用吧
 */

package v1

import (
	"strconv"

	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/cache"
	"github.com/Git-So/blog-api/utils/e"
	"github.com/Git-So/blog-api/utils/helper"
	"github.com/Git-So/blog-api/utils/msg"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

// GetLoginCode 获取登入验证码
func GetLoginCode(c *gin.Context) {
	var msgType msg.Type
	var stat bool

	// 获取发送类型
	messageType := c.DefaultQuery("type", "wechat")
	switch messageType {
	case "wechat":
		msgType = msg.WeChat
	default:
		msgType = msg.XMPP
	}

	// 获取登入码
	loginCode := helper.GetRandomString(6)

	// 获取发送状态
	var count int
	ip := c.ClientIP()
	logger.Debug(ip)
	key := cache.GetKey("blog_loginCodeNum_", ip)
	dataStat, err := cache.Get().Exists(key)
	if err != nil {
		logger.Error("登入码状态获取失败")
		goto fail
	}
	if dataStat {
		// 获取当前次数
		value, err := cache.Get().Get(key)
		if err != nil {
			logger.Error("登入码状态获取失败")
			goto fail
		}
		count, err = strconv.Atoi(value)
		if err != nil {
			logger.Error("登入码状态获取失败")
			goto fail
		}
		if count > 50 {
			logger.Warn("IP: ", ip, " 获取次数已达", value)
			api.Err(e.Fail, "获取登入码太过频繁").Output(c)
			return
		}
	}

	// 更新验证码次数以及时间
	err = cache.Get().SetEx(key, 300, count+1)
	if err != nil {
		logger.Error("验证码信息更新失败")
		goto fail
	}

	// 发送验证码
	msg.SetType(msgType)
	stat = msg.Get().Send("你的登入码为: " + loginCode)
	if !stat {
		goto fail
	}

	// 写入登入验证码
	cache.Get().SetEx("blog_loginCode", 300, loginCode)

	api.Succ().SetMsg("获取登入码成功").Output(c)
	return

fail:
	// 失败
	api.Err(e.Fail, "获取登入码失败").Output(c)
	return
}
