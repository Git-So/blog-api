/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 * API 处理
 */

package api

import (
	"database/sql"
	"net/http"

	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/utils/e"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// New 新建实例
func New(code int, messages ...string) *Response {
	var msg string

	// isMsg
	if len(messages) > 0 {
		msg = messages[0]
	} else {
		msg = e.GetMsg(code)
	}

	return &Response{
		Code: code,
		Msg:  msg,
	}
}

// Succ 默认成功
func Succ() *Response {
	return New(e.Success)
}

// Err 默认错误
func Err(code int, messages ...string) *Response {
	logger.Info("请求错误")
	return New(code, messages...)
}

// ErrValidate 验证错误 快捷函数
func ErrValidate(messages ...string) *Response {
	return Err(e.InvalidParams, messages...)
}

// IsServiceError 查询数据错误检查并输出
func IsServiceError(c *gin.Context, err error) (isNotFound bool, isErr bool) {
	// Not error
	if err == nil {
		return false, false
	}

	// isNotFound
	if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
		return true, false
	}

	// isErr
	New(e.ErrDB).Output(c, http.StatusInternalServerError)

	// log
	logger.Error("查询数据出错：", err)

	return true, true
}
