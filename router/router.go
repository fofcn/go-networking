package router

import (
	"fmt"
	"go-networking/ginh/global"
	"go-networking/ginh/user"
	"go-networking/log"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	log.Info("Init router")
	r.Use(gin.CustomRecovery(global.ErrorHandler))
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(func(c *gin.Context) {
		if c.FullPath() == "/login" {
			log.Info("Passing through the exact URI filter.")
			c.Next()
		} else {
			c.Next()
		}
	})

	// order is very import since auth protection will not work if the middleware is not registered.
	// user.InitAuthMiddleware(r)
	log.Info("Init auth middleware completed")

	v1 := r.Group("/api/v1")

	publicGroup := v1.Group("/")

	protectGroup := v1.Group("/")
	protectGroup.Use(user.AuthMiddleware())

	user.InitRouter(publicGroup, protectGroup)
	log.Info("Init router completed")

	log.Info("Init router completed")
}
