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
	"encoding/base64"
	"math/rand"
	"time"
)

// Enbase64 ...
func Enbase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Debase64 ...
func Debase64(code string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(code)
}

// GetRandomString ...
func GetRandomString(length int) string {
	bytes := []byte("1234567890qwertyuioplkjhgfdsazxcvbnmMNBVCXZASDFGHJKLPOIUYTREWQ")
	result := make([]byte, 0)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
