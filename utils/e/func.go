/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package e

// GetMsg 获取错误信息
func GetMsg(code int) string {
	if msg, ok := Msg[code]; ok {
		return msg
	}
	return Msg[Fail]
}
