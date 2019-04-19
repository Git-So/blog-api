/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package jwt

import (
	"errors"
	"time"

	"github.com/Git-So/blog-api/utils/conf"
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims ...
type CustomClaims struct {
	Code string
	jwt.StandardClaims
}

// CreateToken 生成Token
func CreateToken(loginCode string) (string, error) {
	claims := &CustomClaims{
		Code: loginCode,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),
			ExpiresAt: int64(time.Now().Unix() + conf.Get().Jwt.Expired),
			Issuer:    "So",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJwtSecret())
}

// ParseToken 解析
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if err != nil {
		return nil, validateErr(err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, validateErr(err)
}

func getJwtSecret() []byte {
	return []byte(conf.Get().Jwt.Secret)
}

func validateErr(err error) error {
	errString := "无效Token"
	if e, ok := err.(*jwt.ValidationError); ok {
		switch e.Errors {
		case jwt.ValidationErrorMalformed:
			errString = "Token 格式错误"
		case jwt.ValidationErrorExpired:
			errString = "Token已过期"
		case jwt.ValidationErrorNotValidYet:
			errString = "Token签名验证失败"
		}
	}
	return errors.New(errString)
}
