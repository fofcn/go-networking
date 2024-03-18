package router

import (
	"fmt"
	"go-networking/gin/global"
	"go-networking/gin/handler"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
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
			log.Println("Passing through the exact URI filter.")
			c.Next()
		} else {
			c.Next()
		}
	})

	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/login", handler.Login)
	// auth.DELETE("/logout", handler.Logout)

	file := v1.Group("/file")
	file.POST("/upload", handler.UploadFile)

	user := v1.Group("/user")
	user.GET("/user/:name", handler.RecvParameterFromPath)

	v1.GET("/hello", handler.HelloWorld)
}
