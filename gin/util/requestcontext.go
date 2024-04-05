package util

import "github.com/gin-gonic/gin"

// 该函数用于获取gin.Context中的claims，并将其解析成userId和username返回。

// 参数：

// c：gin.Context，请求的上下文信息。
// 返回值：

// int：用户ID。
// string：用户名。
// 函数首先通过c.Get("claims")获取到claims，然后将其断言为CustomClaims类型，并返回其中的UserId和Username字段的值。
// 函数获取gin.Context中的claims，并解析成userId和username
func GetUserIdAndUsername(c *gin.Context) (int, string) {
	claims, _ := c.Get("claims")
	customClaims := claims.(*CustomClaims)
	return customClaims.UserId, customClaims.Username
}
