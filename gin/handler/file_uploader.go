package handler

import (
	"go-networking/config"
	"go-networking/gin/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// swagger embed files
)

// gin-swagger middleware
// swagger embed files
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
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
