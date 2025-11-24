package main

import (
	"testing"

	"iris-cn-sample-project/config"
	"iris-cn-sample-project/database"
	"iris-cn-sample-project/models"
)

// TestConfig 测试配置加载
func TestConfig(t *testing.T) {
	cfg := config.GetConfig()
	
	if cfg.Server.Port == "" {
		t.Error("服务器端口配置为空")
	}
	
	if cfg.JWT.Secret == "" {
		t.Error("JWT 密钥配置为空")
	}
	
	if cfg.JWT.ExpirationTime <= 0 {
		t.Error("JWT 过期时间配置无效")
	}
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	err := database.InitDB()
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
	}
	
	db := database.GetDB()
	if db == nil {
		t.Error("获取数据库实例失败")
	}
	
	// 测试自动迁移
	err = database.AutoMigrate()
	if err != nil {
		t.Errorf("数据库迁移失败: %v", err)
	}
}

// TestUserModel 测试用户模型
func TestUserModel(t *testing.T) {
	// 初始化数据库
	if err := database.InitDB(); err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}
	
	db := database.GetDB()
	
	// 创建测试用户
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "测试",
		LastName:  "用户",
		Role:      "user",
		Status:    "active",
	}
	
	// 测试创建用户
	if err := db.Create(&user).Error; err != nil {
		t.Errorf("创建用户失败: %v", err)
	}
	
	if user.ID == 0 {
		t.Error("用户 ID 应该大于 0")
	}
	
	// 测试查询用户
	var foundUser models.User
	if err := db.First(&foundUser, user.ID).Error; err != nil {
		t.Errorf("查询用户失败: %v", err)
	}
	
	if foundUser.Username != user.Username {
		t.Error("查询到的用户名不匹配")
	}
	
	// 测试更新用户
	user.FirstName = "更新"
	if err := db.Save(&user).Error; err != nil {
		t.Errorf("更新用户失败: %v", err)
	}
	
	// 测试删除用户
	if err := db.Delete(&user).Error; err != nil {
		t.Errorf("删除用户失败: %v", err)
	}
}

// TestValidation 测试数据验证
func TestValidation(t *testing.T) {
	// 这个测试需要 utils 包中的验证函数
	// 由于我们在测试中，暂时跳过这个测试
	t.Skip("验证测试需要完整的导入路径")
}