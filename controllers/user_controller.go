package controllers

import (
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "iris-cn-sample-project/database"
    "iris-cn-sample-project/middleware"
    "iris-cn-sample-project/models"
    "iris-cn-sample-project/services"

    "github.com/kataras/iris/v12"
)

// Index 首页控制器
func Index(ctx iris.Context) {
    ctx.ViewData("title", "Iris Go 框架学习示例")
    ctx.ViewData("message", "欢迎使用 Iris Go 框架！")
    ctx.View("index.html")
}

// Home 主页控制器
func Home(ctx iris.Context) {
    ctx.JSON(iris.Map{
        "code":    200,
        "message": "欢迎来到 Iris Go 框架学习项目",
        "data": iris.Map{
            "framework": "Iris",
            "version":   "v12",
            "features": []string{
                "高性能路由",
                "中间件支持",
                "模板引擎",
                "静态文件服务",
                "WebSocket 支持",
                "JSON API",
            },
        },
    })
}

// Hello 简单的问候接口
func Hello(ctx iris.Context) {
    // 获取查询参数
    name := ctx.URLParamDefault("name", "访客")
    
    ctx.JSON(iris.Map{
        "code":    200,
        "message": fmt.Sprintf("你好，%s！", name),
        "data": iris.Map{
            "method": ctx.Method(),
            "path":   ctx.Path(),
            "query":  ctx.URLParams(),
        },
    })
}

// GetData 获取数据接口（路径参数示例）
func GetData(ctx iris.Context) {
    // 获取路径参数
    id, err := ctx.Params().GetInt("id")
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "无效的ID参数",
        })
        return
    }

    // 模拟数据
    data := iris.Map{
        "id":       id,
        "title":    fmt.Sprintf("数据项 %d", id),
        "content":  fmt.Sprintf("这是第 %d 条数据的详细内容", id),
        "status":   "active",
        "created":  "2023-01-01 10:00:00",
    }

    ctx.JSON(iris.Map{
        "code":    200,
        "message": "获取数据成功",
        "data":    data,
    })
}

// HandleForm 表单处理接口
func HandleForm(ctx iris.Context) {
    // 获取表单数据
    var formData struct {
        Name    string `form:"name"`
        Email   string `form:"email"`
        Age     int    `form:"age"`
        Comment string `form:"comment"`
    }

    // 绑定表单数据
    if err := ctx.ReadForm(&formData); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "表单数据解析失败: " + err.Error(),
        })
        return
    }

    // 返回处理结果
    ctx.JSON(iris.Map{
        "code":    200,
        "message": "表单提交成功",
        "data":    formData,
    })
}

// UploadFile 文件上传接口
func UploadFile(ctx iris.Context) {
    // 获取上传的文件
    file, info, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "文件上传失败: " + err.Error(),
        })
        return
    }
    defer file.Close()

    // 验证文件类型
    if !isValidFileType(info.Filename) {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "不支持的文件类型",
        })
        return
    }

    // 验证文件大小（最大 10MB）
    if info.Size > 10*1024*1024 {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "文件大小不能超过 10MB",
        })
        return
    }

    // 保存文件
    filename := generateUniqueFilename(info.Filename)
    savePath := filepath.Join("static/uploads", filename)
    
    if err := saveUploadedFile(file, savePath); err != nil {
        ctx.JSON(iris.Map{
            "code":    500,
            "message": "文件保存失败: " + err.Error(),
        })
        return
    }

    // 返回上传结果
    ctx.JSON(iris.Map{
        "code":    200,
        "message": "文件上传成功",
        "data": iris.Map{
            "filename": filename,
            "original": info.Filename,
            "size":     info.Size,
            "url":      "/static/uploads/" + filename,
        },
    })
}

// GetProfile 获取用户资料（需要认证）
func GetProfile(ctx iris.Context) {
    // 从中间件获取用户信息
    userID := ctx.Values().Get("user_id")
    username := ctx.Values().Get("username")
    role := ctx.Values().Get("role")

    ctx.JSON(iris.Map{
        "code":    200,
        "message": "获取用户资料成功",
        "data": iris.Map{
            "user_id":  userID,
            "username": username,
            "role":     role,
        },
    })
}

// UpdateProfile 更新用户资料（需要认证）
func UpdateProfile(ctx iris.Context) {
    // 获取用户ID
    userID := ctx.Values().GetUintDefault("user_id", 0)
    if userID == 0 {
        ctx.JSON(iris.Map{
            "code":    401,
            "message": "无效的用户信息",
        })
        return
    }

    // 解析请求数据
    var updateData models.UpdateUserRequest
    if err := ctx.ReadJSON(&updateData); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "请求数据格式错误: " + err.Error(),
        })
        return
    }

    // 调用服务层更新用户
    user, err := services.UpdateUser(userID, &updateData)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    500,
            "message": "更新用户资料失败: " + err.Error(),
        })
        return
    }

    ctx.JSON(iris.Map{
        "code":    200,
        "message": "用户资料更新成功",
        "data":    user,
    })
}

// GetUsers 获取用户列表（需要管理员权限）
func GetUsers(ctx iris.Context) {
    // 获取分页参数
    page, _ := strconv.Atoi(ctx.URLParamDefault("page", "1"))
    pageSize, _ := strconv.Atoi(ctx.URLParamDefault("page_size", "10"))
    
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }

    // 调用服务层获取用户列表
    users, total, err := services.GetUsers(page, pageSize)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    500,
            "message": "获取用户列表失败: " + err.Error(),
        })
        return
    }

    // 返回分页响应
    ctx.JSON(models.NewPageResponse(200, "获取用户列表成功", users, page, pageSize, total))
}

// GetUser 获取单个用户信息
func GetUser(ctx iris.Context) {
    // 获取用户ID
    userID, err := ctx.Params().GetInt("id")
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "无效的用户ID",
        })
        return
    }

    // 调用服务层获取用户
    user, err := services.GetUserByID(userID)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    404,
            "message": "用户不存在",
        })
        return
    }

    ctx.JSON(models.NewResponse(200, "获取用户信息成功", user))
}

// UpdateUser 更新用户信息
func UpdateUser(ctx iris.Context) {
    // 获取用户ID
    userID, err := ctx.Params().GetInt("id")
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "无效的用户ID",
        })
        return
    }

    // 解析请求数据
    var updateData models.UpdateUserRequest
    if err := ctx.ReadJSON(&updateData); err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "请求数据格式错误: " + err.Error(),
        })
        return
    }

    // 调用服务层更新用户
    user, err := services.UpdateUser(userID, &updateData)
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    500,
            "message": "更新用户失败: " + err.Error(),
        })
        return
    }

    ctx.JSON(models.NewResponse(200, "用户更新成功", user))
}

// DeleteUser 删除用户
func DeleteUser(ctx iris.Context) {
    // 获取用户ID
    userID, err := ctx.Params().GetInt("id")
    if err != nil {
        ctx.JSON(iris.Map{
            "code":    400,
            "message": "无效的用户ID",
        })
        return
    }

    // 调用服务层删除用户
    if err := services.DeleteUser(userID); err != nil {
        ctx.JSON(iris.Map{
            "code":    500,
            "message": "删除用户失败: " + err.Error(),
        })
        return
    }

    ctx.JSON(models.NewResponse(200, "用户删除成功", nil))
}

// UsersPage 用户列表页面
func UsersPage(ctx iris.Context) {
    // 获取用户列表
    users, _, err := services.GetUsers(1, 50)
    if err != nil {
        ctx.ViewData("error", "获取用户列表失败")
        ctx.View("users.html")
        return
    }

    ctx.ViewData("title", "用户列表")
    ctx.ViewData("users", users)
    ctx.View("users.html")
}

// UserPage 用户详情页面
func UserPage(ctx iris.Context) {
    // 获取用户ID
    userID, err := ctx.Params().GetInt("id")
    if err != nil {
        ctx.ViewData("error", "无效的用户ID")
        ctx.View("user.html")
        return
    }

    // 获取用户信息
    user, err := services.GetUserByID(userID)
    if err != nil {
        ctx.ViewData("error", "用户不存在")
        ctx.View("user.html")
        return
    }

    ctx.ViewData("title", "用户详情")
    ctx.ViewData("user", user)
    ctx.View("user.html")
}

// NotFound 404 错误页面
func NotFound(ctx iris.Context) {
    ctx.ViewData("title", "页面未找到")
    ctx.ViewData("message", "抱歉，您访问的页面不存在")
    ctx.ViewData("code", "404")
    ctx.View("error.html")
}

// InternalServerError 500 错误页面
func InternalServerError(ctx iris.Context) {
    ctx.ViewData("title", "服务器错误")
    ctx.ViewData("message", "服务器内部错误，请稍后重试")
    ctx.ViewData("code", "500")
    ctx.View("error.html")
}

// APIDocs API 文档页面
func APIDocs(ctx iris.Context) {
    // 获取所有路由信息
    routes := ctx.App().GetRoutes()
    
    // 构建API文档
    apiDocs := make([]iris.Map, 0)
    for _, route := range routes {
        if strings.HasPrefix(route.Path, "/api") {
            apiDocs = append(apiDocs, iris.Map{
                "method": route.Method,
                "path":   route.Path,
                "name":   route.Name,
            })
        }
    }

    ctx.JSON(iris.Map{
        "code":    200,
        "message": "API 文档",
        "data": iris.Map{
            "title":       "Iris Go 框架 API 文档",
            "version":     "v1.0.0",
            "description": "Iris Go 框架学习项目的 API 接口文档",
            "routes":      apiDocs,
        },
    })
}

// 辅助函数

// isValidFileType 检查文件类型是否有效
func isValidFileType(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt"}
    
    for _, validExt := range validExts {
        if ext == validExt {
            return true
        }
    }
    return false
}

// generateUniqueFilename 生成唯一的文件名
func generateUniqueFilename(originalFilename string) string {
    ext := filepath.Ext(originalFilename)
    name := strings.TrimSuffix(originalFilename, ext)
    return fmt.Sprintf("%s_%d%s", name, 123456789, ext) // 简化的唯一文件名生成
}

// saveUploadedFile 保存上传的文件
func saveUploadedFile(file multipart.File, dst string) error {
    // 确保目录存在
    if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
        return err
    }

    // 创建目标文件
    f, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer f.Close()

    // 复制文件内容
    _, err = io.Copy(f, file)
    return err
}

// createDirIfNotExists 创建目录（如果不存在）
func createDirIfNotExists(dir string) error {
    return os.MkdirAll(dir, 0755)
}