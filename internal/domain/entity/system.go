package entity

import "time"

// SystemStats 系统统计信息
type SystemStats struct {
	ID          uint      `gorm:"primaryKey"`
	CPUUsage    float64   `gorm:"column:cpu_usage"`
	MemoryUsage float64   `gorm:"column:memory_usage"`
	SystemLoad  float64   `gorm:"column:system_load"`
	Timestamp   time.Time `gorm:"column:timestamp"`
}

// TableName 设置表名
func (SystemStats) TableName() string {
	return "system_stats"
}