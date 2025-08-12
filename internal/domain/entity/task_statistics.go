package entity

import "time"

// TaskStatistics 任务统计信息
type TaskStatistics struct {
	TaskID              int     `json:"task_id"`
	TaskName            string  `json:"task_name"`
	TotalExecutions     int64   `json:"total_executions"`     // 总执行次数
	SuccessExecutions   int64   `json:"success_executions"`   // 成功执行次数
	FailureExecutions   int64   `json:"failure_executions"`   // 失败执行次数
	SuccessRate         float64 `json:"success_rate"`         // 成功率 (%)
	AverageExecutionTime float64 `json:"average_execution_time"` // 平均执行时间 (秒)
	LastExecutionTime   *time.Time `json:"last_execution_time"`   // 最后执行时间
	LastExecutionStatus bool    `json:"last_execution_status"` // 最后执行状态
}

// ExecutionTrend 执行趋势数据
type ExecutionTrend struct {
	Date            string `json:"date"`             // 日期 (YYYY-MM-DD)
	TotalExecutions int64  `json:"total_executions"` // 当日总执行次数
	SuccessCount    int64  `json:"success_count"`    // 当日成功次数
	FailureCount    int64  `json:"failure_count"`    // 当日失败次数
	SuccessRate     float64 `json:"success_rate"`    // 当日成功率
}

// TaskExecutionReport 任务执行报表
type TaskExecutionReport struct {
	ReportDate      string            `json:"report_date"`      // 报表日期
	TotalTasks      int64             `json:"total_tasks"`      // 总任务数
	ActiveTasks     int64             `json:"active_tasks"`     // 活跃任务数
	TotalExecutions int64             `json:"total_executions"` // 总执行次数
	SuccessRate     float64           `json:"success_rate"`     // 整体成功率
	TopTasks        []TaskStatistics  `json:"top_tasks"`        // 执行次数最多的任务
	RecentTrends    []ExecutionTrend  `json:"recent_trends"`    // 最近趋势
}

// TaskPerformanceMetrics 任务性能指标
type TaskPerformanceMetrics struct {
	TaskID              int     `json:"task_id"`
	TaskName            string  `json:"task_name"`
	MinExecutionTime    float64 `json:"min_execution_time"`    // 最短执行时间 (秒)
	MaxExecutionTime    float64 `json:"max_execution_time"`    // 最长执行时间 (秒)
	AverageExecutionTime float64 `json:"average_execution_time"` // 平均执行时间 (秒)
	MedianExecutionTime float64 `json:"median_execution_time"`  // 中位数执行时间 (秒)
	ExecutionTimeStdDev float64 `json:"execution_time_std_dev"` // 执行时间标准差
}

// HourlyExecutionStats 小时执行统计
type HourlyExecutionStats struct {
	Hour            int   `json:"hour"`             // 小时 (0-23)
	TotalExecutions int64 `json:"total_executions"` // 该小时总执行次数
	SuccessCount    int64 `json:"success_count"`    // 该小时成功次数
	FailureCount    int64 `json:"failure_count"`    // 该小时失败次数
}

// StatisticsRequest 统计请求参数
type StatisticsRequest struct {
	TaskID    *int       `json:"task_id,omitempty"`    // 可选，特定任务ID
	StartDate *time.Time `json:"start_date,omitempty"` // 开始日期
	EndDate   *time.Time `json:"end_date,omitempty"`   // 结束日期
	Days      int        `json:"days,omitempty"`       // 最近N天，默认30天
}

// NewStatisticsRequest 创建统计请求
func NewStatisticsRequest() *StatisticsRequest {
	return &StatisticsRequest{
		Days: 30, // 默认30天
	}
}