package services

import (
	"errors"
	"fmt"
	"time"

	"iris-cn-sample-project/database"
	"iris-cn-sample-project/models"
	"iris-cn-sample-project/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser 创建用户服务
func CreateUser(req *models.RegisterRequest) (*models.User, error) {
	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		Status:    "active",
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	return &user, nil
}

// LoginUser 用户登录服务
func LoginUser(username, password string) (*models.User, error) {
	db := database.GetDB()

	// 查找用户（支持用户名或邮箱登录）
	var user models.User
	if err := db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, errors.New("用户账户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(userID uint) (*models.UserInfo, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("id = ? AND status = ?", userID, "active").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 转换为用户信息（不包含敏感信息）
	userInfo := models.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userInfo, nil
}

// GetUsers 获取用户列表（分页）
func GetUsers(page, pageSize int) ([]*models.UserInfo, int64, error) {
	db := database.GetDB()

	var users []models.User
	var total int64

	// 获取总数
	if err := db.Model(&models.User{}).Where("status = ?", "active").Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %v", err)
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取用户列表
	if err := db.Where("status = ?", "active").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %v", err)
	}

	// 转换为用户信息列表
	userInfos := make([]*models.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return userInfos, total, nil
}

// UpdateUser 更新用户信息
func UpdateUser(userID uint, req *models.UpdateUserRequest) (*models.UserInfo, error) {
	db := database.GetDB()

	// 查找用户
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 更新字段
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	// 保存更新
	if err := db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("更新用户失败: %v", err)
	}

	// 转换为用户信息
	userInfo := models.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &userInfo, nil
}

// DeleteUser 删除用户（软删除）
func DeleteUser(userID uint) error {
	db := database.GetDB()

	// 软删除用户
	if err := db.Delete(&models.User{}, userID).Error; err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}

	return nil
}

// UpdateUserLastLogin 更新用户最后登录时间
func UpdateUserLastLogin(userID uint) error {
	db := database.GetDB()
	now := time.Now()

	if err := db.Model(&models.User{}).Where("id = ?", userID).Update("last_login", &now).Error; err != nil {
		return fmt.Errorf("更新最后登录时间失败: %v", err)
	}

	return nil
}

// ChangeUserPassword 修改用户密码
func ChangeUserPassword(userID uint, oldPassword, newPassword string) error {
	db := database.GetDB()

	// 查找用户
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	if err := db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	return nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(email string) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}

// SearchUsers 搜索用户
func SearchUsers(keyword string, page, pageSize int) ([]*models.UserInfo, int64, error) {
	db := database.GetDB()

	var users []models.User
	var total int64

	// 构建搜索条件
	searchPattern := "%" + keyword + "%"

	// 获取总数
	if err := db.Model(&models.User{}).
		Where("(username LIKE ? OR email LIKE ? OR first_name LIKE ? OR last_name LIKE ?) AND status = ?",
			searchPattern, searchPattern, searchPattern, searchPattern, "active").
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取搜索结果总数失败: %v", err)
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 搜索用户
	if err := db.Where("(username LIKE ? OR email LIKE ? OR first_name LIKE ? OR last_name LIKE ?) AND status = ?",
		searchPattern, searchPattern, searchPattern, searchPattern, "active").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("搜索用户失败: %v", err)
	}

	// 转换为用户信息列表
	userInfos := make([]*models.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return userInfos, total, nil
}

// ValidateUserCredentials 验证用户凭据
func ValidateUserCredentials(username, password string) (*models.User, error) {
	return LoginUser(username, password)
}

// IsUserExists 检查用户是否存在
func IsUserExists(username, email string) (bool, error) {
	db := database.GetDB()

	var count int64
	if err := db.Model(&models.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查用户是否存在失败: %v", err)
	}

	return count > 0, nil
}

// GetUserStats 获取用户统计信息
func GetUserStats() (map[string]interface{}, error) {
	db := database.GetDB()

	stats := make(map[string]interface{})

	// 总用户数
	var totalUsers int64
	if err := db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, fmt.Errorf("获取总用户数失败: %v", err)
	}
	stats["total_users"] = totalUsers

	// 活跃用户数
	var activeUsers int64
	if err := db.Model(&models.User{}).Where("status = ?", "active").Count(&activeUsers).Error; err != nil {
		return nil, fmt.Errorf("获取活跃用户数失败: %v", err)
	}
	stats["active_users"] = activeUsers

	// 管理员用户数
	var adminUsers int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminUsers).Error; err != nil {
		return nil, fmt.Errorf("获取管理员用户数失败: %v", err)
	}
	stats["admin_users"] = adminUsers

	// 今日注册用户数
	today := time.Now().Format("2006-01-02")
	var todayUsers int64
	if err := db.Model(&models.User{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayUsers).Error; err != nil {
		return nil, fmt.Errorf("获取今日注册用户数失败: %v", err)
	}
	stats["today_users"] = todayUsers

	return stats, nil
}