package user

import (
	"fmt"
	"go-networking/ginh/common"
	"go-networking/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func InitAuthMiddleware(r *gin.Engine) {
	r.Use(AuthMiddleware())
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("AuthMiddleware")
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, common.CommonResp{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		// 这里使用jwt-go解析Claims
		token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 确保token方法与预期一致
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.CommonResp{
				Message: "Unauthorized",
			})
			c.Abort()
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			// 添加claims到上下文
			c.Set("claims", claims)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, common.CommonResp{
				Message: "Unauthorized",
			})
			c.Abort()
		}

	}
}
