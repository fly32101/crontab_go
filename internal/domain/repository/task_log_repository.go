package repository

import "crontab_go/internal/domain/entity"

// TaskLogRepository 任务日志仓库接口
type TaskLogRepository interface {
	// Create 创建任务日志
	Create(log *entity.TaskLog) error
	
	// GetLogsByTaskID 根据任务ID获取任务日志
	GetLogsByTaskID(taskID int) ([]entity.TaskLog, error)
	
	// GetLogsByTaskIDWithPagination 根据任务ID分页获取任务日志
	GetLogsByTaskIDWithPagination(taskID int, req *entity.PaginationRequest) ([]entity.TaskLog, int64, error)

	// GetAllLogs 获取所有任务日志
	GetAllLogs() ([]entity.TaskLog, error)

	// GetAllLogsWithPagination 分页获取所有任务日志
	GetAllLogsWithPagination(page, pageSize int) ([]entity.TaskLog, int64, error)
}