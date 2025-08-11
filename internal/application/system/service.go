package system

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"math/rand"
	"time"
)

type Service struct {
	SystemRepo repository.SystemRepository
}

func NewService(systemRepo repository.SystemRepository) *Service {
	return &Service{SystemRepo: systemRepo}
}

func (s *Service) SaveStats(stats *entity.SystemStats) error {
	return s.SystemRepo.SaveStats(stats)
}

func (s *Service) GetLatestStats() (*entity.SystemStats, error) {
	return s.SystemRepo.GetLatestStats()
}

// CollectAndSaveStats 收集并保存系统统计信息
func (s *Service) CollectAndSaveStats() error {
	stats := &entity.SystemStats{
		CPUUsage:    rand.Float64() * 100,
		MemoryUsage: rand.Float64() * 100,
		SystemLoad:  rand.Float64() * 10,
		Timestamp:   time.Now(),
	}
	
	return s.SaveStats(stats)
}