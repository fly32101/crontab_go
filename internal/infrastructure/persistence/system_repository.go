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