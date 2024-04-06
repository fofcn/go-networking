package user

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Valid 实现 `jwt.Claims` 接口，保证结构满足接口要求。
func (c CustomClaims) Valid() error {
	// 这里验证用户名和ID的正确性
	if c.Username == "" || c.UserId == 0 {
		return jwt.ErrInvalidKey
	}

	return nil
}
