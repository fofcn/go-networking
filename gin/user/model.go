package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        int    `gorm:"primaryKey;autoIncrement`
	Username  string `gorm:"unique;size:100;not null"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
