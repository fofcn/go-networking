package user

import (
	"go-networking/db"
	"go-networking/ginh/common"
	"go-networking/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"gorm.io/gorm"
)

var (
	g          = galidator.New()
	customizer = g.Validator(LoginCmd{})
)

func InitRouter(public *gin.RouterGroup, protect *gin.RouterGroup) {
	public.POST("/auth/login", db.WithDB(Login))
	public.POST("/auth/register", db.WithDB(Register))
	protect.GET("/auth/userinfo", db.WithDB(GetUserInfo))
}

func Login(c *gin.Context, db *gorm.DB) {
	var loginCmd LoginCmd
	if err := c.ShouldBindJSON(&loginCmd); err != nil {
		log.Errorf("Login failed: %v", err)
		c.IndentedJSON(http.StatusBadRequest, common.CommonResp{Data: "Failed", Message: customizer.DecryptErrors(err)})
		return
	}

	userService := GetUserService(db)
	token, err := userService.DoLogin(loginCmd.Username, loginCmd.Password)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		c.IndentedJSON(http.StatusUnauthorized, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
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
		c.IndentedJSON(http.StatusBadRequest, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}

	err := userService.DoRegister(&registerCmd)
	if err != nil {
		log.Errorf("Register failed: %v", err)
		c.IndentedJSON(http.StatusBadRequest, common.CommonResp{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}
	c.JSON(http.StatusOK, common.NoDataSuccessResposne)
}

func GetUserInfo(c *gin.Context, db *gorm.DB) {
	// userService := GetUserService(db)
	if claims, exists := c.Get("claims"); exists {
		if customClaims, ok := claims.(*CustomClaims); ok {
			userinfoDto := UserInfoDto{
				Username: customClaims.Username,
				UserId:   customClaims.UserId,
			}
			c.IndentedJSON(http.StatusOK, common.CommonResp{Data: userinfoDto, Message: ""})
			return
		}

	}

	c.IndentedJSON(http.StatusNotFound, common.NoDataFailureResposne)

}
