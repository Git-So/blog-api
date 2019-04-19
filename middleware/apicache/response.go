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
	"bytes"

	"github.com/gin-gonic/gin"
)

// Response 用于缓存的实例
type Response struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// NewResponse 新建缓存实例
func NewResponse(ResponseWriter gin.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: ResponseWriter,
		body:           bytes.NewBufferString(""),
	}
}

// Write 写入值
func (resp *Response) Write(b []byte) (n int, err error) {
	if n, err = resp.body.Write(b); err != nil {
		return
	}
	return resp.ResponseWriter.Write(b)
}
