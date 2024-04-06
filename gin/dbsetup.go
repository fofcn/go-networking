package gin

import (
	"fmt"
	"go-networking/log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GDB *gorm.DB = nil

// DBConfig 用于存储数据库配置信息
type DBConfig struct {
	DbDriver   string
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
}

// LoadDBConfig 从环境变量或默认值加载数据库配置
func LoadDBConfig() *DBConfig {
	config := &DBConfig{
		DbDriver:   getEnv("DB_DRIVER", "mysql"), // 默认值为 "mysql"
		DbHost:     getEnv("DB_HOST", "localhost"),
		DbUser:     getEnv("DB_USER", "root"),
		DbPassword: getEnv("DB_PASSWORD", "password"),
		DbName:     getEnv("DB_NAME", "mydb"),
		DbPort:     getEnv("DB_PORT", "3306"),
	}
	return config
}

// getEnv 用于获取环境变量的值，如果未设置，则返回默认值
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// InitDB 初始化数据库连接
func InitDB() {
	log.Info("init DB")
	err := godotenv.Load(".dbenv")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := LoadDBConfig()

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       DBURL,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent SQL logging
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		// 这里可以添加更多的错误处理逻辑，比如重试连接
		os.Exit(1) // 显式退出程序
	}

	GDB = db
	// 这里添加AutoMigrate代码
	// GDB.AutoMigrate(&User{})
	log.Info("init DB completed")
}
