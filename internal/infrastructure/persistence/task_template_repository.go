package persistence

import (
	"database/sql"
	"encoding/json"

	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"gorm.io/gorm"
)

type taskTemplateRepository struct {
	db *gorm.DB
}

func NewTaskTemplateRepository(db *gorm.DB) repository.TaskTemplateRepository {
	return &taskTemplateRepository{db: db}
}

func (r *taskTemplateRepository) Create(template *entity.TaskTemplate) error {
	return r.db.Create(template).Error
}

func (r *taskTemplateRepository) Update(template *entity.TaskTemplate) error {
	return r.db.Save(template).Error
}

func (r *taskTemplateRepository) Delete(id int) error {
	return r.db.Delete(&entity.TaskTemplate{}, id).Error
}

func (r *taskTemplateRepository) FindByID(id int) (*entity.TaskTemplate, error) {
	var template entity.TaskTemplate
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *taskTemplateRepository) FindAll() ([]*entity.TaskTemplate, error) {
	var templates []*entity.TaskTemplate
	err := r.db.Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *taskTemplateRepository) FindByCategory(category string) ([]*entity.TaskTemplate, error) {
	var templates []*entity.TaskTemplate
	err := r.db.Where("category = ?", category).Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *taskTemplateRepository) FindByCreator(creatorID int) ([]*entity.TaskTemplate, error) {
	var templates []*entity.TaskTemplate
	err := r.db.Where("created_by = ?", creatorID).Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *taskTemplateRepository) FindPublicTemplates() ([]*entity.TaskTemplate, error) {
	var templates []*entity.TaskTemplate
	err := r.db.Where("is_public = ?", true).Order("usage_count DESC, created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *taskTemplateRepository) Search(req *entity.TaskTemplateSearchRequest) ([]*entity.TaskTemplateWithStats, int64, error) {
	query := r.db.Model(&entity.TaskTemplate{}).
		Select(`task_templates.*, 
				task_template_categories.name as category_name,
				users.username as creator_name`).
		Joins("LEFT JOIN task_template_categories ON task_templates.category = task_template_categories.name").
		Joins("LEFT JOIN users ON task_templates.created_by = users.id")

	// 应用搜索条件
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("task_templates.name LIKE ? OR task_templates.description LIKE ?", keyword, keyword)
	}

	if req.Category != "" {
		query = query.Where("task_templates.category = ?", req.Category)
	}

	if req.IsPublic != nil {
		query = query.Where("task_templates.is_public = ?", *req.IsPublic)
	}

	if req.CreatedBy != nil {
		query = query.Where("task_templates.created_by = ?", *req.CreatedBy)
	}

	if req.Tags != "" {
		query = query.Where("task_templates.tags LIKE ?", "%"+req.Tags+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Order("task_templates.usage_count DESC, task_templates.created_at DESC")

	var results []struct {
		entity.TaskTemplate
		CategoryName string `json:"category_name"`
		CreatorName  string `json:"creator_name"`
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, 0, err
	}

	// 转换结果
	var templates []*entity.TaskTemplateWithStats
	for _, result := range results {
		template := &entity.TaskTemplateWithStats{
			TaskTemplate: result.TaskTemplate,
			CategoryName: result.CategoryName,
			CreatorName:  result.CreatorName,
		}

		// 解析标签
		if result.Tags != "" {
			var tags []string
			if err := json.Unmarshal([]byte(result.Tags), &tags); err == nil {
				template.TagList = tags
			}
		}

		templates = append(templates, template)
	}

	return templates, total, nil
}

func (r *taskTemplateRepository) GetPopularTemplates(limit int) ([]*entity.PopularTemplate, error) {
	var templates []*entity.PopularTemplate
	err := r.db.Model(&entity.TaskTemplate{}).
		Select("id, name, description, category, usage_count").
		Where("is_public = ?", true).
		Order("usage_count DESC").
		Limit(limit).
		Find(&templates).Error
	return templates, err
}

func (r *taskTemplateRepository) IncrementUsageCount(id int) error {
	return r.db.Model(&entity.TaskTemplate{}).Where("id = ?", id).UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error
}

func (r *taskTemplateRepository) GetStats() (*entity.TemplateStats, error) {
	stats := &entity.TemplateStats{}

	// 总模板数
	if err := r.db.Model(&entity.TaskTemplate{}).Count(&stats.TotalTemplates).Error; err != nil {
		return nil, err
	}

	// 公共模板数
	if err := r.db.Model(&entity.TaskTemplate{}).Where("is_public = ?", true).Count(&stats.PublicTemplates).Error; err != nil {
		return nil, err
	}

	// 私有模板数
	stats.PrivateTemplates = stats.TotalTemplates - stats.PublicTemplates

	// 总使用次数
	var totalUsage sql.NullInt64
	if err := r.db.Model(&entity.TaskTemplate{}).Select("SUM(usage_count)").Scan(&totalUsage).Error; err != nil {
		return nil, err
	}
	if totalUsage.Valid {
		stats.TotalUsage = totalUsage.Int64
	}

	// 分类数
	if err := r.db.Model(&entity.TaskTemplateCategory{}).Count(&stats.CategoriesCount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// TaskTemplateCategoryRepository 实现
type taskTemplateCategoryRepository struct {
	db *gorm.DB
}

func NewTaskTemplateCategoryRepository(db *gorm.DB) repository.TaskTemplateCategoryRepository {
	return &taskTemplateCategoryRepository{db: db}
}

func (r *taskTemplateCategoryRepository) Create(category *entity.TaskTemplateCategory) error {
	return r.db.Create(category).Error
}

func (r *taskTemplateCategoryRepository) Update(category *entity.TaskTemplateCategory) error {
	return r.db.Save(category).Error
}

func (r *taskTemplateCategoryRepository) Delete(id int) error {
	return r.db.Delete(&entity.TaskTemplateCategory{}, id).Error
}

func (r *taskTemplateCategoryRepository) FindByID(id int) (*entity.TaskTemplateCategory, error) {
	var category entity.TaskTemplateCategory
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *taskTemplateCategoryRepository) FindAll() ([]*entity.TaskTemplateCategory, error) {
	var categories []*entity.TaskTemplateCategory
	err := r.db.Order("sort_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

func (r *taskTemplateCategoryRepository) FindByName(name string) (*entity.TaskTemplateCategory, error) {
	var category entity.TaskTemplateCategory
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}