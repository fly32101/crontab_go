package template

import (
	"errors"
	"time"

	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
)

type Service struct {
	templateRepo repository.TaskTemplateRepository
	categoryRepo repository.TaskTemplateCategoryRepository
	taskRepo     repository.TaskRepository
}

func NewService(
	templateRepo repository.TaskTemplateRepository,
	categoryRepo repository.TaskTemplateCategoryRepository,
	taskRepo repository.TaskRepository,
) *Service {
	return &Service{
		templateRepo: templateRepo,
		categoryRepo: categoryRepo,
		taskRepo:     taskRepo,
	}
}

// CreateTemplate 创建任务模板
func (s *Service) CreateTemplate(template *entity.TaskTemplate) error {
	// 验证分类是否存在
	if template.Category != "" {
		if _, err := s.categoryRepo.FindByName(template.Category); err != nil {
			// 如果分类不存在，创建默认分类
			template.Category = "general"
		}
	}

	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	return s.templateRepo.Create(template)
}

// UpdateTemplate 更新任务模板
func (s *Service) UpdateTemplate(template *entity.TaskTemplate) error {
	// 检查模板是否存在
	existing, err := s.templateRepo.FindByID(template.ID)
	if err != nil {
		return err
	}

	// 保留创建时间和创建者
	template.CreatedAt = existing.CreatedAt
	template.CreatedBy = existing.CreatedBy
	template.UpdatedAt = time.Now()

	return s.templateRepo.Update(template)
}

// DeleteTemplate 删除任务模板
func (s *Service) DeleteTemplate(id int) error {
	return s.templateRepo.Delete(id)
}

// GetTemplate 获取任务模板
func (s *Service) GetTemplate(id int) (*entity.TaskTemplate, error) {
	return s.templateRepo.FindByID(id)
}

// ListTemplates 获取模板列表
func (s *Service) ListTemplates() ([]*entity.TaskTemplate, error) {
	return s.templateRepo.FindAll()
}

// ListTemplatesByCategory 按分类获取模板
func (s *Service) ListTemplatesByCategory(category string) ([]*entity.TaskTemplate, error) {
	return s.templateRepo.FindByCategory(category)
}

// ListMyTemplates 获取我的模板
func (s *Service) ListMyTemplates(userID int) ([]*entity.TaskTemplate, error) {
	return s.templateRepo.FindByCreator(userID)
}

// ListPublicTemplates 获取公共模板
func (s *Service) ListPublicTemplates() ([]*entity.TaskTemplate, error) {
	return s.templateRepo.FindPublicTemplates()
}

// SearchTemplates 搜索模板
func (s *Service) SearchTemplates(req *entity.TaskTemplateSearchRequest) ([]*entity.TaskTemplateWithStats, int64, error) {
	return s.templateRepo.Search(req)
}

// GetPopularTemplates 获取热门模板
func (s *Service) GetPopularTemplates(limit int) ([]*entity.PopularTemplate, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.templateRepo.GetPopularTemplates(limit)
}

// CreateTaskFromTemplate 从模板创建任务
func (s *Service) CreateTaskFromTemplate(req *entity.CreateTaskFromTemplateRequest, userID int) (*entity.Task, error) {
	// 获取模板
	template, err := s.templateRepo.FindByID(req.TemplateID)
	if err != nil {
		return nil, err
	}

	// 检查权限（如果是私有模板，只有创建者可以使用）
	if !template.IsPublic && template.CreatedBy != userID {
		return nil, errors.New("无权限使用此模板")
	}

	// 创建任务
	task := &entity.Task{
		Name:               req.TaskName,
		Schedule:           template.Schedule,
		Command:            template.Command,
		Method:             template.Method,
		Headers:            template.Headers,
		Enabled:            req.Enabled,
		Description:        template.Description,
		NotifyOnSuccess:    template.NotifyOnSuccess,
		NotifyOnFailure:    template.NotifyOnFailure,
		NotificationTypes:  template.NotificationTypes,
		NotificationConfig: template.NotificationConfig,
	}

	// 应用覆盖设置
	if req.Overrides != nil {
		s.applyOverrides(task, req.Overrides)
	}

	// 创建任务
	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	// 增加模板使用次数
	s.templateRepo.IncrementUsageCount(req.TemplateID)

	return task, nil
}

// applyOverrides 应用覆盖设置
func (s *Service) applyOverrides(task *entity.Task, overrides map[string]interface{}) {
	if schedule, ok := overrides["schedule"].(string); ok {
		task.Schedule = schedule
	}
	if command, ok := overrides["command"].(string); ok {
		task.Command = command
	}
	if method, ok := overrides["method"].(string); ok {
		task.Method = method
	}
	if headers, ok := overrides["headers"].(string); ok {
		task.Headers = headers
	}
	if description, ok := overrides["description"].(string); ok {
		task.Description = description
	}
	if notifyOnSuccess, ok := overrides["notify_on_success"].(bool); ok {
		task.NotifyOnSuccess = notifyOnSuccess
	}
	if notifyOnFailure, ok := overrides["notify_on_failure"].(bool); ok {
		task.NotifyOnFailure = notifyOnFailure
	}
}

// GetTemplateStats 获取模板统计
func (s *Service) GetTemplateStats() (*entity.TemplateStats, error) {
	return s.templateRepo.GetStats()
}

// 分类管理方法

// CreateCategory 创建分类
func (s *Service) CreateCategory(category *entity.TaskTemplateCategory) error {
	return s.categoryRepo.Create(category)
}

// UpdateCategory 更新分类
func (s *Service) UpdateCategory(category *entity.TaskTemplateCategory) error {
	return s.categoryRepo.Update(category)
}

// DeleteCategory 删除分类
func (s *Service) DeleteCategory(id int) error {
	return s.categoryRepo.Delete(id)
}

// GetCategory 获取分类
func (s *Service) GetCategory(id int) (*entity.TaskTemplateCategory, error) {
	return s.categoryRepo.FindByID(id)
}

// ListCategories 获取分类列表
func (s *Service) ListCategories() ([]*entity.TaskTemplateCategory, error) {
	return s.categoryRepo.FindAll()
}

// InitializeDefaultCategories 初始化默认分类
func (s *Service) InitializeDefaultCategories() error {
	defaultCategories := []*entity.TaskTemplateCategory{
		{
			Name:        "general",
			Description: "通用任务",
			Icon:        "AppstoreOutlined",
			Color:       "#1890ff",
			SortOrder:   1,
		},
		{
			Name:        "backup",
			Description: "备份任务",
			Icon:        "DatabaseOutlined",
			Color:       "#52c41a",
			SortOrder:   2,
		},
		{
			Name:        "monitoring",
			Description: "监控任务",
			Icon:        "MonitorOutlined",
			Color:       "#faad14",
			SortOrder:   3,
		},
		{
			Name:        "cleanup",
			Description: "清理任务",
			Icon:        "DeleteOutlined",
			Color:       "#f5222d",
			SortOrder:   4,
		},
		{
			Name:        "notification",
			Description: "通知任务",
			Icon:        "BellOutlined",
			Color:       "#722ed1",
			SortOrder:   5,
		},
		{
			Name:        "api",
			Description: "API调用",
			Icon:        "ApiOutlined",
			Color:       "#13c2c2",
			SortOrder:   6,
		},
	}

	for _, category := range defaultCategories {
		// 检查分类是否已存在
		if _, err := s.categoryRepo.FindByName(category.Name); err != nil {
			// 分类不存在，创建它
			if err := s.categoryRepo.Create(category); err != nil {
				return err
			}
		}
	}

	return nil
}

// CreateDefaultTemplates 创建默认模板
func (s *Service) CreateDefaultTemplates() error {
	defaultTemplates := []*entity.TaskTemplate{
		{
			Name:        "数据库备份",
			Description: "定期备份数据库",
			Category:    "backup",
			Schedule:    "0 0 2 * * *", // 每天凌晨2点
			Command:     "mysqldump -u root -p database_name > /backup/db_$(date +%Y%m%d).sql",
			Tags:        `["backup", "database", "mysql"]`,
			IsPublic:    true,
		},
		{
			Name:        "日志清理",
			Description: "清理7天前的日志文件",
			Category:    "cleanup",
			Schedule:    "0 0 3 * * *", // 每天凌晨3点
			Command:     "find /var/log -name '*.log' -mtime +7 -delete",
			Tags:        `["cleanup", "logs"]`,
			IsPublic:    true,
		},
		{
			Name:        "系统监控",
			Description: "检查系统资源使用情况",
			Category:    "monitoring",
			Schedule:    "0 */10 * * * *", // 每10分钟
			Command:     "df -h && free -m && uptime",
			Tags:        `["monitoring", "system"]`,
			IsPublic:    true,
		},
		{
			Name:        "健康检查",
			Description: "检查服务健康状态",
			Category:    "monitoring",
			Schedule:    "0 */5 * * * *", // 每5分钟
			Command:     "https://api.example.com/health",
			Method:      "GET",
			Tags:        `["monitoring", "health", "api"]`,
			IsPublic:    true,
		},
		{
			Name:            "每日报告",
			Description:     "发送每日系统报告",
			Category:        "notification",
			Schedule:        "0 0 9 * * *", // 每天上午9点
			Command:         "python /scripts/daily_report.py",
			Tags:            `["notification", "report"]`,
			IsPublic:        true,
			NotifyOnSuccess: true,
			NotifyOnFailure: true,
		},
	}

	for _, template := range defaultTemplates {
		template.CreatedAt = time.Now()
		template.UpdatedAt = time.Now()
		template.CreatedBy = 1 // 假设管理员用户ID为1

		// 检查模板是否已存在
		templates, err := s.templateRepo.FindAll()
		if err != nil {
			return err
		}

		exists := false
		for _, existing := range templates {
			if existing.Name == template.Name {
				exists = true
				break
			}
		}

		if !exists {
			if err := s.templateRepo.Create(template); err != nil {
				return err
			}
		}
	}

	return nil
}
