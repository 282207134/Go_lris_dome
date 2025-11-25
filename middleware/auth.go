package middleware

import (
	"strings"

	"iris-cn-sample-project/config"
	"iris-cn-sample-project/utils"

	"github.com/kataras/iris/v12"
)

// JWTAuthentication JWT 认证中间件
func JWTAuthentication() iris.Handler {
	return func(ctx iris.Context) {
		// 从请求头获取 Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(iris.Map{
				"code":    401,
				"message": "缺少认证令牌",
			})
			ctx.StatusCode(iris.StatusUnauthorized)
			return
		}

		// 检查 Bearer 前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			ctx.JSON(iris.Map{
				"code":    401,
				"message": "认证令牌格式错误",
			})
			ctx.StatusCode(iris.StatusUnauthorized)
			return
		}

		// 提取令牌
		tokenString := authHeader[len(bearerPrefix):]
		
		// 验证令牌
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			ctx.JSON(iris.Map{
				"code":    401,
				"message": "无效的认证令牌: " + err.Error(),
			})
			ctx.StatusCode(iris.StatusUnauthorized)
			return
		}

		// 将用户信息存储到上下文中
		ctx.Values().Set("user_id", claims.UserID)
		ctx.Values().Set("username", claims.Username)
		ctx.Values().Set("role", claims.Role)

		// 继续处理请求
		ctx.Next()
	}
}

// RequireRole 角色验证中间件
func RequireRole(roles ...string) iris.Handler {
	return func(ctx iris.Context) {
		// 获取用户角色
		userRole, ok := ctx.Values().GetString("role")
		if !ok {
			ctx.JSON(iris.Map{
				"code":    403,
				"message": "无法获取用户角色信息",
			})
			ctx.StatusCode(iris.StatusForbidden)
			return
		}

		// 检查用户角色是否在允许的角色列表中
		hasPermission := false
		for _, role := range roles {
			if userRole == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			ctx.JSON(iris.Map{
				"code":    403,
				"message": "权限不足",
			})
			ctx.StatusCode(iris.StatusForbidden)
			return
		}

		ctx.Next()
	}
}

// RequireAdmin 管理员权限中间件
func RequireAdmin() iris.Handler {
	return RequireRole("admin")
}

// OptionalAuthentication 可选认证中间件
func OptionalAuthentication() iris.Handler {
	return func(ctx iris.Context) {
		// 从请求头获取 Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			// 没有认证头，继续处理请求
			ctx.Next()
			return
		}

		// 检查 Bearer 前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			// 格式错误，继续处理请求
			ctx.Next()
			return
		}

		// 提取令牌
		tokenString := authHeader[len(bearerPrefix):]
		
		// 验证令牌
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			// 令牌无效，继续处理请求
			ctx.Next()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Values().Set("user_id", claims.UserID)
		ctx.Values().Set("username", claims.Username)
		ctx.Values().Set("role", claims.Role)

		// 继续处理请求
		ctx.Next()
	}
}