package utils

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

// ResponseUtil 响应工具结构体
type ResponseUtil struct {
	ctx iris.Context
}

// NewResponseUtil 创建响应工具实例
func NewResponseUtil(ctx iris.Context) *ResponseUtil {
	return &ResponseUtil{ctx: ctx}
}

// Success 返回成功响应
func (r *ResponseUtil) Success(data interface{}) {
	response := map[string]interface{}{
		"code":    200,
		"message": "操作成功",
		"data":    data,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.JSON(response)
}

// SuccessWithMessage 返回带消息的成功响应
func (r *ResponseUtil) SuccessWithMessage(message string, data interface{}) {
	response := map[string]interface{}{
		"code":      200,
		"message":   message,
		"data":      data,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.JSON(response)
}

// Error 返回错误响应
func (r *ResponseUtil) Error(code int, message string) {
	response := map[string]interface{}{
		"code":      code,
		"message":   message,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(code)
	r.ctx.JSON(response)
}

// ErrorWithData 返回带数据的错误响应
func (r *ResponseUtil) ErrorWithData(code int, message string, data interface{}) {
	response := map[string]interface{}{
		"code":      code,
		"message":   message,
		"data":      data,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(code)
	r.ctx.JSON(response)
}

// ValidationError 返回验证错误响应
func (r *ResponseUtil) ValidationError(errors interface{}) {
	response := map[string]interface{}{
		"code":      400,
		"message":   "输入数据验证失败",
		"errors":    errors,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(iris.StatusBadRequest)
	r.ctx.JSON(response)
}

// PageResponse 返回分页响应
func (r *ResponseUtil) PageResponse(data interface{}, page, pageSize int, total int64) {
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	response := map[string]interface{}{
		"code":    200,
		"message": "获取数据成功",
		"data":    data,
		"page": map[string]interface{}{
			"current":   page,
			"page_size": pageSize,
			"total":     total,
			"total_page": totalPage,
		},
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.JSON(response)
}

// Created 返回创建成功响应
func (r *ResponseUtil) Created(data interface{}) {
	response := map[string]interface{}{
		"code":      201,
		"message":   "创建成功",
		"data":      data,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(iris.StatusCreated)
	r.ctx.JSON(response)
}

// Updated 返回更新成功响应
func (r *ResponseUtil) Updated(data interface{}) {
	response := map[string]interface{}{
		"code":      200,
		"message":   "更新成功",
		"data":      data,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.JSON(response)
}

// Deleted 返回删除成功响应
func (r *ResponseUtil) Deleted() {
	response := map[string]interface{}{
		"code":      200,
		"message":   "删除成功",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.JSON(response)
}

// Unauthorized 返回未授权响应
func (r *ResponseUtil) Unauthorized(message string) {
	if message == "" {
		message = "未授权访问"
	}
	response := map[string]interface{}{
		"code":      401,
		"message":   message,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(iris.StatusUnauthorized)
	r.ctx.JSON(response)
}

// Forbidden 返回禁止访问响应
func (r *ResponseUtil) Forbidden(message string) {
	if message == "" {
		message = "权限不足"
	}
	response := map[string]interface{}{
		"code":      403,
		"message":   message,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(iris.StatusForbidden)
	r.ctx.JSON(response)
}

// NotFound 返回未找到响应
func (r *ResponseUtil) NotFound(message string) {
	if message == "" {
		message = "资源不存在"
	}
	response := map[string]interface{}{
		"code":      404,
		"message":   message,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(iris.StatusNotFound)
	r.ctx.JSON(response)
}

// InternalServerError 返回服务器内部错误响应
func (r *ResponseUtil) InternalServerError(message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	response := map[string]interface{}{
		"code":      500,
		"message":   message,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	r.ctx.StatusCode(iris.StatusInternalServerError)
	r.ctx.JSON(response)
}

// GetClientIP 获取客户端真实IP地址
func GetClientIP(r *http.Request) string {
	// 尝试从 X-Forwarded-For 头获取
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" && ip != "unknown" {
				return ip
			}
		}
	}

	// 尝试从 X-Real-IP 头获取
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		ip := strings.TrimSpace(xRealIP)
		if ip != "" && ip != "unknown" {
			return ip
		}
	}

	// 尝试从 X-Forwarded 头获取
	xForwarded := r.Header.Get("X-Forwarded")
	if xForwarded != "" {
		ip := strings.TrimSpace(xForwarded)
		if ip != "" && ip != "unknown" {
			return ip
		}
	}

	// 尝试从 Forwarded-For 头获取
	forwardedFor := r.Header.Get("Forwarded-For")
	if forwardedFor != "" {
		ip := strings.TrimSpace(forwardedFor)
		if ip != "" && ip != "unknown" {
			return ip
		}
	}

	// 最后使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// IsAJAXRequest 检查是否为AJAX请求
func IsAJAXRequest(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// IsAPIRequest 检查是否为API请求
func IsAPIRequest(path string) bool {
	return strings.HasPrefix(path, "/api/")
}

// GetRequestInfo 获取请求信息
func GetRequestInfo(r *http.Request) map[string]interface{} {
	info := make(map[string]interface{})
	
	info["method"] = r.Method
	info["path"] = r.URL.Path
	info["query"] = r.URL.RawQuery
	info["proto"] = r.Proto
	info["host"] = r.Host
	info["remote_addr"] = r.RemoteAddr
	info["user_agent"] = r.Header.Get("User-Agent")
	info["referer"] = r.Header.Get("Referer")
	info["content_type"] = r.Header.Get("Content-Type")
	info["content_length"] = r.ContentLength
	info["client_ip"] = GetClientIP(r)
	info["is_ajax"] = IsAJAXRequest(r)
	info["is_api"] = IsAPIRequest(r.URL.Path)
	
	return info
}

// PrettyJSON 格式化JSON输出
func PrettyJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToJSON 将数据转换为JSON字符串
func ToJSON(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串解析数据
func FromJSON(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// SanitizeString 清理字符串
func SanitizeString(s string) string {
	// 简单的字符串清理，去除首尾空白
	return strings.TrimSpace(s)
}

// SanitizeSQL 防止SQL注入（简单版本）
func SanitizeSQL(s string) string {
	// 简单的SQL注入防护，实际项目中应该使用参数化查询
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, ";", "")
	s = strings.ReplaceAll(s, "--", "")
	s = strings.ReplaceAll(s, "/*", "")
	s = strings.ReplaceAll(s, "*/", "")
	return s
}

// GenerateTimestamp 生成时间戳
func GenerateTimestamp() int64 {
	return time.Now().Unix()
}

// GenerateDateTime 生成日期时间字符串
func GenerateDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GenerateDate 生成日期字符串
func GenerateDate() string {
	return time.Now().Format("2006-01-02")
}

// GenerateTime 生成时间字符串
func GenerateTime() string {
	return time.Now().Format("15:04:05")
}

// ParseDateTime 解析日期时间字符串
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", dateTimeStr)
}

// ParseDate 解析日期字符串
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04:05", timeStr)
}

// FormatDateTime 格式化日期时间
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate 格式化日期
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// TimeAgo 计算时间差（多久之前）
func TimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return "刚刚"
	} else if duration < time.Hour {
		return fmt.Sprintf("%d分钟前", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d小时前", int(duration.Hours()))
	} else if duration < 30*24*time.Hour {
		return fmt.Sprintf("%d天前", int(duration.Hours()/24))
	} else if duration < 12*30*24*time.Hour {
		return fmt.Sprintf("%d个月前", int(duration.Hours()/24/30))
	} else {
		return fmt.Sprintf("%d年前", int(duration.Hours()/24/365))
	}
}

// Contains 检查字符串切片是否包含指定字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates 去除字符串切片中的重复项
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// InStringArray 检查字符串是否在数组中
func InStringArray(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

// InIntArray 检查整数是否在数组中
func InIntArray(num int, arr []int) bool {
	for _, item := range arr {
		if item == num {
			return true
		}
	}
	return false
}

// MergeMaps 合并两个map
func MergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	for k, v := range map1 {
		result[k] = v
	}
	
	for k, v := range map2 {
		result[k] = v
	}
	
	return result
}

// GetMapValue 安全获取map值
func GetMapValue(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// SetMapValue 设置map值
func SetMapValue(m map[string]interface{}, key string, value interface{}) {
	if m == nil {
		m = make(map[string]interface{})
	}
	m[key] = value
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	
	return string(result)
}

// GenerateUUID 生成简单的UUID（简化版本）
func GenerateUUID() string {
	return GenerateRandomString(32)
}

// Pagination 计算分页信息
func Pagination(page, pageSize int, total int64) (offset int, totalPages int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	
	offset = (page - 1) * pageSize
	totalPages = int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return offset, totalPages
}