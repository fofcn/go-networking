package handler

import (
	"go-networking/config"
	"go-networking/gin/common"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
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
