/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 * 错误处理
 */

package e

const (
	// Success 响应成功
	Success = 200
	// Fail 响应失败
	Fail = 500
	// InvalidParams 无效参数
	InvalidParams = 400
	// NotFound 404
	NotFound = 404
	// ErrDB 数据库响应错误
	ErrDB = 30001
	// ErrAuth 检权失败
	ErrAuth = 50001
	// ErrCaptcha 验证码错误
	ErrCaptcha = 50002
	// ErrNotFoundData 未找到数据
	ErrNotFoundData = 60001
	// ErrNotFoundArticle 未找到文章
	ErrNotFoundArticle = 60010
	// ErrExistsArticle 文章已存在
	ErrExistsArticle = 60011
	// ErrNotFoundSubject 未找到专题
	ErrNotFoundSubject = 60020
	// ErrExistsSubject 专题已存在
	ErrExistsSubject = 60021
	// ErrExistsTag 标签已存在
	ErrExistsTag = 60030
	// ErrUsedTag 标签已被使用
	ErrUsedTag = 60031
	// ErrExistsLink 友链已存在
	ErrExistsLink = 60040
	// ErrNotFoundLink 未找到友链
	ErrNotFoundLink = 60041
	// ErrNotFoundComment 未找到评论
	ErrNotFoundComment = 60050
	// ErrNotFoundParentComment 未找到父评论
	ErrNotFoundParentComment = 60051
)

// Msg 错误信息
var Msg = map[int]string{
	Success:                  "响应成功",
	Fail:                     "响应失败",
	InvalidParams:            "无效参数",
	NotFound:                 "请求地址不存在",
	ErrDB:                    "数据库响应错误",
	ErrAuth:                  "检权失败",
	ErrCaptcha:               "验证码错误",
	ErrNotFoundData:          "未找到数据",
	ErrNotFoundArticle:       "未找到文章",
	ErrNotFoundSubject:       "未找到专题",
	ErrExistsArticle:         "文章已存在",
	ErrExistsTag:             "标签已存在",
	ErrUsedTag:               "标签已被使用",
	ErrExistsSubject:         "专题已存在",
	ErrExistsLink:            "友链已存在",
	ErrNotFoundLink:          "未找到友链",
	ErrNotFoundComment:       "未找到评论",
	ErrNotFoundParentComment: "未找到父评论",
}
