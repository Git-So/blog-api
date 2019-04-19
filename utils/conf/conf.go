/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package conf

import (
	"io/ioutil"
	"log"

	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/utils/helper"
	yaml "gopkg.in/yaml.v2"
)

var (
	// ConfigFile 配置文件地址
	ConfigFile = ".blog/conf.yaml"

	// 配置数据
	config *ConfigData
)

// Get 获取配置实例
func Get() *ConfigData {
	if config == nil {
		new()
	}
	return config
}

func new() {
	// 检查配置文件
	_, isNotExists, err := helper.IsUsablePath(ConfigFile)
	if err != nil {
		logger.Fatal("配置文件初始化失败：", err)
		return
	}
	if isNotExists {
		// 配置文件初始化
		err := helper.WriteFile(ConfigFile, configYaml)
		if err != nil {
			logger.Fatal("配置文件初始化无法写入：", err)
			return
		}
	}

	// 读取配置文件
	confFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatalf("config file err: %v", err)
	}

	config = &ConfigData{}
	err = yaml.Unmarshal(confFile, &config)
	if err != nil {
		log.Fatalf("config content err: %v", err)
	}
	return
}
