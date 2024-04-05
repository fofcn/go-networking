package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitAuthMiddleware() {
	r.Use(AuthMiddleware())
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil {
			c.JSON()
		}

		// 添加claims到上下文
		c.Set("claims", claims)

		c.Next()
	}
}
