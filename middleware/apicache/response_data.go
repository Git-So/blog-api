/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package apicache

import "encoding/json"

// ResponseData 接口数据
type ResponseData struct {
	Code int
	Msg  string
}

func responseDataParse(data []byte) (resp *ResponseData, err error) {
	err = json.Unmarshal(data, &resp)
	return
}
