package file

import (
	"go-networking/db"
	"go-networking/ginh/common"
	"go-networking/ginh/user"
	"go-networking/ginh/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(protect *gin.RouterGroup) {
	protect.GET("/file/list", db.WithDB(ListFile))
	protect.GET("/file/:id", db.WithDB(GetFile))
	protect.POST("/file", db.WithDB(UploadFile))
}

func ListFile(c *gin.Context, db *gorm.DB) {
	page := c.Param("page")
	pageSize := c.Param("page_size")
	fileId := c.Param("file_id")

	if len(page) == 0 {
		page = "1"
	}
	if len(pageSize) == 0 {
		pageSize = "10"
	}
	if len(fileId) == 0 {
		fileId = "0"
	}
	service := NewFileService(db)
	cmd := &ListFileCmd{
		Pagenation: common.Pagenation{
			Page:     util.StringToInt(page),
			PageSize: util.StringToInt(pageSize),
		},
		FileId: util.StringToUint(fileId),
	}

	service.ListFile(cmd, user.GetUserId(c))

	var files []FileInfo
	db.Find(&files)
	c.JSON(200, files)
}

func GetFile(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var file FileInfo
	db.First(&file, id)
	c.JSON(200, file)
}

func UploadFile(c *gin.Context, db *gorm.DB) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, file)
}
