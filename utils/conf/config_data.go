/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package conf

import "time"

// ConfigData 配置数据
type ConfigData struct {
	Dev      *Dev      `yaml:"dev"`
	Server   *Server   `yaml:"server"`
	Database *Database `yaml:"database"`
	Cache    *Cache    `yaml:"cache"`
	Page     *Page     `yaml:"page"`
	Jwt      *Jwt      `yaml:"jwt"`
	WeChat   *WeChat   `yaml:"wechat"`
	XMPP     *XMPP     `yaml:"xmpp"`
}

// Dev 开发相关配置
type Dev struct {
	RunMode string `yaml:"run_mode"` //运行模式
}

// Server 服务配置
type Server struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// Database 数据库配置
type Database struct {
	Type        string `yaml:"type"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Passwd      string `yaml:"passwd"`
	Name        string `yaml:"name"`
	TablePrefix string `yaml:"table_prefix"`
}

// Cache 缓存配置
type Cache struct {
	Type         string `yaml:"type"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Expired      int64  `yaml:"expired"`
	Prefix       string `yaml:"prefix"`
	APICacheStat bool   `yaml:"api_cache_state"`
}

// Page 页数配置
type Page struct {
	Article    uint `yaml:"article"`
	HotArticle uint `yaml:"hot_article"`
	Comment    uint `yaml:"comment"`
	Subject    uint `yaml:"subject"`
	Tag        uint `yaml:"tag"`
	Link       uint `yaml:"link"`
}

// Jwt 应用配置
type Jwt struct {
	Secret  string `yaml:"secret"`  //  安全码
	Expired int64  `yaml:"expired"` //  有效期
}

// WeChat 微信开发账号配置
type WeChat struct {
	AppID      string `yaml:"app_id"`
	AppSecret  string `yaml:"app_secret"`
	ToUser     string `yaml:"touser"`
	TemplateID string `yaml:"template_id"`
	Color      string `yaml:"color"`
}

// XMPP xmpp服务设置
type XMPP struct {
	Host          string `yaml:"host"`
	ToUser        string `yaml:"touser"`
	User          string `yaml:"user"`
	Passwd        string `yaml:"passwd"`
	NoTLS         bool   `yaml:"no_tls"`
	Session       bool   `yaml:"session"`
	Status        string `yaml:"status"`
	StatusMessage string `yaml:"status_messgae"`
}
