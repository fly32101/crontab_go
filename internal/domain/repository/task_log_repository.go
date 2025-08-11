package repository

import "crontab_go/internal/domain/entity"

// TaskLogRepository 任务日志仓库接口
type TaskLogRepository interface {
	// Create 创建任务日志
	Create(log *entity.TaskLog) error
	
	// GetLogsByTaskID 根据任务ID获取任务日志
	GetLogsByTaskID(taskID int) ([]entity.TaskLog, error)
}