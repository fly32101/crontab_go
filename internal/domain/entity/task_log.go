package entity

import (
	"time"
)

// TaskLog 任务执行日志
type TaskLog struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    int       `gorm:"not null"` // 关联的任务ID
	TaskName  string    `gorm:"not null"` // 任务名称（冗余存储，便于查询）
	StartTime time.Time `gorm:"not null"` // 任务开始执行时间
	EndTime   time.Time `gorm:"not null"` // 任务执行结束时间
	Success   bool      `gorm:"not null"` // 执行是否成功
	Output    string    `gorm:"type:text"` // 任务输出
	Error     string    `gorm:"type:text"` // 错误信息（如果有的话）
}

// TableName 设置表名
func (TaskLog) TableName() string {
	return "task_logs"
}