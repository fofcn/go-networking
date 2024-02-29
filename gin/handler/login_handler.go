package handler

import (
	"go-networking/gin/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Loginrequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var embeddedCredentials = []Loginrequest{
	{Username: "hello", Password: "world"},
}

func Login(c *gin.Context) {
	var request Loginrequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusOK, common.SuccessResponse{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}

	if request.Username == embeddedCredentials[0].Username && request.Password == embeddedCredentials[0].Password {
		c.IndentedJSON(http.StatusOK, common.NoDataSuccessResposne)
	} else {
		c.IndentedJSON(http.StatusOK, common.SuccessResponse{Data: "Failed", Message: "Username or password is incorrect."})
		return
	}
}
