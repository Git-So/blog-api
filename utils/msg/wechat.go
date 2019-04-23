/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 * 使用微信快,但测试号可能过期,看微信官方心情把
 */

package msg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/utils/cache"
	"github.com/Git-So/blog-api/utils/conf"
)

// 缓存 key
var key = cache.GetKey("blog_wechat_accessToken")

type wechat struct {
}

// ErrAccessToken 错误信息
// var ErrAccessToken = map[int]string{
// 	-1:    "系统繁忙",
// 	0:     "请求成功",
// 	40001: "AppSecret错误",
// 	40002: "请确保grant_type字段值为client_credential",
// 	40164: "调用接口的IP地址不在白名单中",
// }

// AccessToken 全局唯一接口调用凭据
// succ{"access_token":"ACCESS_TOKEN","expires_in":7200}
// fail {"errcode":40013,"errmsg":"invalid appid"}
type AccessToken struct {
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   int    `json:"expires_in"`   // 凭证有效时间，单位：秒
	Errcode     int    `json:"errcode"`      // 错误码
	Errmsg      int    `json:"errmsg"`       // 错误信息
}

// Template 模板消息
type Template struct {
	ToUser     string `json:"touser"`
	TemplateID string `json:"template_id"`
	Data       Data   `json:"data"`
}

// Data 模板数据
type Data struct {
	Value Value `json:"value"`
}

// Value 数据实例
type Value struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// TemplateResponse 请求发送模板返回数据
type TemplateResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgID   int    `json:"msgid"`
}

// Send 发送微信模板消息
func (we *wechat) Send(message string) (stat bool) {
	logger.Debug("发送微信模板消息")
	// 获取 accessToken
	isExist, _ := cache.Get().Exists(key)
	logger.Debug("发送微信模板消息")
	if !isExist {
		stat = we.getAccessToken()
		if !stat {
			logger.Error("生成微信全局唯一接口调用凭据失败")
			return
		}
	}
	logger.Debug("获取 AccessToken")
	token, err := cache.Get().Get(key)
	if err != nil {
		logger.Error("缓存获取微信全局唯一接口调用凭据失败")
	}

	// 整理传送数据
	wechatConf := conf.Get().WeChat
	data := Template{
		ToUser:     wechatConf.ToUser,
		TemplateID: wechatConf.TemplateID,
		Data: Data{
			Value: Value{
				Value: message,
				Color: "#" + wechatConf.Color,
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("数据转换错误")
		return
	}

	// 发送数据
	logger.Debug("发送数据")
	uri := `https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=` + token
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("请求发送模板消息出错", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("获取 body 数据失败", err)
		return
	}

	// 解析数据
	logger.Debug("解析数据")
	var response TemplateResponse
	if err = json.Unmarshal(body, &response); err != nil {
		logger.Error("response数据错误", err)
		return
	}
	logger.Debug("判断数据")
	if response.ErrCode != 0 {
		// 保存请求与结果到日志
		logger.Error("请求地址:", uri)
		logger.Error("请求数据:", string(jsonData))
		logger.Error("请求结果:", string(body))
		return
	}
	return true
}

// getAccessToken 获取微信全局唯一接口调用凭据
func (we *wechat) getAccessToken() (stat bool) {
	// 地址
	uri := fmt.Sprintf(
		`https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s`,
		conf.Get().WeChat.AppID,
		conf.Get().WeChat.AppSecret,
	)
	logger.Debug("发送微信模板消息")
	// 请求
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Get(uri)
	if err != nil {
		logger.Error("获取AccessToken出错")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("获取 body 数据失败")
		return
	}

	// 解析数据
	var token AccessToken
	if err = json.Unmarshal(body, &token); err != nil {
		logger.Error("response数据错误")
		return
	}

	// 获取过期时间
	if token.ExpiresIn < 1 {
		// 保存请求与结果到日志
		logger.Error("请求地址:", uri)
		logger.Error("请求结果:", string(body))
		return
	}

	// 缓存 accessToken 时间少30s 防止网络请求浪费时间
	err = cache.Get().SetEx(key, token.ExpiresIn-30, token.AccessToken)
	if err != nil {
		logger.Error("保存 accessToken 失败")
		return
	}
	return true
}
