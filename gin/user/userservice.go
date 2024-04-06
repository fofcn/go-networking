package user

import (
	"fmt"
	"go-networking/gin"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
}

func (us *UserService) DoLogin(username string, password string) (string, error) {
	// 查询数据库
	var user User
	err := gin.GDB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	if user.ID == 0 {
		return "", fmt.Errorf("user not exist")
	}
	// 校验密码，算法BCRYPT
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("password error")
	}

	// 生成token
	tokenService := &TokenService{}
	token, err := tokenService.GenerateToken(user.Id, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *UserService) DoRegister(username string, password string) error {
	// 查询数据库
	var user User
	err := gin.GDB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}
	if user.ID != 0 {
		return fmt.Errorf("user already exist")
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return gin.GDB.Create(&User{
		Username: username,
		Password: string(hashedPassword),
	}).Error
}
