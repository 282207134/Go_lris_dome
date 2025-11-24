package main

import (
	"log"
	"os"

	"iris-cn-sample-project/config"
	"iris-cn-sample-project/controllers"
	"iris-cn-sample-project/database"
	"iris-cn-sample-project/middleware"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

// main 应用程序入口函数
func main() {
	// 创建 Iris 应用实例
	app := iris.New()

	// 配置应用
	configureApp(app)

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置路由
	setupRoutes(app)

	// 启动服务器
	port := config.GetConfig().Server.Port
	app.Listen(":" + port, iris.WithOptimizations)
}

// configureApp 配置应用程序
func configureApp(app *iris.Application) {
	// 设置应用配置
	app.Configure(iris.WithConfiguration(iris.Configuration{
		DisableInterruptHandler:          false,
		DisablePathCorrection:            false,
		EnablePathIntelligence:           true,
		EnablePathEscape:                 true,
		FireMethodNotAllowed:             true,
		DisableBodyConsumptionOnReturn:   false,
		DisableAutoFireStatusCode:        false,
		ResetOnFireErrorCode:             false,
		EnableOptimizations:              true,
		TimeFormat:                       "2006-01-02 15:04:05",
		Charset:                          "UTF-8",
		PostMaxMemory:                    32 << 20, // 32 MB
		TranslateFunctionContextKey:      "translate",
		TranslateLanguageContextKey:      "language",
		ViewLayoutContextKey:             "layout",
		ViewDataContextKey:               "data",
		RemoteAddrHeaders:                []string{"X-Forwarded-For"},
		RemoteAddrHeadersForce:           false,
		EnableOptimizationsUse:           true,
		EnableProtoJSON:                   true,
		DisableStartupLog:                false,
		DisableBanner:                    false,
		IgnoreServerErrors:               []string{},
		DisablePathCorrectionRedirection: false,
	}))

	// 添加全局中间件
	app.Use(recover.New())
	app.Use(logger.New())

	// 设置模板引擎
	setupTemplates(app)

	// 设置静态文件服务
	setupStaticFiles(app)
}

// setupRoutes 设置所有路由
func setupRoutes(app *iris.Application) {
	// 主页路由
	app.Handle("GET", "/", controllers.Index)
	app.Handle("GET", "/home", controllers.Home)

	// API 路由组
	api := app.Party("/api")
	{
		// 添加 CORS 中间件
		api.Use(middleware.CORS())

		// 基础示例接口
		api.Get("/hello", controllers.Hello)
		api.Get("/data/{id:int}", controllers.GetData)
		api.Post("/form", controllers.HandleForm)
		api.Post("/upload", controllers.UploadFile)

		// 认证相关接口
		auth := api.Party("/auth")
		{
			auth.Post("/login", controllers.Login)
			auth.Post("/register", controllers.Register)
			auth.Post("/refresh", controllers.RefreshToken)
		}

		// 需要认证的接口
		protected := api.Party("/protected")
		protected.Use(middleware.JWTAuthentication())
		{
			protected.Get("/profile", controllers.GetProfile)
			protected.Put("/profile", controllers.UpdateProfile)
		}

		// 用户管理接口
		users := api.Party("/users")
		users.Use(middleware.JWTAuthentication())
		{
			users.Get("/", controllers.GetUsers)
			users.Get("/{id:int}", controllers.GetUser)
			users.Put("/{id:int}", controllers.UpdateUser)
			users.Delete("/{id:int}", controllers.DeleteUser)
		}

		// API 文档
		api.Get("/docs", controllers.APIDocs)
	}

	// 用户页面路由
	pages := app.Party("/pages")
	{
		pages.Get("/users", controllers.UsersPage)
		pages.Get("/user/{id:int}", controllers.UserPage)
	}

	// 错误处理路由
	app.OnErrorCode(iris.StatusNotFound, controllers.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, controllers.InternalServerError)
}

// setupTemplates 设置模板引擎
func setupTemplates(app *iris.Application) {
	// 创建模板目录
	if err := os.MkdirAll("templates", 0755); err != nil {
		log.Printf("创建模板目录失败: %v", err)
	}

	// 注册 HTML 模板引擎
	tmpl := iris.HTML("./templates", ".html")
	tmpl.Layout("layouts/layout.html")
	tmpl.AddFunc("formatTime", func(t interface{}) string {
		return "2023-01-01" // 简化的时间格式化函数
	})
	app.RegisterView(tmpl)
}

// setupStaticFiles 设置静态文件服务
func setupStaticFiles(app *iris.Application) {
	// 创建静态文件目录
	dirs := []string{"static", "static/css", "static/js", "static/images"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("创建静态目录失败 %s: %v", dir, err)
		}
	}

	// 设置静态文件服务
	app.HandleDir("/static", iris.Dir("./static"))
	app.HandleDir("/assets", iris.Dir("./static"))
}