package repository

import "crontab_go/internal/domain/entity"

type TaskRepository interface {
	Create(task *entity.Task) error
	Update(task *entity.Task) error
	Delete(id int) error
	FindByID(id int) (*entity.Task, error)
	FindAll() ([]*entity.Task, error)
	FindEnabled() ([]*entity.Task, error)
}

type SystemRepository interface {
	SaveStats(stats *entity.SystemStats) error
	GetLatestStats() (*entity.SystemStats, error)
}