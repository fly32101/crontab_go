package repository

import "crontab_go/internal/domain/entity"

// SystemRepository 系统仓储接口
type SystemRepository interface {
	// SaveStats 保存系统统计信息
	SaveStats(stats *entity.SystemStats) error

	// GetLatestStats 获取最新的系统统计信息
	GetLatestStats() (*entity.SystemStats, error)

	// GetStatsHistory 获取系统统计历史数据
	GetStatsHistory(limit int, offset int) ([]entity.SystemStats, error)

	// GetStatsCount 获取系统统计数据的总数
	GetStatsCount() (int64, error)
}
