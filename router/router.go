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

	r.GET("/hello", handler.HelloWorld)
	r.POST("/login", handler.Login)
	r.POST("/upload", handler.UploadFile)
	r.GET("/user/:name", handler.RecvParameterFromPath)
}

func Add() {

}
