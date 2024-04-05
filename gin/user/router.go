package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.POST("/login", Login)
	r.POST("/register", Register)
}

func Login(c *gin.Context) {
	var loginCmd LoginCmd
	if err := c.ShouldBindJSON(&loginCmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := UserService.Login(loginCmd.Username, loginCmd.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func Register(c *gin.Context) {
	UserService.Register(c.PostForm("username"), c.PostForm("password"))
	c.JSON(http.StatusOK, gin.H{"message": "register success"})
}
