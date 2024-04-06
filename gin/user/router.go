package user

import (
	"go-networking/gin/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
	r.POST("/login", Login)
	r.POST("/register", Register)
}

var userService = &UserService{}

func Login(c *gin.Context) {
	var loginCmd LoginCmd
	if err := c.ShouldBindJSON(&loginCmd); err != nil {
		c.IndentedJSON(http.StatusOK, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}

	token, err := userService.DoLogin(loginCmd.Username, loginCmd.Password)
	if err != nil {
		c.IndentedJSON(http.StatusOK, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}

	c.IndentedJSON(http.StatusOK,
		common.CommonResp{
			Data:    &LoginDto{Token: token, Exp: 0},
			Message: "Username or password is incorrect."},
	)

}

func Register(c *gin.Context) {
	userService.DoRegister(c.PostForm("username"), c.PostForm("password"))
	c.JSON(http.StatusOK, gin.H{"message": "register success"})
}
