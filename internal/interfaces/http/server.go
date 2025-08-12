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
	}
}

func (s *Server) Start() {
	// 提供前端静态文件
	log.Println("Starting server on :8080")
	s.engine.Static("/static", "./web/static")
	s.engine.LoadHTMLFiles("./web/index.html")

	// 添加根路由以提供前端页面
	s.engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	
	// 登录页面路由
	s.engine.GET("/login.html", func(c *gin.Context) {
		c.File("./web/login.html")
	})

	if err := s.engine.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}