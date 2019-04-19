/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 * 使用 Xmpp 慢,我使用的是xxx@xmpp.jp
 */

package msg

import (
	"github.com/wonderivan/logger"

	"github.com/Git-So/blog-api/utils/conf"
	goxmpp "github.com/mattn/go-xmpp"
)

type xmpp struct {
}

var talk *goxmpp.Client

// Send 发送xmpp信息
func (xp *xmpp) Send(message string) (stat bool) {
	// 获取实例
	if stat := xp.talk(); !stat {
		logger.Error("无法获取XMPP连接实例")
		return false
	}

	// 发送消息
	newChat := goxmpp.Chat{
		Remote: conf.Get().XMPP.ToUser,
		Type:   "Chat",
		Text:   message,
	}
	talk.Send(newChat)

	return true
}

func (xp *xmpp) talk() (stat bool) {
	var err error
	xmppConf := conf.Get().XMPP
	options := goxmpp.Options{
		Host:          xmppConf.Host,
		User:          xmppConf.User,
		Password:      xmppConf.Passwd,
		NoTLS:         xmppConf.NoTLS,
		Debug:         conf.Get().Dev.RunMode == "dev",
		Session:       xmppConf.Session,
		Status:        xmppConf.Status,
		StatusMessage: xmppConf.StatusMessage,
	}

	talk, err = options.NewClient()

	if err != nil {
		logger.Error("XMPP 连接失败")
		return
	}

	return true
}
