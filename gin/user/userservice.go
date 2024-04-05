package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
}

func (us *UserService) Login(username string, password string) (string, error) {
	// 查询数据库
	var user User
	err := GDB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	if user.ID == 0 {
		return "", fmt.Errorf("user not exist")
	}
	// 校验密码，算法BCRYPT
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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
