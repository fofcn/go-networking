package user

import (
	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) uint {
	// userService := GetUserService(db)
	if claims, exists := c.Get("claims"); exists {
		if customClaims, ok := claims.(*CustomClaims); ok {
			return customClaims.UserId
		}

	}

	return 0
}
