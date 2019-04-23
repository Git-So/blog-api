/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package service

import (
	"github.com/gin-gonic/gin"
)

// Service struct
type Service struct {
	// 是否使用缓存,不使用将更新缓存
	IsCache bool
}

// New 实例化
func New(c *gin.Context) *Service {
	var s Service
	// IsCache
	s.IsCache = true
	if value := c.GetHeader("Is-Cache"); value == "off" {
		s.IsCache = false
	}
	return &s
}

// 是否存在数据库错误
func isErrDB(err error) bool {
	// Not error
	if err == nil {
		return false
	}
	return true
}
