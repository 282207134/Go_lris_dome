package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/kataras/iris/v12"
)

// Recovery 恢复中间件，用于捕获 panic 并优雅处理
func Recovery() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 信息
				logPanic(ctx, err)
				
				// 返回友好的错误响应
				handlePanic(ctx, err)
			}
		}()

		ctx.Next()
	}
}

// logPanic 记录 panic 信息
func logPanic(ctx iris.Context, err interface{}) {
	// 构建错误信息
	errorInfo := map[string]interface{}{
		"error":      fmt.Sprintf("%v", err),
		"stack":      string(debug.Stack()),
		"method":     ctx.Method(),
		"path":       ctx.Path(),
		"query":      ctx.URLParams(),
		"headers":    ctx.Request().Header,
		"remote_addr": ctx.RemoteAddr(),
		"user_agent": ctx.GetHeader("User-Agent"),
	}

	// 添加用户信息（如果有认证）
	if userID := ctx.Values().Get("user_id"); userID != nil {
		errorInfo["user_id"] = userID
	}
	if username := ctx.Values().Get("username"); username != nil {
		errorInfo["username"] = username
	}

	// 输出错误日志（这里简化处理）
	println("[PANIC] 服务器内部错误:")
	println("错误信息:", fmt.Sprintf("%v", err))
	println("请求方法:", ctx.Method())
	println("请求路径:", ctx.Path())
	println("客户端IP:", ctx.RemoteAddr())
	println("堆栈信息:", string(debug.Stack()))
}

// handlePanic 处理 panic 响应
func handlePanic(ctx iris.Context, err interface{}) {
	// 根据请求类型返回不同的响应格式
	accept := ctx.GetHeader("Accept")
	
	// 如果是 API 请求或接受 JSON
	if isAPIRequest(ctx.Path()) || containsJSON(accept) {
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "服务器内部错误，请稍后重试",
			"error":   "internal_server_error",
		})
	} else {
		// 返回 HTML 错误页面
		ctx.HTML(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>服务器错误</title>
				<meta charset="UTF-8">
				<style>
					body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
					.error { color: #e74c3c; font-size: 48px; margin-bottom: 20px; }
					.message { color: #333; font-size: 18px; }
				</style>
			</head>
			<body>
				<div class="error">500</div>
				<div class="message">服务器内部错误，请稍后重试</div>
			</body>
			</html>
		`)
	}

	ctx.StatusCode(iris.StatusInternalServerError)
}

// isAPIRequest 判断是否为 API 请求
func isAPIRequest(path string) bool {
	return len(path) >= 4 && path[:4] == "/api"
}

// containsJSON 检查是否包含 JSON
func containsJSON(accept string) bool {
	return len(accept) >= 4 && 
		   (accept[:4] == "appl" || 
		    accept[:4] == "text" ||
		    accept[:4] == "*/")
}

// CustomRecovery 自定义恢复中间件
func CustomRecovery(customHandler func(ctx iris.Context, err interface{})) iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 信息
				logPanic(ctx, err)
				
				// 使用自定义处理器
				if customHandler != nil {
					customHandler(ctx, err)
				} else {
					// 使用默认处理器
					handlePanic(ctx, err)
				}
			}
		}()

		ctx.Next()
	}
}

// RecoveryWithLogger 带自定义日志记录器的恢复中间件
func RecoveryWithLogger(logger func(ctx iris.Context, err interface{})) iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 使用自定义日志记录器
				if logger != nil {
					logger(ctx, err)
				} else {
					logPanic(ctx, err)
				}
				
				// 处理响应
				handlePanic(ctx, err)
			}
		}()

		ctx.Next()
	}
}