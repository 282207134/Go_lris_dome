package controllers

import (
	"fmt"
	"strconv"
	"time"

	"iris-cn-sample-project/models"
	"iris-cn-sample-project/utils"

	"github.com/kataras/iris/v12"
)

// APIDocs API 文档接口
func APIDocs(ctx iris.Context) {
	// 获取所有路由信息
	routes := ctx.App().GetRoutes()
	
	// 构建API文档
	apiDocs := make([]iris.Map, 0)
	for _, route := range routes {
		if len(route.Path) > 4 && route.Path[:4] == "/api" {
			apiDocs = append(apiDocs, iris.Map{
				"method": route.Method,
				"path":   route.Path,
				"name":   route.Name,
				"description": getRouteDescription(route.Path, route.Method),
			})
		}
	}

	// 组织完整的API文档
	documentation := iris.Map{
		"title":       "Iris Go 框架 API 文档",
		"version":     "v1.0.0",
		"description": "Iris Go 框架学习项目的 API 接口文档",
		"base_url":    "http://localhost:8080",
		"created_at":  time.Now().Format("2006-01-02 15:04:05"),
		"endpoints":   apiDocs,
		"categories": iris.Map{
			"认证相关": []string{
				"POST /api/auth/login",
				"POST /api/auth/register", 
				"POST /api/auth/refresh",
				"POST /api/auth/logout",
				"GET /api/auth/info",
				"POST /api/auth/change-password",
				"GET /api/auth/validate",
			},
			"用户管理": []string{
				"GET /api/users",
				"GET /api/users/{id}",
				"PUT /api/users/{id}",
				"DELETE /api/users/{id}",
			},
			"示例接口": []string{
				"GET /api/hello",
				"GET /api/data/{id}",
				"POST /api/form",
				"POST /api/upload",
			},
			"受保护接口": []string{
				"GET /api/protected/profile",
				"PUT /api/protected/profile",
			},
		},
		"authentication": iris.Map{
			"type": "Bearer Token",
			"description": "使用 JWT Bearer Token 进行身份验证",
			"header": "Authorization: Bearer <token>",
		},
		"response_format": iris.Map{
			"success": iris.Map{
				"code":    200,
				"message": "操作成功",
				"data":    "响应数据",
			},
			"error": iris.Map{
				"code":    "错误代码",
				"message": "错误信息",
				"errors":  "详细错误信息（可选）",
			},
		},
	}

	ctx.JSON(models.NewResponse(200, "API 文档", documentation))
}

// HealthCheck 健康检查接口
func HealthCheck(ctx iris.Context) {
	health := iris.Map{
		"status":    "healthy",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"uptime":    "0h 0m 0s", // 简化的运行时间
		"version":   "1.0.0",
		"service":   "iris-sample-project",
		"checks": iris.Map{
			"database": "ok",
			"memory":   "ok",
			"disk":     "ok",
		},
	}

	ctx.JSON(models.NewResponse(200, "健康检查", health))
}

// Metrics 系统指标接口
func Metrics(ctx iris.Context) {
	metrics := iris.Map{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"requests": iris.Map{
			"total":   1000,
			"success": 950,
			"error":   50,
		},
		"response_time": iris.Map{
			"avg":     "120ms",
			"p50":     "100ms",
			"p95":     "200ms",
			"p99":     "300ms",
		},
		"system": iris.Map{
			"cpu_usage":    "15.5%",
			"memory_usage": "45.2%",
			"disk_usage":   "22.8%",
			"goroutines":   25,
		},
	}

	ctx.JSON(models.NewResponse(200, "系统指标", metrics))
}

// Echo 回显接口（用于测试）
func Echo(ctx iris.Context) {
	// 获取请求信息
	requestInfo := iris.Map{
		"method":     ctx.Method(),
		"path":       ctx.Path(),
		"query":      ctx.URLParams(),
		"headers":    ctx.Request().Header,
		"remote_addr": ctx.RemoteAddr(),
		"user_agent": ctx.GetHeader("User-Agent"),
		"content_type": ctx.GetContentType(),
		"content_length": ctx.GetContentLength(),
	}

	// 尝试读取请求体
	if ctx.GetContentLength() > 0 {
		body := ctx.Request().Body
		if body != nil {
			// 这里简化处理，实际应该读取body内容
			requestInfo["body"] = "请求体内容（已省略）"
		}
	}

	// 添加时间戳
	requestInfo["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	ctx.JSON(models.NewResponse(200, "请求回显", requestInfo))
}

// Delay 延迟接口（用于测试超时）
func Delay(ctx iris.Context) {
	// 获取延迟参数（秒）
	secondsStr := ctx.URLParamDefault("seconds", "1")
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil || seconds < 0 || seconds > 30 {
		seconds = 1
	}

	// 模拟延迟
	time.Sleep(time.Duration(seconds) * time.Second)

	ctx.JSON(models.NewResponse(200, fmt.Sprintf("延迟 %d 秒后的响应", seconds), iris.Map{
		"delay_seconds": seconds,
		"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
	}))
}

// Headers 获取请求头信息
func Headers(ctx iris.Context) {
	headers := make(map[string]string)
	for name, values := range ctx.Request().Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	ctx.JSON(models.NewResponse(200, "请求头信息", iris.Map{
		"headers":   headers,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}))
}

// IP 获取客户端IP信息
func IP(ctx iris.Context) {
	clientIP := ctx.RemoteAddr()
	realIP := ctx.GetHeader("X-Real-IP")
	forwardedFor := ctx.GetHeader("X-Forwarded-For")

	ipInfo := iris.Map{
		"remote_addr":    clientIP,
		"real_ip":        realIP,
		"forwarded_for":  forwardedFor,
		"client_ip":      utils.GetClientIP(ctx.Request()),
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(models.NewResponse(200, "IP 信息", ipInfo))
}

// Cookies Cookie 操作接口
func Cookies(ctx iris.Context) {
	method := ctx.Method()

	switch method {
	case "GET":
		// 获取所有 Cookie
		cookies := ctx.Request().Cookies()
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie.Value
		}

		ctx.JSON(models.NewResponse(200, "Cookie 信息", iris.Map{
			"cookies":  cookieMap,
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		}))

	case "POST":
		// 设置 Cookie
		var cookieData struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}

		if err := ctx.ReadJSON(&cookieData); err != nil {
			ctx.JSON(models.NewResponse(400, "请求数据格式错误", nil))
			return
		}

		// 设置 Cookie
		ctx.SetCookieKV(cookieData.Name, cookieData.Value)

		ctx.JSON(models.NewResponse(200, "Cookie 设置成功", iris.Map{
			"name":     cookieData.Name,
			"value":    cookieData.Value,
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		}))

	default:
		ctx.JSON(models.NewResponse(405, "不支持的请求方法", nil))
	}
}

// Status 返回指定状态码
func Status(ctx iris.Context) {
	// 获取状态码参数
	codeStr := ctx.Params().Get("code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		code = 200
	}

	// 构建响应消息
	message := getStatusMessage(code)
	
	ctx.StatusCode(code)
	ctx.JSON(models.NewResponse(code, message, iris.Map{
		"status_code": code,
		"message":     message,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	}))
}

// 辅助函数

// getRouteDescription 获取路由描述
func getRouteDescription(path, method string) string {
	descriptions := map[string]string{
		"GET /api/hello":        "简单的问候接口，支持查询参数",
		"GET /api/data/{id}":    "获取指定ID的数据，路径参数示例",
		"POST /api/form":        "表单数据处理接口",
		"POST /api/upload":      "文件上传接口",
		"POST /api/auth/login":  "用户登录接口",
		"POST /api/auth/register": "用户注册接口",
		"GET /api/users":        "获取用户列表（需要认证）",
		"GET /api/users/{id}":   "获取指定用户信息（需要认证）",
		"PUT /api/users/{id}":   "更新用户信息（需要认证）",
		"DELETE /api/users/{id}": "删除用户（需要认证）",
	}

	key := method + " " + path
	if desc, exists := descriptions[key]; exists {
		return desc
	}

	return "API 接口"
}

// getStatusMessage 根据状态码获取消息
func getStatusMessage(code int) string {
	messages := map[int]string{
		200: "OK - 请求成功",
		201: "Created - 资源创建成功",
		400: "Bad Request - 请求参数错误",
		401: "Unauthorized - 未授权",
		403: "Forbidden - 禁止访问",
		404: "Not Found - 资源不存在",
		405: "Method Not Allowed - 请求方法不允许",
		500: "Internal Server Error - 服务器内部错误",
		502: "Bad Gateway - 网关错误",
		503: "Service Unavailable - 服务不可用",
	}

	if msg, exists := messages[code]; exists {
		return msg
	}

	return fmt.Sprintf("HTTP Status %d", code)
}