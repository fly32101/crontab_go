package persistence

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"gorm.io/gorm"
)

type SQLiteSystemRepository struct {
	DB *gorm.DB
}

func NewSystemRepository(db *gorm.DB) repository.SystemRepository {
	return &SQLiteSystemRepository{DB: db}
}

func (r *SQLiteSystemRepository) SaveStats(stats *entity.SystemStats) error {
	return r.DB.Create(stats).Error
}

func (r *SQLiteSystemRepository) GetLatestStats() (*entity.SystemStats, error) {
	var stats entity.SystemStats
	if err := r.DB.Order("timestamp DESC").First(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *SQLiteSystemRepository) GetStatsHistory(limit int, offset int) ([]entity.SystemStats, error) {
	var stats []entity.SystemStats
	if err := r.DB.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *SQLiteSystemRepository) GetStatsCount() (int64, error) {
	var count int64
	if err := r.DB.Model(&entity.SystemStats{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *SQLiteSystemRepository) KeepLatestStats(count int) error {
	// 获取总记录数
	var totalCount int64
	if err := r.DB.Model(&entity.SystemStats{}).Count(&totalCount).Error; err != nil {
		return err
	}
	
	// 如果记录数小于等于要保留的数量，不需要删除
	if totalCount <= int64(count) {
		return nil
	}
	
	// 删除多余的旧记录
	// 先获取要保留的最新记录的最小ID
	var minID uint
	if err := r.DB.Model(&entity.SystemStats{}).
		Select("id").
		Order("timestamp DESC").
		Limit(count).
		Offset(count - 1).
		Scan(&minID).Error; err != nil {
		return err
	}
	
	// 删除ID小于minID的记录
	return r.DB.Where("id < ?", minID).Delete(&entity.SystemStats{}).Error
}