package repository

import "crontab_go/internal/domain/entity"

type TaskTemplateRepository interface {
	Create(template *entity.TaskTemplate) error
	Update(template *entity.TaskTemplate) error
	Delete(id int) error
	FindByID(id int) (*entity.TaskTemplate, error)
	FindAll() ([]*entity.TaskTemplate, error)
	FindByCategory(category string) ([]*entity.TaskTemplate, error)
	FindByCreator(creatorID int) ([]*entity.TaskTemplate, error)
	FindPublicTemplates() ([]*entity.TaskTemplate, error)
	Search(req *entity.TaskTemplateSearchRequest) ([]*entity.TaskTemplateWithStats, int64, error)
	GetPopularTemplates(limit int) ([]*entity.PopularTemplate, error)
	IncrementUsageCount(id int) error
	GetStats() (*entity.TemplateStats, error)
}

type TaskTemplateCategoryRepository interface {
	Create(category *entity.TaskTemplateCategory) error
	Update(category *entity.TaskTemplateCategory) error
	Delete(id int) error
	FindByID(id int) (*entity.TaskTemplateCategory, error)
	FindAll() ([]*entity.TaskTemplateCategory, error)
	FindByName(name string) (*entity.TaskTemplateCategory, error)
}