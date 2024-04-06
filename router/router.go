package router

import (
	"fmt"
	"go-networking/gin/global"
	"go-networking/gin/user"
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

	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	user.InitRouter(auth)
	user.InitAuthMiddleware(r)

	log.Info("Init router completed")
}
