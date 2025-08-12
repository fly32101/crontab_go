package statistics

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"math"
	"sort"
	"time"
)

type Service struct {
	taskRepo    repository.TaskRepository
	taskLogRepo repository.TaskLogRepository
}

func NewService(taskRepo repository.TaskRepository, taskLogRepo repository.TaskLogRepository) *Service {
	return &Service{
		taskRepo:    taskRepo,
		taskLogRepo: taskLogRepo,
	}
}

// GetTaskStatistics 获取任务统计信息
func (s *Service) GetTaskStatistics(req *entity.StatisticsRequest) ([]entity.TaskStatistics, error) {
	// 获取所有任务
	tasks, err := s.taskRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var statistics []entity.TaskStatistics

	for _, task := range tasks {
		stat, err := s.calculateTaskStatistics(task, req)
		if err != nil {
			continue // 跳过错误的任务，继续处理其他任务
		}
		statistics = append(statistics, *stat)
	}

	return statistics, nil
}

// GetTaskStatisticsByID 获取特定任务的统计信息
func (s *Service) GetTaskStatisticsByID(taskID int, req *entity.StatisticsRequest) (*entity.TaskStatistics, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	return s.calculateTaskStatistics(task, req)
}

// calculateTaskStatistics 计算任务统计信息
func (s *Service) calculateTaskStatistics(task *entity.Task, req *entity.StatisticsRequest) (*entity.TaskStatistics, error) {
	// 获取任务日志
	logs, err := s.getTaskLogsInRange(task.ID, req)
	if err != nil {
		return nil, err
	}

	stat := &entity.TaskStatistics{
		TaskID:   task.ID,
		TaskName: task.Name,
	}

	if len(logs) == 0 {
		return stat, nil
	}

	// 计算基本统计
	stat.TotalExecutions = int64(len(logs))
	var totalDuration time.Duration
	var lastExecution *entity.TaskLog

	for _, log := range logs {
		if log.Success {
			stat.SuccessExecutions++
		} else {
			stat.FailureExecutions++
		}

		duration := log.EndTime.Sub(log.StartTime)
		totalDuration += duration

		// 找到最后执行的任务
		if lastExecution == nil || log.StartTime.After(lastExecution.StartTime) {
			lastExecution = &log
		}
	}

	// 计算成功率
	if stat.TotalExecutions > 0 {
		stat.SuccessRate = float64(stat.SuccessExecutions) / float64(stat.TotalExecutions) * 100
	}

	// 计算平均执行时间
	if stat.TotalExecutions > 0 {
		stat.AverageExecutionTime = totalDuration.Seconds() / float64(stat.TotalExecutions)
	}

	// 设置最后执行信息
	if lastExecution != nil {
		stat.LastExecutionTime = &lastExecution.StartTime
		stat.LastExecutionStatus = lastExecution.Success
	}

	return stat, nil
}

// GetExecutionTrends 获取执行趋势数据
func (s *Service) GetExecutionTrends(req *entity.StatisticsRequest) ([]entity.ExecutionTrend, error) {
	// 确定日期范围
	endDate := time.Now()
	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	startDate := endDate.AddDate(0, 0, -req.Days)
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	// 获取日期范围内的所有日志
	logs, err := s.getLogsInDateRange(startDate, endDate, req.TaskID)
	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	dailyStats := make(map[string]*entity.ExecutionTrend)

	// 初始化所有日期
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dailyStats[dateStr] = &entity.ExecutionTrend{
			Date: dateStr,
		}
	}

	// 统计每日数据
	for _, log := range logs {
		dateStr := log.StartTime.Format("2006-01-02")
		if trend, exists := dailyStats[dateStr]; exists {
			trend.TotalExecutions++
			if log.Success {
				trend.SuccessCount++
			} else {
				trend.FailureCount++
			}
		}
	}

	// 计算成功率并转换为切片
	var trends []entity.ExecutionTrend
	for _, trend := range dailyStats {
		if trend.TotalExecutions > 0 {
			trend.SuccessRate = float64(trend.SuccessCount) / float64(trend.TotalExecutions) * 100
		}
		trends = append(trends, *trend)
	}

	// 按日期排序
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Date < trends[j].Date
	})

	return trends, nil
}

// GetTaskExecutionReport 获取任务执行报表
func (s *Service) GetTaskExecutionReport(req *entity.StatisticsRequest) (*entity.TaskExecutionReport, error) {
	report := &entity.TaskExecutionReport{
		ReportDate: time.Now().Format("2006-01-02"),
	}

	// 获取所有任务
	tasks, err := s.taskRepo.FindAll()
	if err != nil {
		return nil, err
	}

	report.TotalTasks = int64(len(tasks))

	// 计算活跃任务数（最近30天有执行记录的任务）
	activeTasks := 0
	var allStats []entity.TaskStatistics
	var totalExecutions int64
	var totalSuccessExecutions int64

	for _, task := range tasks {
		stat, err := s.calculateTaskStatistics(task, req)
		if err != nil {
			continue
		}

		if stat.TotalExecutions > 0 {
			activeTasks++
			allStats = append(allStats, *stat)
			totalExecutions += stat.TotalExecutions
			totalSuccessExecutions += stat.SuccessExecutions
		}
	}

	report.ActiveTasks = int64(activeTasks)
	report.TotalExecutions = totalExecutions

	// 计算整体成功率
	if totalExecutions > 0 {
		report.SuccessRate = float64(totalSuccessExecutions) / float64(totalExecutions) * 100
	}

	// 获取执行次数最多的前10个任务
	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].TotalExecutions > allStats[j].TotalExecutions
	})

	topCount := 10
	if len(allStats) < topCount {
		topCount = len(allStats)
	}
	report.TopTasks = allStats[:topCount]

	// 获取最近趋势
	trends, err := s.GetExecutionTrends(req)
	if err != nil {
		return nil, err
	}
	report.RecentTrends = trends

	return report, nil
}

// GetTaskPerformanceMetrics 获取任务性能指标
func (s *Service) GetTaskPerformanceMetrics(req *entity.StatisticsRequest) ([]entity.TaskPerformanceMetrics, error) {
	tasks, err := s.taskRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var metrics []entity.TaskPerformanceMetrics

	for _, task := range tasks {
		metric, err := s.calculateTaskPerformanceMetrics(task, req)
		if err != nil {
			continue
		}
		if metric.AverageExecutionTime > 0 { // 只包含有执行记录的任务
			metrics = append(metrics, *metric)
		}
	}

	return metrics, nil
}

// calculateTaskPerformanceMetrics 计算任务性能指标
func (s *Service) calculateTaskPerformanceMetrics(task *entity.Task, req *entity.StatisticsRequest) (*entity.TaskPerformanceMetrics, error) {
	logs, err := s.getTaskLogsInRange(task.ID, req)
	if err != nil {
		return nil, err
	}

	metric := &entity.TaskPerformanceMetrics{
		TaskID:   task.ID,
		TaskName: task.Name,
	}

	if len(logs) == 0 {
		return metric, nil
	}

	// 计算执行时间
	var durations []float64
	var totalDuration float64

	for _, log := range logs {
		duration := log.EndTime.Sub(log.StartTime).Seconds()
		durations = append(durations, duration)
		totalDuration += duration
	}

	// 排序以计算中位数
	sort.Float64s(durations)

	// 基本统计
	metric.MinExecutionTime = durations[0]
	metric.MaxExecutionTime = durations[len(durations)-1]
	metric.AverageExecutionTime = totalDuration / float64(len(durations))

	// 中位数
	n := len(durations)
	if n%2 == 0 {
		metric.MedianExecutionTime = (durations[n/2-1] + durations[n/2]) / 2
	} else {
		metric.MedianExecutionTime = durations[n/2]
	}

	// 标准差
	var variance float64
	for _, duration := range durations {
		variance += math.Pow(duration-metric.AverageExecutionTime, 2)
	}
	variance /= float64(len(durations))
	metric.ExecutionTimeStdDev = math.Sqrt(variance)

	return metric, nil
}

// GetHourlyExecutionStats 获取小时执行统计
func (s *Service) GetHourlyExecutionStats(req *entity.StatisticsRequest) ([]entity.HourlyExecutionStats, error) {
	// 获取日志数据
	logs, err := s.getLogsInRange(req)
	if err != nil {
		return nil, err
	}

	// 初始化24小时统计
	hourlyStats := make([]entity.HourlyExecutionStats, 24)
	for i := 0; i < 24; i++ {
		hourlyStats[i] = entity.HourlyExecutionStats{Hour: i}
	}

	// 统计每小时数据
	for _, log := range logs {
		hour := log.StartTime.Hour()
		hourlyStats[hour].TotalExecutions++
		if log.Success {
			hourlyStats[hour].SuccessCount++
		} else {
			hourlyStats[hour].FailureCount++
		}
	}

	return hourlyStats, nil
}

// 辅助方法：获取指定任务在时间范围内的日志
func (s *Service) getTaskLogsInRange(taskID int, req *entity.StatisticsRequest) ([]entity.TaskLog, error) {
	// 这里需要扩展 TaskLogRepository 来支持日期范围查询
	// 暂时使用现有方法获取所有日志，然后过滤
	logs, err := s.taskLogRepo.GetLogsByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	return s.filterLogsByDateRange(logs, req), nil
}

// 辅助方法：获取日期范围内的日志
func (s *Service) getLogsInDateRange(startDate, endDate time.Time, taskID *int) ([]entity.TaskLog, error) {
	var logs []entity.TaskLog
	var err error

	if taskID != nil {
		logs, err = s.taskLogRepo.GetLogsByTaskID(*taskID)
	} else {
		logs, err = s.taskLogRepo.GetAllLogs()
	}

	if err != nil {
		return nil, err
	}

	// 过滤日期范围
	var filteredLogs []entity.TaskLog
	for _, log := range logs {
		if (log.StartTime.After(startDate) || log.StartTime.Equal(startDate)) &&
			(log.StartTime.Before(endDate) || log.StartTime.Equal(endDate)) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs, nil
}

// 辅助方法：获取范围内的所有日志
func (s *Service) getLogsInRange(req *entity.StatisticsRequest) ([]entity.TaskLog, error) {
	endDate := time.Now()
	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	startDate := endDate.AddDate(0, 0, -req.Days)
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	return s.getLogsInDateRange(startDate, endDate, req.TaskID)
}

// 辅助方法：按日期范围过滤日志
func (s *Service) filterLogsByDateRange(logs []entity.TaskLog, req *entity.StatisticsRequest) []entity.TaskLog {
	if req == nil {
		return logs
	}

	endDate := time.Now()
	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	startDate := endDate.AddDate(0, 0, -req.Days)
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	var filteredLogs []entity.TaskLog
	for _, log := range logs {
		if (log.StartTime.After(startDate) || log.StartTime.Equal(startDate)) &&
			(log.StartTime.Before(endDate) || log.StartTime.Equal(endDate)) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs
}