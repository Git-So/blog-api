/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package routers

import (
	"github.com/Git-So/blog-api/utils/api"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/e"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

// Get 获取路由实例
func Get() *gin.Engine {
	if router == nil {
		new()
	}
	return router
}

func new() {
	config := conf.Get()

	gin.SetMode(config.Dev.RunMode)

	router = gin.New()

	// 中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(cachemiddleware.Cache())

	// 404
	router.NoRoute(NotFound)

	// api.v1
	v1router(router)

	return
}

// NotFound 无效路由
func NotFound(c *gin.Context) {
	api.Err(e.NotFound).Output(c)
}
