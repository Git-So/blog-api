/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package api

import (
	"net/http"

	"github.com/wonderivan/logger"

	"github.com/gin-gonic/gin"
)

// Response response数据
type Response struct {
	Code int
	Msg  string
	Data interface{} `json:"data,omitempty"`
}

// SetCode ...
func (resp *Response) SetCode(code int) *Response {
	resp.Code = code
	return resp
}

// SetMsg ...
func (resp *Response) SetMsg(msg string) *Response {
	resp.Msg = msg
	return resp
}

// SetData ...
func (resp *Response) SetData(data interface{}) *Response {
	resp.Data = data
	return resp
}

// Output response输出
func (resp *Response) Output(c *gin.Context, httpStatusCode ...int) {
	httpCode := http.StatusOK
	// isHttpStatusCode
	if len(httpStatusCode) > 0 {
		httpCode = httpStatusCode[0]
	}

	logger.Debug("resp", resp)
	logger.Debug("httpCode", httpCode)
	c.JSON(httpCode, resp)
	return
}
