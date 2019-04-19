/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package service

func isErrDB(err error) bool {
	// Not error
	if err == nil {
		return false
	}
	return true
}
