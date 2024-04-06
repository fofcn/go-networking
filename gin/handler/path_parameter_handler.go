package handler

// import (
// 	"go-networking/gin/common"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type RecvParamDto struct {
// 	Name string `json:"name"`
// }

// func RecvParameterFromPath(c *gin.Context) {
// 	name := c.Param("name")
// 	recvParamResp := RecvParamDto{Name: name}
// 	c.IndentedJSON(http.StatusOK, common.SuccessResponse{Data: recvParamResp})
// }
