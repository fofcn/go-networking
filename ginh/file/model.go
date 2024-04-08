package file

import (
	"time"

	"gorm.io/gorm"
)

func InitFileMigrate(db *gorm.DB) {
	db.AutoMigrate(&FileInfo{})
}

type FileInfo struct {
	gorm.Model
	UserId        uint   `gorm:"index;not null"`
	ParentId      uint   `gorm:"index;not null"`
	OrgFilename   string `gorm:"size:255;not null"`
	Fileame       string `gorm:"size:255;not null"`
	FileType      uint8  `gorm:"not null"`
	ContentType   uint8  `gorm:"not null"`
	Deleted       uint8  `gorm:"not null"`
	FileCreatedAt time.Time
	FileUpdatedAt time.Time
}

func (f FileInfo) TableName() string {
	return "file_info"
}
