package http

import (
	"crontab_go/internal/application/auth"
	"crontab_go/internal/application/statistics"
	"crontab_go/internal/application/system"
	"crontab_go/internal/application/task"
	"crontab_go/internal/application/template"
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/service"
	"crontab_go/internal/infrastructure/persistence"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CORSMiddleware 处理跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type Handler struct {
	taskService       *task.Service
	systemService     *system.Service
	authService       *auth.Service
	statisticsService *statistics.Service
	templateService   *template.Service
}

func NewHandler(db *gorm.DB) *Handler {
	taskRepo := persistence.NewTaskRepository(db)
	taskLogRepo := persistence.NewTaskLogRepository(db)
	taskService := task.NewService(taskRepo, taskLogRepo)

	systemRepo := persistence.NewSystemRepository(db)
	systemService := system.NewService(systemRepo)

	userRepo := persistence.NewUserRepository(db)
	authService := auth.NewService(userRepo)

	statisticsService := statistics.NewService(taskRepo, taskLogRepo)

	templateRepo := persistence.NewTaskTemplateRepository(db)
	categoryRepo := persistence.NewTaskTemplateCategoryRepository(db)
	templateService := template.NewService(templateRepo, categoryRepo, taskRepo)

	return &Handler{
		taskService:       taskService,
		systemService:     systemService,
		authService:       authService,
		statisticsService: statisticsService,
		templateService:   templateService,
	}
}

func (h *Handler) CreateTask(c *gin.Context) {
	var task entity.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := h.taskService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) ListTasks(c *gin.Context) {
	tasks, err := h.taskService.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// ListTasksWithPagination 分页获取任务列表
func (h *Handler) ListTasksWithPagination(c *gin.Context) {
	var req entity.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	paginationReq := entity.NewPaginationRequest(req.Page, req.PageSize)

	response, err := h.taskService.ListTasksWithPagination(paginationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task entity.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = id
	if err := h.taskService.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.taskService.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// GetTaskLogs 获取任务执行日志
func (h *Handler) GetTaskLogs(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	logs, err := h.taskService.GetTaskLogs(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetTaskLogsWithPagination 分页获取任务执行日志
func (h *Handler) GetTaskLogsWithPagination(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req entity.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	paginationReq := entity.NewPaginationRequest(req.Page, req.PageSize)

	response, err := h.taskService.GetTaskLogsWithPagination(taskID, paginationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSystemStats(c *gin.Context) {
	// 返回实时系统统计数据
	stats, err := h.systemService.GetRealTimeStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 对百分比字段保留两位小数
	stats.CPUUsage = float64(int(stats.CPUUsage*100)) / 100
	stats.MemoryUsage = float64(int(stats.MemoryUsage*100)) / 100
	stats.DiskUsage = float64(int(stats.DiskUsage*100)) / 100

	c.JSON(http.StatusOK, stats)
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req entity.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": user})
}

// GetCurrentUser 获取当前用户信息
func (h *Handler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ExecuteTask 立即执行任务
func (h *Handler) ExecuteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.taskService.ExecuteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task executed successfully"})
}

// GetAllLogs 获取所有任务执行日志
func (h *Handler) GetAllLogs(c *gin.Context) {
	logs, err := h.taskService.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetAllLogsWithPagination 分页获取所有任务执行日志
func (h *Handler) GetAllLogsWithPagination(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 确保页码和每页大小有效
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	logs, total, err := h.taskService.GetAllLogsWithPagination(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// TestNotification 测试通知配置
func (h *Handler) TestNotification(c *gin.Context) {
	var req struct {
		NotificationTypes  []string                  `json:"notification_types"`
		NotificationConfig entity.NotificationConfig `json:"notification_config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建测试消息
	message := &entity.NotificationMessage{
		TaskName:  "测试任务",
		Success:   true,
		StartTime: time.Now().Add(-time.Minute).Format("2006-01-02 15:04:05"),
		EndTime:   time.Now().Format("2006-01-02 15:04:05"),
		Duration:  "1m0s",
		Output:    "这是一条测试通知消息",
	}

	// 发送测试通知
	notificationService := service.NewNotificationService()
	notificationService.SendNotification(&req.NotificationConfig, message, req.NotificationTypes)

	c.JSON(http.StatusOK, gin.H{"message": "测试通知已发送"})
}

// GetTaskStatistics 获取任务统计信息
func (h *Handler) GetTaskStatistics(c *gin.Context) {
	req := entity.NewStatisticsRequest()

	// 解析查询参数
	if days := c.Query("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d > 0 {
			req.Days = d
		}
	}

	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		if taskID, err := strconv.Atoi(taskIDStr); err == nil {
			req.TaskID = &taskID
		}
	}

	statistics, err := h.statisticsService.GetTaskStatistics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statistics)
}

// GetTaskStatisticsByID 获取特定任务的统计信息
func (h *Handler) GetTaskStatisticsByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	req := entity.NewStatisticsRequest()
	if days := c.Query("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d > 0 {
			req.Days = d
		}
	}

	statistics, err := h.statisticsService.GetTaskStatisticsByID(taskID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statistics)
}

// GetExecutionTrends 获取执行趋势数据
func (h *Handler) GetExecutionTrends(c *gin.Context) {
	req := entity.NewStatisticsRequest()

	// 解析查询参数
	if days := c.Query("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d > 0 {
			req.Days = d
		}
	}

	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		if taskID, err := strconv.Atoi(taskIDStr); err == nil {
			req.TaskID = &taskID
		}
	}

	trends, err := h.statisticsService.GetExecutionTrends(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetTaskExecutionReport 获取任务执行报表
func (h *Handler) GetTaskExecutionReport(c *gin.Context) {
	req := entity.NewStatisticsRequest()

	// 解析查询参数
	if days := c.Query("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d > 0 {
			req.Days = d
		}
	}

	report, err := h.statisticsService.GetTaskExecutionReport(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetTaskPerformanceMetrics 获取任务性能指标
func (h *Handler) GetTaskPerformanceMetrics(c *gin.Context) {
	req := entity.NewStatisticsRequest()

	// 解析查询参数
	if days := c.Query("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d > 0 {
			req.Days = d
		}
	}

	metrics, err := h.statisticsService.GetTaskPerformanceMetrics(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetHourlyExecutionStats 获取小时执行统计
func (h *Handler) GetHourlyExecutionStats(c *gin.Context) {
	req := entity.NewStatisticsRequest()

	// 解析查询参数
	if days := c.Query("days"); days != "" {
		if d, err := strconv.Atoi(days); err == nil && d > 0 {
			req.Days = d
		}
	}

	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		if taskID, err := strconv.Atoi(taskIDStr); err == nil {
			req.TaskID = &taskID
		}
	}

	stats, err := h.statisticsService.GetHourlyExecutionStats(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// 任务模板相关接口

// CreateTemplate 创建任务模板
func (h *Handler) CreateTemplate(c *gin.Context) {
	var template entity.TaskTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userEntity := user.(*entity.User)
	template.CreatedBy = int(userEntity.ID)

	if err := h.templateService.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetTemplate 获取任务模板
func (h *Handler) GetTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := h.templateService.GetTemplate(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// ListTemplates 获取模板列表
func (h *Handler) ListTemplates(c *gin.Context) {
	templates, err := h.templateService.ListTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// ListPublicTemplates 获取公共模板列表
func (h *Handler) ListPublicTemplates(c *gin.Context) {
	templates, err := h.templateService.ListPublicTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// ListMyTemplates 获取我的模板列表
func (h *Handler) ListMyTemplates(c *gin.Context) {
	// 获取当前用户ID
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userEntity := user.(*entity.User)

	templates, err := h.templateService.ListMyTemplates(int(userEntity.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// SearchTemplates 搜索模板
func (h *Handler) SearchTemplates(c *gin.Context) {
	var req entity.TaskTemplateSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	templates, total, err := h.templateService.SearchTemplates(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      templates,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// UpdateTemplate 更新任务模板
func (h *Handler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var template entity.TaskTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID = id
	if err := h.templateService.UpdateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate 删除任务模板
func (h *Handler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	if err := h.templateService.DeleteTemplate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// GetPopularTemplates 获取热门模板
func (h *Handler) GetPopularTemplates(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	templates, err := h.templateService.GetPopularTemplates(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// CreateTaskFromTemplate 从模板创建任务
func (h *Handler) CreateTaskFromTemplate(c *gin.Context) {
	var req entity.CreateTaskFromTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userEntity := user.(*entity.User)

	task, err := h.templateService.CreateTaskFromTemplate(&req, int(userEntity.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetTemplateStats 获取模板统计
func (h *Handler) GetTemplateStats(c *gin.Context) {
	stats, err := h.templateService.GetTemplateStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// 模板分类相关接口

// CreateCategory 创建模板分类
func (h *Handler) CreateCategory(c *gin.Context) {
	var category entity.TaskTemplateCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.templateService.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// ListCategories 获取分类列表
func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.templateService.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory 更新分类
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category entity.TaskTemplateCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = id
	if err := h.templateService.UpdateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory 删除分类
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.templateService.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
