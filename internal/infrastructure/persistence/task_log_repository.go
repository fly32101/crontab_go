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

// GetLogsByTaskIDWithPagination 根据任务ID分页获取任务日志
func (r *SQLiteTaskLogRepository) GetLogsByTaskIDWithPagination(taskID int, req *entity.PaginationRequest) ([]entity.TaskLog, int64, error) {
	var logs []entity.TaskLog
	var total int64
	
	// 获取总数
	if err := r.DB.Model(&entity.TaskLog{}).Where("task_id = ?", taskID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	if err := r.DB.Where("task_id = ?", taskID).
		Order("start_time DESC").
		Offset(req.GetOffset()).
		Limit(req.PageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	
	return logs, total, nil
}

// GetAllLogs 获取所有任务日志
func (r *SQLiteTaskLogRepository) GetAllLogs() ([]entity.TaskLog, error) {
	var logs []entity.TaskLog
	if err := r.DB.Order("start_time DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetAllLogsWithPagination 分页获取所有任务日志
func (r *SQLiteTaskLogRepository) GetAllLogsWithPagination(page, pageSize int) ([]entity.TaskLog, int64, error) {
	var logs []entity.TaskLog
	var total int64

	// 获取总数
	if err := r.DB.Model(&entity.TaskLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := r.DB.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}