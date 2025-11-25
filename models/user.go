package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50" validate:"required,min=3,max=50"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100" validate:"required,email"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	FirstName string         `json:"first_name" gorm:"size:50"`
	LastName  string         `json:"last_name" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Role      string         `json:"role" gorm:"default:user;size:20"`
	Status    string         `json:"status" gorm:"default:active;size:20"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}

// AfterFind 查询后的钩子
func (u *User) AfterFind(tx *gorm.DB) error {
	// 可以在这里添加查询后的逻辑
	return nil
}

// GetFullName 获取用户全名
func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.Username
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}