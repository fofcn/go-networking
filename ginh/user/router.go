package user

import (
	"go-networking/db"
	"go-networking/ginh/common"
	"go-networking/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(r *gin.RouterGroup) {
	r.POST("/login", db.WithDB(Login))
	r.POST("/register", db.WithDB(Register))
}

func Login(c *gin.Context, db *gorm.DB) {
	var loginCmd LoginCmd
	if err := c.ShouldBindJSON(&loginCmd); err != nil {
		log.Errorf("Login failed: %v", err)
		c.IndentedJSON(http.StatusOK, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}

	userService := GetUserService(db)
	token, err := userService.DoLogin(loginCmd.Username, loginCmd.Password)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		c.IndentedJSON(http.StatusOK, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}

	c.IndentedJSON(http.StatusOK,
		common.CommonResp{
			Data:    &LoginDto{Token: token, Exp: 0},
			Message: "ok"},
	)

}

func Register(c *gin.Context, db *gorm.DB) {
	userService := GetUserService(db)
	var registerCmd RegisterCmd
	if err := c.ShouldBindJSON(&registerCmd); err != nil {
		log.Errorf("Register failed: %v", err)
		c.IndentedJSON(http.StatusOK, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}

	err := userService.DoRegister(&registerCmd)
	if err != nil {
		log.Errorf("Register failed: %v", err)
		c.IndentedJSON(http.StatusOK, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}
	c.JSON(http.StatusOK, common.NoDataSuccessResposne)
}
