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
	"github.com/Git-So/blog-api/models"
)

// GetConfigList .
func GetConfigList(isAdmin bool, where []interface{}) (configList []*models.Config, err error) {
	// 查询数据
	var confing models.Config
	configList, err = confing.List(isAdmin, where)
	if isErrDB(err) {
		return nil, err
	}

	return
}

// ConfigUpdateAll 配置批量更新
func ConfigUpdateAll(cfs []*models.Config) (err error) {
	var cf models.Config
	return cf.UpdateAll(cfs)
}
