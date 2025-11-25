package database

import (
	"fmt"
	"log"

	"iris-cn-sample-project/config"
	"iris-cn-sample-project/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.GetConfig()
	
	// 配置数据库连接
	var err error
	logLevel := logger.Silent
	if cfg.Log.Level == "debug" {
		logLevel = logger.Info
	}

	switch cfg.Database.Driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(cfg.Database.Database), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
	default:
		return fmt.Errorf("不支持的数据库驱动: %s", cfg.Database.Driver)
	}

	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 获取底层的 sql.DB 对象进行连接池配置
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 自动迁移数据库表结构
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	// 初始化基础数据
	if err := seedData(); err != nil {
		log.Printf("初始化基础数据失败: %v", err)
	}

	log.Println("数据库初始化成功")
	return nil
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
	)
}

// seedData 初始化基础数据
func seedData() error {
	// 检查是否已有用户数据
	var count int64
	if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	// 如果没有数据，创建示例用户
	if count == 0 {
		users := []models.User{
			{
				Username: "admin",
				Email:    "admin@example.com",
				Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.IY4Jyy5.R9WOM4.7O1HqK4C9/Q9/pW", // password: admin123
				Role:     "admin",
				Status:   "active",
			},
			{
				Username: "user",
				Email:    "user@example.com",
				Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.IY4Jyy5.R9WOM4.7O1HqK4C9/Q9/pW", // password: user123
				Role:     "user",
				Status:   "active",
			},
		}

		for _, user := range users {
			if err := DB.Create(&user).Error; err != nil {
				return fmt.Errorf("创建示例用户失败: %v", err)
			}
		}

		log.Println("示例用户数据创建成功")
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}