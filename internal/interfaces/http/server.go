package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
	handler *Handler
}

func NewServer(db *gorm.DB) *Server {
	engine := gin.Default()
	
	// 应用CORS中间件
	engine.Use(CORSMiddleware())
	
	handler := NewHandler(db)

	// 注册路由
	registerRoutes(engine, handler)

	return &Server{
		engine:  engine,
		handler: handler,
	}
}

func registerRoutes(engine *gin.Engine, handler *Handler) {
	api := engine.Group("/api/v1")
	
	// 认证相关路由（无需认证）
	auth := api.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
	}

	// 系统监控路由（无需认证，只读）
	system := api.Group("/system")
	{
		system.GET("/stats", handler.GetSystemStats)
	}

	// 需要认证的路由
	authMiddleware := AuthMiddleware(handler.authService)
	authenticated := api.Group("")
	authenticated.Use(authMiddleware)
	{
		// 用户信息
		authenticated.GET("/user", handler.GetCurrentUser)
		
		// 任务相关路由（需要认证）
		tasks := authenticated.Group("/tasks")
		{
			tasks.POST("", handler.CreateTask)           // 创建任务需要认证
			tasks.GET("", handler.ListTasks)             // 查看任务需要认证
			tasks.GET("/paginated", handler.ListTasksWithPagination)
			tasks.GET(":id", handler.GetTask)
			tasks.PUT(":id", handler.UpdateTask)         // 更新任务需要认证
			tasks.DELETE(":id", handler.DeleteTask)      // 删除任务需要认证
			tasks.GET(":id/logs", handler.GetTaskLogs)
			tasks.GET(":id/logs/paginated", handler.GetTaskLogsWithPagination)
			tasks.POST(":id/execute", handler.ExecuteTask) // 执行任务需要认证
		}

		// 日志相关路由（需要认证）
		logs := authenticated.Group("/logs")
		{
			logs.GET("", handler.GetAllLogs)             // 获取所有日志
			logs.GET("/paginated", handler.GetAllLogsWithPagination) // 分页获取所有日志
		}

		// 通知相关路由（需要认证）
		notifications := authenticated.Group("/notifications")
		{
			notifications.POST("/test", handler.TestNotification) // 测试通知
		}
	}
}

func (s *Server) Start() {
	log.Println("Starting server on :8080")
	
	// 提供前端构建后的静态文件
	s.engine.Static("/assets", "./web/dist/assets")
	s.engine.StaticFile("/favicon.ico", "./web/dist/favicon.ico")
	
	// 对于所有非API路由，返回 index.html (SPA路由)
	s.engine.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		// 否则返回前端应用的入口文件
		c.File("./web/dist/index.html")
	})

	if err := s.engine.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}