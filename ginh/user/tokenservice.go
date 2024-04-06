package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
}

func (ts *TokenService) GenerateToken(userid uint, username string) (string, error) {
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
