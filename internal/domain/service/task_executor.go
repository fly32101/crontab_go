package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
)

type TaskExecutor struct {
	taskRepo     repository.TaskRepository
	taskLogRepo  repository.TaskLogRepository
	cron         *cron.Cron
	runningTasks map[int]cron.EntryID
}

func NewTaskExecutor(taskRepo repository.TaskRepository, taskLogRepo repository.TaskLogRepository) *TaskExecutor {
	return &TaskExecutor{
		taskRepo:     taskRepo,
		taskLogRepo:  taskLogRepo,
		cron:         cron.New(),
		runningTasks: make(map[int]cron.EntryID),
	}
}

func (te *TaskExecutor) Start() {
	// 加载已启用的任务
	tasks, err := te.taskRepo.FindEnabled()
	if err != nil {
		log.Printf("Failed to load tasks: %v", err)
		return
	}

	for _, task := range tasks {
		te.scheduleTask(task)
	}

	te.cron.Start()
}

func (te *TaskExecutor) Stop() {
	te.cron.Stop()
}

func (te *TaskExecutor) scheduleTask(task *entity.Task) {
	entryID, err := te.cron.AddFunc(task.Schedule, func() {
		te.executeTask(task)
	})
	if err != nil {
		log.Printf("Failed to schedule task %s: %v", task.Name, err)
		return
	}

	te.runningTasks[task.ID] = entryID
	log.Printf("Scheduled task %s with schedule %s", task.Name, task.Schedule)
}

func (te *TaskExecutor) executeTask(task *entity.Task) {
	log.Printf("Executing task: %s", task.Name)

	// 检查是否为HTTP请求
	if strings.HasPrefix(task.Command, "http://") || strings.HasPrefix(task.Command, "https://") {
		te.ExecuteHTTPRequest(task)
	} else {
		// 执行系统命令
		te.ExecuteSystemCommand(task)
	}
}

func (te *TaskExecutor) ExecuteSystemCommand(task *entity.Task) {
	// 记录开始时间
	startTime := time.Now()
	
	// 分割命令和参数
	parts := strings.Fields(task.Command)
	if len(parts) == 0 {
		log.Printf("Invalid command for task %s", task.Name)
		
		// 记录日志到数据库
		taskLog := &entity.TaskLog{
			TaskID:    task.ID,
			TaskName:  task.Name,
			StartTime: startTime,
			EndTime:   time.Now(),
			Success:   false,
			Error:     "Invalid command",
		}
		if err := te.taskLogRepo.Create(taskLog); err != nil {
			log.Printf("Failed to save task log for task %s: %v", task.Name, err)
		}
		
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	endTime := time.Now()
	
	if err != nil {
		log.Printf("Task %s failed: %v\nOutput: %s", task.Name, err, string(output))
		
		// 记录日志到数据库
		taskLog := &entity.TaskLog{
			TaskID:    task.ID,
			TaskName:  task.Name,
			StartTime: startTime,
			EndTime:   endTime,
			Success:   false,
			Output:    string(output),
			Error:     err.Error(),
		}
		if err := te.taskLogRepo.Create(taskLog); err != nil {
			log.Printf("Failed to save task log for task %s: %v", task.Name, err)
		}
	} else {
		log.Printf("Task %s completed successfully\nOutput: %s", task.Name, string(output))
		
		// 记录日志到数据库
		taskLog := &entity.TaskLog{
			TaskID:    task.ID,
			TaskName:  task.Name,
			StartTime: startTime,
			EndTime:   endTime,
			Success:   true,
			Output:    string(output),
		}
		if err := te.taskLogRepo.Create(taskLog); err != nil {
			log.Printf("Failed to save task log for task %s: %v", task.Name, err)
		}
	}
}

func (te *TaskExecutor) ExecuteHTTPRequest(task *entity.Task) {
	// 记录开始时间
	startTime := time.Now()
	
	// 创建HTTP客户端
	client := &http.Client{Timeout: 30 * time.Second}

	// 确定请求方法，默认为GET
	method := task.Method
	if method == "" {
		method = "GET"
	}

	// 创建请求
	req, err := http.NewRequest(method, task.Command, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request for task %s: %v", task.Name, err)
		
		// 记录日志到数据库
		taskLog := &entity.TaskLog{
			TaskID:    task.ID,
			TaskName:  task.Name,
			StartTime: startTime,
			EndTime:   time.Now(),
			Success:   false,
			Error:     err.Error(),
		}
		if err := te.taskLogRepo.Create(taskLog); err != nil {
			log.Printf("Failed to save task log for task %s: %v", task.Name, err)
		}
		
		return
	}

	// 添加请求头
	if task.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(task.Headers), &headers); err != nil {
			log.Printf("Failed to parse headers for task %s: %v", task.Name, err)
		} else {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}
	}

	// 执行请求
	resp, err := client.Do(req)
	endTime := time.Now()
	
	if err != nil {
		log.Printf("HTTP request failed for task %s: %v", task.Name, err)
		
		// 记录日志到数据库
		taskLog := &entity.TaskLog{
			TaskID:    task.ID,
			TaskName:  task.Name,
			StartTime: startTime,
			EndTime:   endTime,
			Success:   false,
			Error:     err.Error(),
		}
		if err := te.taskLogRepo.Create(taskLog); err != nil {
			log.Printf("Failed to save task log for task %s: %v", task.Name, err)
		}
		
		return
	}
	defer resp.Body.Close()

	log.Printf("HTTP request for task %s completed with status %s", task.Name, resp.Status)
	
	// 记录日志到数据库
	taskLog := &entity.TaskLog{
		TaskID:    task.ID,
		TaskName:  task.Name,
		StartTime: startTime,
		EndTime:   endTime,
		Success:   true,
	}
	if err := te.taskLogRepo.Create(taskLog); err != nil {
		log.Printf("Failed to save task log for task %s: %v", task.Name, err)
	}
}