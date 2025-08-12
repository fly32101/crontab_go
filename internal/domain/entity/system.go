package entity

import "time"

// SystemStats 系统统计信息
type SystemStats struct {
	ID              uint      `gorm:"primaryKey"`
	CPUUsage        float64   `gorm:"column:cpu_usage"`
	MemoryUsage     float64   `gorm:"column:memory_usage"`
	MemoryTotal     uint64    `gorm:"column:memory_total"`     // 总内存 (MB)
	MemoryUsed      uint64    `gorm:"column:memory_used"`      // 已用内存 (MB)
	MemoryFree      uint64    `gorm:"column:memory_free"`      // 空闲内存 (MB)
	DiskUsage       float64   `gorm:"column:disk_usage"`       // 磁盘使用率
	DiskTotal       uint64    `gorm:"column:disk_total"`       // 总磁盘空间 (GB)
	DiskUsed        uint64    `gorm:"column:disk_used"`        // 已用磁盘空间 (GB)
	DiskFree        uint64    `gorm:"column:disk_free"`        // 空闲磁盘空间 (GB)
	SystemLoad      float64   `gorm:"column:system_load"`
	NetworkRxBytes  uint64    `gorm:"column:network_rx_bytes"` // 网络接收字节数
	NetworkTxBytes  uint64    `gorm:"column:network_tx_bytes"` // 网络发送字节数
	ProcessCount    int       `gorm:"column:process_count"`    // 进程数量
	GoroutineCount  int       `gorm:"column:goroutine_count"`  // Goroutine数量
	Uptime          uint64    `gorm:"column:uptime"`           // 系统运行时间 (秒)
	Timestamp       time.Time `gorm:"column:timestamp"`
}

// TableName 设置表名
func (SystemStats) TableName() string {
	return "system_stats"
}