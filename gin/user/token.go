package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId   int    `json:"user_id"`
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

type TokenService struct {
}

func (ts *TokenService) GenerateToken(userid int, username string) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 设置token过期时间
			Issuer:    "nasapp",                                           // 设置token发行人
		},
		Username: username,
		UserId:   userid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret")) //  todo : secret
}
