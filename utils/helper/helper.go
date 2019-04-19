/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package helper

import (
	"os"
	"strings"

	"github.com/wonderivan/logger"
)

// CreatePath 创建给定路径
func CreatePath(path string, checkPath string) (err error) {
	paths := make(map[int]string, 2)

	index := strings.LastIndex(checkPath, "/")
	if index > 0 {
		// 目录
		paths[0] = checkPath[0:index]
		// 文件
		paths[1] = checkPath[index:]
	}
	logger.Debug(paths)

	for key, value := range paths {
		// 忽略空路径
		if value == "" {
			continue
		}

		// 保存路径
		path += value
		logger.Debug("保存路径: ", path)

		// 路径是否存在
		isUsable, _, err := IsUsablePath(path)
		if isUsable {
			continue
		}
		if err != nil {
			logger.Error("获取路径可用性失败：", err)
			return err
		}

		// 文件
		if key == 1 {
			f, err := os.Create(path)
			if err != nil {
				logger.Error("新建文件失败 ", err)
				return err
			}
			defer f.Close()
			logger.Info("文件", path, "已创建")
			continue
		}

		// 目录路径
		err = os.MkdirAll(path, 0755)
		if err != nil {
			logger.Error("新建文件失败 ", err)
			return err
		}
		logger.Info("目录", path, "已创建")
	}
	return nil
}

// WriteFile 写入内容到文件
func WriteFile(path string, text string) (err error) {
	// 创建路径
	err = CreatePath("", path)
	if err != nil {
		logger.Error("获取文件失败：", err)
	}

	// 打开文件
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("打开文件失败: ", err)
		return
	}
	defer file.Close()

	// 写入文件数据
	_, err = file.WriteString(text)
	if err != nil {
		logger.Error("写入文件数据失败: ", err)
		return
	}

	return
}

// IsUsablePath 路径是否存在且可操纵
func IsUsablePath(path string) (isUsable, isNotExist bool, err error) {
	// 路径是否存在
	_, err = os.Stat(path)
	switch {
	case os.IsNotExist(err):
		logger.Info("文件", path, "不存在")
		return false, true, nil
	case os.IsPermission(err):
		logger.Error("权限不足：", err)
		return
	case err != nil:
		logger.Error("获取路径状态失败：", err)
		return
	}
	return true, false, nil
}
