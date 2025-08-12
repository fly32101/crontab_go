package system

import (
	"context"
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"log"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
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

// GetRealTimeStats 获取实时系统统计信息（不存储到数据库）
func (s *Service) GetRealTimeStats() (*entity.SystemStats, error) {
	ctx := context.Background()
	
	stats := &entity.SystemStats{
		CPUUsage:       s.getCPUUsage(ctx),
		MemoryUsage:    s.getMemoryUsage(ctx),
		MemoryTotal:    s.getMemoryTotal(ctx),
		MemoryUsed:     s.getMemoryUsed(ctx),
		MemoryFree:     s.getMemoryFree(ctx),
		DiskUsage:      s.getDiskUsage(ctx),
		DiskTotal:      s.getDiskTotal(ctx),
		DiskUsed:       s.getDiskUsed(ctx),
		DiskFree:       s.getDiskFree(ctx),
		SystemLoad:     s.getSystemLoad(ctx),
		NetworkRxBytes: s.getNetworkRxBytes(ctx),
		NetworkTxBytes: s.getNetworkTxBytes(ctx),
		ProcessCount:   s.getProcessCount(ctx),
		GoroutineCount: runtime.NumGoroutine(),
		Uptime:         s.getUptime(ctx),
		Timestamp:      time.Now(),
	}
	
	return stats, nil
}

// CleanOldStats 只保留最新的100条数据
func (s *Service) CleanOldStats() error {
	return s.SystemRepo.KeepLatestStats(100)
}

// CollectAndSaveStats 收集并保存系统统计信息
func (s *Service) CollectAndSaveStats() error {
	ctx := context.Background()
	
	stats := &entity.SystemStats{
		CPUUsage:       s.getCPUUsage(ctx),
		MemoryUsage:    s.getMemoryUsage(ctx),
		MemoryTotal:    s.getMemoryTotal(ctx),
		MemoryUsed:     s.getMemoryUsed(ctx),
		MemoryFree:     s.getMemoryFree(ctx),
		DiskUsage:      s.getDiskUsage(ctx),
		DiskTotal:      s.getDiskTotal(ctx),
		DiskUsed:       s.getDiskUsed(ctx),
		DiskFree:       s.getDiskFree(ctx),
		SystemLoad:     s.getSystemLoad(ctx),
		NetworkRxBytes: s.getNetworkRxBytes(ctx),
		NetworkTxBytes: s.getNetworkTxBytes(ctx),
		ProcessCount:   s.getProcessCount(ctx),
		GoroutineCount: runtime.NumGoroutine(),
		Uptime:         s.getUptime(ctx),
		Timestamp:      time.Now(),
	}
	
	return s.SaveStats(stats)
}

// 获取CPU使用率
func (s *Service) getCPUUsage(ctx context.Context) float64 {
	percentages, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		log.Printf("获取CPU使用率失败: %v", err)
		return 0.0
	}
	if len(percentages) > 0 {
		return percentages[0]
	}
	return 0.0
}

// 获取内存使用率
func (s *Service) getMemoryUsage(ctx context.Context) float64 {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		log.Printf("获取内存信息失败: %v", err)
		return 0.0
	}
	return vmStat.UsedPercent
}

// 获取总内存 (MB)
func (s *Service) getMemoryTotal(ctx context.Context) uint64 {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		log.Printf("获取内存信息失败: %v", err)
		return 0
	}
	return vmStat.Total / 1024 / 1024 // 转换为MB
}

// 获取已用内存 (MB)
func (s *Service) getMemoryUsed(ctx context.Context) uint64 {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		log.Printf("获取内存信息失败: %v", err)
		return 0
	}
	return vmStat.Used / 1024 / 1024 // 转换为MB
}

// 获取空闲内存 (MB)
func (s *Service) getMemoryFree(ctx context.Context) uint64 {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		log.Printf("获取内存信息失败: %v", err)
		return 0
	}
	return vmStat.Available / 1024 / 1024 // 转换为MB
}

// 获取磁盘使用率
func (s *Service) getDiskUsage(ctx context.Context) float64 {
	// 获取根目录的磁盘使用情况
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:\\"
	}
	
	diskStat, err := disk.UsageWithContext(ctx, path)
	if err != nil {
		log.Printf("获取磁盘信息失败: %v", err)
		return 0.0
	}
	return diskStat.UsedPercent
}

// 获取总磁盘空间 (GB)
func (s *Service) getDiskTotal(ctx context.Context) uint64 {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:\\"
	}
	
	diskStat, err := disk.UsageWithContext(ctx, path)
	if err != nil {
		log.Printf("获取磁盘信息失败: %v", err)
		return 0
	}
	return diskStat.Total / 1024 / 1024 / 1024 // 转换为GB
}

// 获取已用磁盘空间 (GB)
func (s *Service) getDiskUsed(ctx context.Context) uint64 {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:\\"
	}
	
	diskStat, err := disk.UsageWithContext(ctx, path)
	if err != nil {
		log.Printf("获取磁盘信息失败: %v", err)
		return 0
	}
	return diskStat.Used / 1024 / 1024 / 1024 // 转换为GB
}

// 获取空闲磁盘空间 (GB)
func (s *Service) getDiskFree(ctx context.Context) uint64 {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:\\"
	}
	
	diskStat, err := disk.UsageWithContext(ctx, path)
	if err != nil {
		log.Printf("获取磁盘信息失败: %v", err)
		return 0
	}
	return diskStat.Free / 1024 / 1024 / 1024 // 转换为GB
}

// 获取系统负载
func (s *Service) getSystemLoad(ctx context.Context) float64 {
	loadStat, err := load.AvgWithContext(ctx)
	if err != nil {
		log.Printf("获取系统负载失败: %v", err)
		// 在Windows上可能不支持，返回CPU使用率作为替代
		if runtime.GOOS == "windows" {
			return s.getCPUUsage(ctx) / 100 * 4 // 假设4核CPU
		}
		return 0.0
	}
	return loadStat.Load1 // 1分钟平均负载
}

// 获取网络接收字节数
func (s *Service) getNetworkRxBytes(ctx context.Context) uint64 {
	netStats, err := net.IOCountersWithContext(ctx, false)
	if err != nil {
		log.Printf("获取网络信息失败: %v", err)
		return 0
	}
	if len(netStats) > 0 {
		return netStats[0].BytesRecv
	}
	return 0
}

// 获取网络发送字节数
func (s *Service) getNetworkTxBytes(ctx context.Context) uint64 {
	netStats, err := net.IOCountersWithContext(ctx, false)
	if err != nil {
		log.Printf("获取网络信息失败: %v", err)
		return 0
	}
	if len(netStats) > 0 {
		return netStats[0].BytesSent
	}
	return 0
}

// 获取进程数量
func (s *Service) getProcessCount(ctx context.Context) int {
	pids, err := process.PidsWithContext(ctx)
	if err != nil {
		log.Printf("获取进程信息失败: %v", err)
		return 0
	}
	return len(pids)
}

// 获取系统运行时间
func (s *Service) getUptime(ctx context.Context) uint64 {
	hostStat, err := host.InfoWithContext(ctx)
	if err != nil {
		log.Printf("获取主机信息失败: %v", err)
		return 0
	}
	return hostStat.Uptime
}