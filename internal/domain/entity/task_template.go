package entity

import "time"

// TaskTemplate 任务模板
type TaskTemplate struct {
	ID                 int       `json:"id" gorm:"primaryKey"`
	Name               string    `json:"name" gorm:"not null"`                    // 模板名称
	Description        string    `json:"description"`                             // 模板描述
	Category           string    `json:"category" gorm:"default:'general'"`       // 模板分类
	Schedule           string    `json:"schedule" gorm:"not null"`                // Cron表达式
	Command            string    `json:"command" gorm:"not null"`                 // 命令或URL
	Method             string    `json:"method" gorm:"default:'GET'"`             // HTTP请求方法
	Headers            string    `json:"headers"`                                 // HTTP请求头，JSON格式存储
	NotifyOnSuccess    bool      `json:"notify_on_success" gorm:"default:false"` // 成功时是否通知
	NotifyOnFailure    bool      `json:"notify_on_failure" gorm:"default:true"`  // 失败时是否通知
	NotificationTypes  string    `json:"notification_types"`                     // 通知类型，JSON格式存储
	NotificationConfig string    `json:"notification_config"`                    // 通知配置，JSON格式存储
	Tags               string    `json:"tags"`                                    // 标签，JSON格式存储
	IsPublic           bool      `json:"is_public" gorm:"default:false"`         // 是否为公共模板
	CreatedBy          int       `json:"created_by"`                              // 创建者ID
	UsageCount         int       `json:"usage_count" gorm:"default:0"`           // 使用次数
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (TaskTemplate) TableName() string {
	return "task_templates"
}

// TaskTemplateCategory 任务模板分类
type TaskTemplateCategory struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description"`
	Icon        string `json:"icon"`        // 图标
	Color       string `json:"color"`       // 颜色
	SortOrder   int    `json:"sort_order"`  // 排序
}

func (TaskTemplateCategory) TableName() string {
	return "task_template_categories"
}

// CreateTaskFromTemplateRequest 从模板创建任务的请求
type CreateTaskFromTemplateRequest struct {
	TemplateID  int                    `json:"template_id" binding:"required"`
	TaskName    string                 `json:"task_name" binding:"required"`
	Overrides   map[string]interface{} `json:"overrides,omitempty"` // 覆盖的字段
	Enabled     bool                   `json:"enabled"`
}

// TaskTemplateSearchRequest 模板搜索请求
type TaskTemplateSearchRequest struct {
	Keyword    string `json:"keyword,omitempty"`
	Category   string `json:"category,omitempty"`
	IsPublic   *bool  `json:"is_public,omitempty"`
	CreatedBy  *int   `json:"created_by,omitempty"`
	Tags       string `json:"tags,omitempty"`
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
}

// TaskTemplateWithStats 带统计信息的任务模板
type TaskTemplateWithStats struct {
	TaskTemplate
	CategoryName string `json:"category_name"`
	CreatorName  string `json:"creator_name"`
	TagList      []string `json:"tag_list"`
}

// PopularTemplate 热门模板
type PopularTemplate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	UsageCount  int    `json:"usage_count"`
}

// TemplateStats 模板统计
type TemplateStats struct {
	TotalTemplates   int64 `json:"total_templates"`
	PublicTemplates  int64 `json:"public_templates"`
	PrivateTemplates int64 `json:"private_templates"`
	TotalUsage       int64 `json:"total_usage"`
	CategoriesCount  int64 `json:"categories_count"`
}