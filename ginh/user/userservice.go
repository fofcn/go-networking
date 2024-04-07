package user

import (
	"errors"
	"fmt"
	"go-networking/log"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

var (
	userService *UserService
	once        sync.Once
)

func GetUserService(db *gorm.DB) *UserService {
	once.Do(func() {
		userService = &UserService{
			db: db,
		}
	})

	return userService
}

func (us *UserService) DoLogin(username string, password string) (string, error) {
	// 查询数据库
	var user User
	err := us.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Errorf("find user error, %v", username)
		return "", err
	}
	if user.ID == 0 {
		log.Errorf("user could not be found, %v", username)
		return "", fmt.Errorf("user not exist")
	}
	// 校验密码，算法BCRYPT
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Errorf("password is inccorrect, %v", username)
		return "", fmt.Errorf("password error")
	}

	// 生成token
	tokenService := &TokenService{}
	token, err := tokenService.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Errorf("token generation failed, %v", username)
		return "", err
	}

	return token, nil
}

func (us *UserService) DoRegister(cmd *RegisterCmd) error {
	// 查询数据库
	var user User
	err := us.db.Where("username = ?", cmd.Username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("find user error, %v %v", cmd.Username, err)
		return err
	}

	if user.ID != 0 {
		log.Errorf("use has already registerd, %v", cmd.Username)
		return UserHasRegisteredErr
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("password format is curious. %v", err)
		return PwdFmtIncorrect
	}
	return us.db.Create(&User{
		Username: cmd.Username,
		Password: string(hashedPassword),
	}).Error
}
