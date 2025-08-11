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
	// 任务相关路由
	tasks := engine.Group("/api/v1/tasks")
	{
		tasks.POST("", handler.CreateTask)
		tasks.GET("", handler.ListTasks)
		tasks.GET(":id", handler.GetTask)
		tasks.PUT(":id", handler.UpdateTask)
		tasks.DELETE(":id", handler.DeleteTask)
		tasks.GET(":id/logs", handler.GetTaskLogs)
		tasks.POST(":id/execute", handler.ExecuteTask)
	}

	// 系统监控路由
	system := engine.Group("/api/v1/system")
	{
		system.GET("/stats", handler.GetSystemStats)
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

	if err := s.engine.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}