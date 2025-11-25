package models

import (
    "time"
)

// Response 通用响应结构体
type Response struct {
    Code    int         `json:"code"`    // 响应状态码
    Message string      `json:"message"` // 响应消息
    Data    interface{} `json:"data"`    // 响应数据
}

// PageResponse 分页响应结构体
type PageResponse struct {
    Code    int         `json:"code"`    // 响应状态码
    Message string      `json:"message"` // 响应消息
    Data    interface{} `json:"data"`    // 响应数据
    Page    PageInfo    `json:"page"`    // 分页信息
}

// PageInfo 分页信息
type PageInfo struct {
    Current   int   `json:"current"`   // 当前页码
    PageSize  int   `json:"page_size"` // 每页大小
    Total     int64 `json:"total"`     // 总记录数
    TotalPage int   `json:"total_page"` // 总页数
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    User      *UserInfo `json:"user"`
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
    Username  string `json:"username" validate:"required,min=3,max=50"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=6"`
    FirstName string `json:"first_name" validate:"max=50"`
    LastName  string `json:"last_name" validate:"max=50"`
}

// UserInfo 用户信息结构体（不包含敏感信息）
type UserInfo struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Avatar    string    `json:"avatar"`
    Role      string    `json:"role"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
    FirstName string `json:"first_name" validate:"max=50"`
    LastName  string `json:"last_name" validate:"max=50"`
    Avatar    string `json:"avatar"`
    Role      string `json:"role" validate:"omitempty,oneof=admin user"`
    Status    string `json:"status" validate:"omitempty,oneof=active inactive"`
}

// ChangePasswordRequest 修改密码请求结构体
type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=6"`
}

// FileUploadResponse 文件上传响应结构体
type FileUploadResponse struct {
    FileName string `json:"file_name"`
    FileSize int64  `json:"file_size"`
    FileURL  string `json:"file_url"`
    Success  bool   `json:"success"`
    Message  string `json:"message"`
}

// ErrorResponse 错误响应结构体
type ErrorResponse struct {
    Code      int                    `json:"code"`
    Message   string                 `json:"message"`
    Errors    map[string]interface{} `json:"errors,omitempty"`
    Timestamp string                 `json:"timestamp"`
    Path      string                 `json:"path"`
}

// SuccessResponse 成功响应结构体
type SuccessResponse struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp string      `json:"timestamp"`
}

// NewResponse 创建通用响应
func NewResponse(code int, message string, data interface{}) *Response {
    return &Response{
        Code:    code,
        Message: message,
        Data:    data,
    }
}

// NewPageResponse 创建分页响应
func NewPageResponse(code int, message string, data interface{}, current, pageSize int, total int64) *PageResponse {
    totalPage := int(total) / pageSize
    if int(total)%pageSize > 0 {
        totalPage++
    }

    return &PageResponse{
        Code:    code,
        Message: message,
        Data:    data,
        Page: PageInfo{
            Current:   current,
            PageSize:  pageSize,
            Total:     total,
            TotalPage: totalPage,
        },
    }
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, errors map[string]interface{}, path string) *ErrorResponse {
    return &ErrorResponse{
        Code:      code,
        Message:   message,
        Errors:    errors,
        Timestamp: "2023-01-01 00:00:00", // 简化的时间戳
        Path:      path,
    }
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(code int, message string, data interface{}) *SuccessResponse {
    return &SuccessResponse{
        Code:      code,
        Message:   message,
        Data:      data,
        Timestamp: "2023-01-01 00:00:00", // 简化的时间戳
    }
}