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
	"fmt"

	"github.com/Git-So/blog-api/utils/helper"
)

var configYaml = fmt.Sprintf(`#开发相关
dev:
  run_mode: debug # release debug

#页数配置
page:
  article: 8
  hot_article: 15
  comment: 10
  subject: 8
  tag: 200
  link: 10

#服务配置
server:
  port: 8099
  read_timeout: 60
  write_timeout: 60

#数据库配置
database:
  type: mysql # mysql sqlite3
  host: 127.0.0.1
  port: 3306
  user: root
  passwd: root
  name: blog
  table_prefix: blog_

#缓存配置
cache:
  type: redis
  host: 127.0.0.1
  port: 6379
  expired: 259200 # 秒
  prefix: blog_
  api_cache_state: true

#Jwt配置
jwt:
  secret: %s
  expired: 7200 # 秒

#WeChat配置
wechat:
  app_id: 
  app_secret: 
  touser: 
  template_id: 
  color: 000fff


#XMPP配置
xmpp:
  host: xmpp.jp:5222
  touser: sooo.site@xmpp.jp
  user: 
  passwd: 
  no_tls: true
  session: true
  status: xa
  status_messgae: I for one welcome our new codebot overlords.
`, helper.GetRandomString(50))
