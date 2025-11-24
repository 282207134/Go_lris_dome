package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/kataras/iris/v12"
)

// Logger 自定义日志中间件
func Logger() iris.Handler {
	return func(ctx iris.Context) {
		// 记录开始时间
		start := time.Now()

		// 读取请求体（用于日志记录）
		var requestBody []byte
		if ctx.Request().Body != nil {
			requestBody, _ = io.ReadAll(ctx.Request().Body)
			ctx.Request().Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		ctx.Next()

		// 计算处理时间
		duration := time.Since(start)

		// 记录请求日志
		logRequest(ctx, requestBody, duration)
	}
}

// logRequest 记录请求日志
func logRequest(ctx iris.Context, requestBody []byte, duration time.Duration) {
	// 获取请求信息
	method := ctx.Method()
	path := ctx.Path()
	statusCode := ctx.GetStatusCode()
	clientIP := ctx.RemoteAddr()
	userAgent := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")

	// 获取响应大小
	responseSize := ctx.ResponseWriter().Written()

	// 构建日志消息
	logMessage := map[string]interface{}{
		"timestamp":    time.Now().Format("2006-01-02 15:04:05"),
		"method":       method,
		"path":         path,
		"status_code":  statusCode,
		"duration":     duration.String(),
		"client_ip":    clientIP,
		"user_agent":   userAgent,
		"referer":      referer,
		"response_size": responseSize,
	}

	// 添加请求体（仅对特定方法和路径）
	if shouldLogRequestBody(method, path) && len(requestBody) > 0 {
		logMessage["request_body"] = string(requestBody)
	}

	// 添加用户信息（如果有认证）
	if userID := ctx.Values().Get("user_id"); userID != nil {
		logMessage["user_id"] = userID
	}
	if username := ctx.Values().Get("username"); username != nil {
		logMessage["username"] = username
	}

	// 根据状态码确定日志级别
	logLevel := getLogLevel(statusCode)

	// 输出日志（这里简化处理，实际项目中可以使用 logrus、zap 等日志库）
	printLog(logLevel, logMessage)
}

// shouldLogRequestBody 判断是否应该记录请求体
func shouldLogRequestBody(method, path string) bool {
	// 不记录敏感接口的请求体
	sensitivePaths := []string{
		"/api/auth/login",
		"/api/auth/register",
		"/api/auth/refresh",
	}

	for _, sensitivePath := range sensitivePaths {
		if path == sensitivePath {
			return false
		}
	}

	// 只记录 POST、PUT、PATCH 请求的请求体
	return method == "POST" || method == "PUT" || method == "PATCH"
}

// getLogLevel 根据状态码获取日志级别
func getLogLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "ERROR"
	case statusCode >= 400:
		return "WARN"
	case statusCode >= 300:
		return "INFO"
	default:
		return "INFO"
	}
}

// printLog 打印日志（简化版本）
func printLog(level string, message map[string]interface{}) {
	// 这里简化处理，实际项目中应该使用专业的日志库
	switch level {
	case "ERROR":
		println("[ERROR]", formatLogMessage(message))
	case "WARN":
		println("[WARN]", formatLogMessage(message))
	default:
		println("[INFO]", formatLogMessage(message))
	}
}

// formatLogMessage 格式化日志消息
func formatLogMessage(message map[string]interface{}) string {
	// 简化的日志格式化
	return message["timestamp"].(string) + " " + 
		   message["method"].(string) + " " + 
		   message["path"].(string) + " " + 
		   string(rune(message["status_code"].(int))) + " " + 
		   message["duration"].(string)
}

// RequestID 请求ID中间件
func RequestID() iris.Handler {
	return func(ctx iris.Context) {
		// 生成或获取请求ID
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置请求ID到响应头和上下文
		ctx.Header("X-Request-ID", requestID)
		ctx.Values().Set("request_id", requestID)

		ctx.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简化的请求ID生成
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}