package persistence

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"gorm.io/gorm"
)

// SQLiteTaskLogRepository SQLite任务日志仓库实现
type SQLiteTaskLogRepository struct {
	DB *gorm.DB
}

// NewTaskLogRepository 创建任务日志仓库实例
func NewTaskLogRepository(db *gorm.DB) repository.TaskLogRepository {
	return &SQLiteTaskLogRepository{DB: db}
}

// Create 创建任务日志
func (r *SQLiteTaskLogRepository) Create(log *entity.TaskLog) error {
	return r.DB.Create(log).Error
}

// GetLogsByTaskID 根据任务ID获取任务日志
func (r *SQLiteTaskLogRepository) GetLogsByTaskID(taskID int) ([]entity.TaskLog, error) {
	var logs []entity.TaskLog
	if err := r.DB.Where("task_id = ?", taskID).Order("start_time DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}