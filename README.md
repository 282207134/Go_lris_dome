# Iris Go 框架学习示例项目

这是一个完整的 Iris Go Web 框架学习项目，包含了各种常用功能和最佳实践。

## 项目结构

```
iris-cn-sample-project/
├── main.go                 # 应用程序入口
├── go.mod                  # Go 模块文件
├── go.sum                  # Go 模块依赖版本锁定
├── README.md               # 项目说明文档
├── docs/                   # 详细文档目录
│   ├── 01-基础入门.md
│   ├── 02-路由系统.md
│   ├── 03-中间件使用.md
│   ├── 04-请求处理.md
│   ├── 05-响应处理.md
│   ├── 06-静态文件.md
│   ├── 07-模板引擎.md
│   ├── 08-数据库集成.md
│   ├── 09-身份验证.md
│   └── 10-错误处理.md
├── config/                 # 配置文件
│   └── config.go
├── controllers/            # 控制器
│   ├── user_controller.go
│   ├── auth_controller.go
│   └── api_controller.go
├── middleware/             # 中间件
│   ├── auth.go
│   ├── cors.go
│   ├── logger.go
│   └── recovery.go
├── models/                 # 数据模型
│   ├── user.go
│   └── response.go
├── services/               # 服务层
│   ├── user_service.go
│   └── auth_service.go
├── utils/                  # 工具函数
│   ├── jwt.go
│   ├── validator.go
│   └── response.go
├── static/                 # 静态文件
│   ├── css/
│   ├── js/
│   └── images/
├── templates/              # 模板文件
│   ├── index.html
│   ├── users.html
│   └── layouts/
└── database/               # 数据库相关
    └── database.go
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行项目

```bash
go run main.go
```

### 3. 访问应用

- 主页: http://localhost:8080
- API 文档: http://localhost:8080/api/docs
- 用户管理: http://localhost:8080/users

## 功能特性

### 核心功能
- ✅ RESTful API 设计
- ✅ JWT 身份验证
- ✅ 数据库集成 (GORM + SQLite)
- ✅ 请求验证
- ✅ 错误处理
- ✅ 日志记录
- ✅ CORS 跨域支持
- ✅ 静态文件服务
- ✅ HTML 模板渲染

### Iris 框架特性展示
- ✅ 路由系统（GET、POST、PUT、DELETE 等）
- ✅ 中间件机制
- ✅ 请求上下文处理
- ✅ 参数绑定和验证
- ✅ 响应格式化
- ✅ 文件上传下载
- ✅ WebSocket 支持
- ✅ 重定向和重写

## API 接口文档

### 认证相关
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/refresh` - 刷新令牌

### 用户管理
- `GET /api/users` - 获取用户列表
- `GET /api/users/:id` - 获取用户详情
- `PUT /api/users/:id` - 更新用户信息
- `DELETE /api/users/:id` - 删除用户

### 示例接口
- `GET /api/hello` - 简单问候接口
- `GET /api/data/:id` - 路径参数示例
- `POST /api/form` - 表单数据处理
- `POST /api/upload` - 文件上传

## 学习路径

建议按照以下顺序学习：

1. **基础入门** - 了解 Iris 基本概念和项目结构
2. **路由系统** - 学习路由配置和参数处理
3. **中间件使用** - 掌握中间件的编写和使用
4. **请求处理** - 学习各种请求类型的处理
5. **响应处理** - 掌握响应格式化和状态码
6. **静态文件** - 学习静态资源服务
7. **模板引擎** - 掌握 HTML 模板渲染
8. **数据库集成** - 学习数据库操作
9. **身份验证** - 掌握 JWT 认证机制
10. **错误处理** - 学习错误处理和日志记录

## 技术栈

- **Web 框架**: Iris v12
- **数据库**: SQLite (使用 GORM)
- **身份验证**: JWT
- **验证器**: go-playground/validator
- **Go 版本**: 1.21+

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个学习项目。

## 许可证

MIT License

## 相关链接

- [Iris 官方文档](https://iris-go.com/)
- [Iris GitHub 仓库](https://github.com/kataras/iris)
- [Iris 中文文档](https://github.com/kataras/iris/blob/main/README_ZH_HANS.md)