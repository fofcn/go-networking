package helloworld

import (
	"go-networking/ginh/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(public *gin.Engine) {
	public.GET("/", HelloWorld)
}

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, common.NewCommonResp("world", "hello"))
}
