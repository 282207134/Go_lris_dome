package controllers

import (
    "time"

    "iris-cn-sample-project/models"
    "iris-cn-sample-project/services"
    "iris-cn-sample-project/utils"

    "github.com/kataras/iris/v12"
)

// Login 用户登录接口
func Login(ctx iris.Context) {
    // 解析登录请求数据
    var loginReq models.LoginRequest
    if err := ctx.ReadJSON(&loginReq); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "请求数据格式错误: " + err.Error(),
        })
        return
    }

    // 验证输入数据
    if err := utils.ValidateStruct(&loginReq); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "输入数据验证失败",
            "errors":  err,
        })
        return
    }

    // 调用服务层进行登录验证
    user, err := services.LoginUser(loginReq.Username, loginReq.Password)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "用户名或密码错误",
        })
        return
    }

    // 生成 JWT 令牌
    token, expiresAt, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    500,
            "message": "生成令牌失败: " + err.Error(),
        })
        return
    }

    // 更新用户最后登录时间
    services.UpdateUserLastLogin(user.ID)

    // 返回登录成功响应
    loginResp := models.LoginResponse{
        Token:     token,
        ExpiresAt: expiresAt,
        User: &models.UserInfo{
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
        },
    }

    ctx.JSON(models.NewResponse(200, "登录成功", loginResp))
}

// Register 用户注册接口
func Register(ctx iris.Context) {
    // 解析注册请求数据
    var registerReq models.RegisterRequest
    if err := ctx.ReadJSON(&registerReq); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "请求数据格式错误: " + err.Error(),
        })
        return
    }

    // 验证输入数据
    if err := utils.ValidateStruct(&registerReq); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "输入数据验证失败",
            "errors":  err,
        })
        return
    }

    // 调用服务层创建用户
    user, err := services.CreateUser(&registerReq)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "用户注册失败: " + err.Error(),
        })
        return
    }

    // 返回注册成功响应（不包含敏感信息）
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

    ctx.JSON(models.NewResponse(201, "用户注册成功", userInfo))
}

// RefreshToken 刷新令牌接口
func RefreshToken(ctx iris.Context) {
    // 从请求头获取当前令牌
    authHeader := ctx.GetHeader("Authorization")
    if authHeader == "" {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "缺少认证令牌",
        })
        return
    }

    // 提取令牌
    const bearerPrefix = "Bearer "
    if len(authHeader) <= len(bearerPrefix) {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "认证令牌格式错误",
        })
        return
    }

    tokenString := authHeader[len(bearerPrefix):]

    // 验证当前令牌（即使过期也要能解析出用户信息）
    claims, err := utils.ParseJWTWithoutValidation(tokenString)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "无效的认证令牌",
        })
        return
    }

    // 检查用户是否仍然存在且有效
    user, err := services.GetUserByID(claims.UserID)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "用户不存在或已被禁用",
        })
        return
    }

    if !user.IsActive() {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "用户账户已被禁用",
        })
        return
    }

    // 生成新的令牌
    newToken, expiresAt, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    500,
            "message": "生成新令牌失败: " + err.Error(),
        })
        return
    }

    // 返回新令牌
    ctx.JSON(models.NewResponse(200, "令牌刷新成功", iris.Map{
        "token":      newToken,
        "expires_at": expiresAt,
        "type":       "Bearer",
    }))
}

// Logout 用户登出接口
func Logout(ctx iris.Context) {
    // 在实际应用中，这里可以将令牌加入黑名单
    // 由于 JWT 是无状态的，真正的登出需要在客户端删除令牌
    
    ctx.JSON(models.NewResponse(200, "登出成功", iris.Map{
        "message": "请在客户端删除存储的令牌",
    }))
}

// ChangePassword 修改密码接口
func ChangePassword(ctx iris.Context) {
    // 获取当前用户ID
    userID := ctx.Values().GetUintDefault("user_id", 0)
    if userID == 0 {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "无效的用户信息",
        })
        return
    }

    // 解析修改密码请求数据
    var changePwdReq models.ChangePasswordRequest
    if err := ctx.ReadJSON(&changePwdReq); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "请求数据格式错误: " + err.Error(),
        })
        return
    }

    // 验证输入数据
    if err := utils.ValidateStruct(&changePwdReq); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "输入数据验证失败",
            "errors":  err,
        })
        return
    }

    // 调用服务层修改密码
    if err := services.ChangeUserPassword(userID, changePwdReq.OldPassword, changePwdReq.NewPassword); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "密码修改失败: " + err.Error(),
        })
        return
    }

    ctx.JSON(models.NewResponse(200, "密码修改成功", nil))
}

// GetAuthInfo 获取当前认证用户信息
func GetAuthInfo(ctx iris.Context) {
    // 从中间件获取用户信息
    userID := ctx.Values().GetUintDefault("user_id", 0)
    username := ctx.Values().GetStringDefault("username", "")
    role := ctx.Values().GetStringDefault("role", "")

    // 获取完整的用户信息
    user, err := services.GetUserByID(userID)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    404,
            "message": "用户信息不存在",
        })
        return
    }

    // 返回用户信息
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

    ctx.JSON(models.NewResponse(200, "获取认证信息成功", iris.Map{
        "user": userInfo,
        "token_info": iris.Map{
            "user_id":  userID,
            "username": username,
            "role":     role,
            "issued_at": time.Now().Format("2006-01-02 15:04:05"),
        },
    }))
}

// ValidateToken 验证令牌接口
func ValidateToken(ctx iris.Context) {
    // 从请求头获取令牌
    authHeader := ctx.GetHeader("Authorization")
    if authHeader == "" {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "缺少认证令牌",
        })
        return
    }

    // 提取令牌
    const bearerPrefix = "Bearer "
    if len(authHeader) <= len(bearerPrefix) {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "认证令牌格式错误",
        })
        return
    }

    tokenString := authHeader[len(bearerPrefix):]

    // 验证令牌
    claims, err := utils.ValidateJWT(tokenString)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "令牌验证失败: " + err.Error(),
        })
        return
    }

    // 检查用户是否仍然有效
    user, err := services.GetUserByID(claims.UserID)
    if err != nil || !user.IsActive() {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "用户不存在或已被禁用",
        })
        return
    }

    ctx.JSON(models.NewResponse(200, "令牌验证成功", iris.Map{
        "valid":      true,
        "user_id":    claims.UserID,
        "username":   claims.Username,
        "role":       claims.Role,
        "expires_at": claims.ExpiresAt,
    }))
}