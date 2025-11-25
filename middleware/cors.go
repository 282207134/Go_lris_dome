package middleware

import (
	"github.com/kataras/iris/v12"
)

// CORS 跨域资源共享中间件
func CORS() iris.Handler {
	return func(ctx iris.Context) {
		// 设置 CORS 头
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if ctx.Method() == "OPTIONS" {
			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// CORSWithOrigin 允许指定源的 CORS 中间件
func CORSWithOrigin(allowedOrigins ...string) iris.Handler {
	return func(ctx iris.Context) {
		origin := ctx.GetHeader("Origin")
		
		// 检查请求源是否在允许列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		// 如果允许则设置对应的源，否则不设置
		if allowed {
			ctx.Header("Access-Control-Allow-Origin", origin)
		} else {
			ctx.Header("Access-Control-Allow-Origin", "")
		}

		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if ctx.Method() == "OPTIONS" {
			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// CORSWithCredentials 带凭据的 CORS 中间件
func CORSWithCredentials(allowedOrigins ...string) iris.Handler {
	return func(ctx iris.Context) {
		origin := ctx.GetHeader("Origin")
		
		// 检查请求源是否在允许列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		// 如果允许则设置对应的源，否则不设置
		if allowed {
			ctx.Header("Access-Control-Allow-Origin", origin)
		} else {
			ctx.Header("Access-Control-Allow-Origin", "")
		}

		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Vary", "Origin")

		// 处理预检请求
		if ctx.Method() == "OPTIONS" {
			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	}
}