/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package cache_test

import (
	"testing"

	"github.com/Git-So/blog-api/utils/cache"

	"github.com/smartystreets/goconvey/convey"
)

func Test_redis(t *testing.T) {
	convey.Convey("redis:", t, func() {
		obj := cache.New()
		obj.Connect()
		cache.Get().Set("Name", "123")
	})
}
