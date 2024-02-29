package handler

import (
	"go-networking/config"
	"go-networking/gin/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	err := c.SaveUploadedFile(file, config.GetAppStorePath()+file.Filename)
	if err != nil {
		log.Printf("Upload file failed: %s, reason: %s", file.Filename, err)
		return
	}
	c.IndentedJSON(http.StatusOK, common.NoDataSuccessResposne)
}
