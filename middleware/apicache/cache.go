/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package apicache

import (
	"net/http"

	"github.com/Git-So/blog-api/utils/cache"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/helper"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var (
	jsonType  = "application/json"
	jsonpType = "application/javascript"
)

// Cache 路由缓存
func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取缓存状态
		cacheConf := conf.Get().Cache
		cacheInstance := cache.Get()

		// 缓存
		key := helper.Enbase64([]byte(c.Request.URL.String()))
		stat, _ := cacheInstance.Exists(key)
		logger.Debug("当前API缓存状态", cacheConf.APICacheStat)
		if c.Request.Method == "GET" && cacheConf.APICacheStat && stat {
			value, err := cacheInstance.Get(key)
			if err == nil {
				logger.Debug("使用缓存: ", value)
				c.Data(http.StatusOK, jsonType, []byte(value))
				c.Abort()
				return
			}
		}

		// 请求处理
		cacheResponse := NewResponse(c.Writer)
		c.Writer = cacheResponse

		// api处理
		c.Next()
		logger.Debug("请求处理：", string(cacheResponse.body.Bytes()))

		// 解析response
		response, err := responseDataParse(cacheResponse.body.Bytes())
		if err != nil {
			logger.Error("response 解析错误：", err)
		}

		// response返回正常非错信息缓存
		if response.Code == 200 {
			value := string(cacheResponse.body.Bytes())
			cacheInstance.SetEx(key, cacheConf.Expired, value)
		}

		return
	}
}
