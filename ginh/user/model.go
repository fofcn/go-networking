package user

import "gorm.io/gorm"

func InitAutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;size:100;not null"`
	Password string `gorm:"size:255;not null"`
}

// TableName 会设置 GORM 对这个结构体使用的表名。
func (User) TableName() string {
	return "nas_user"
}
