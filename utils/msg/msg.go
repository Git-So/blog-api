/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 * 为啥接口就一个方法?因为我就使用到一个方法
 */

package msg

// Type 消息类型
type Type int

const (
	// WeChat ...
	WeChat Type = iota
	// XMPP ...
	XMPP
)

var (
	// MsgInterface ...
	MsgInterface *Msg
	msgType      = WeChat
)

// Msg 消息接口
type Msg interface {
	// 只有一个接口,用完即销毁
	Send(message string) (stat bool)
}

// 新建消息接口实例
func new(msg Msg) *Msg {
	MsgInterface = &msg
	return MsgInterface
}

// Get 获取缓存实例
func Get() Msg {
	switch msgType {
	case WeChat:
		MsgInterface := &wechat{}
		new(MsgInterface)
	case XMPP:
		MsgInterface := &xmpp{}
		new(MsgInterface)
	}
	return *MsgInterface
}

// SetType 设置消息类型
func SetType(t Type) {
	msgType = t
}
