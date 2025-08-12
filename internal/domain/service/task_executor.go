package service

import (
	"encoding/json"
	"fmt"
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
	taskRepo            repository.TaskRepository
	taskLogRepo         repository.TaskLogRepository
	cron                *cron.Cron
	runningTasks        map[int]cron.EntryID
	notificationService *NotificationService
}

func NewTaskExecutor(taskRepo repository.TaskRepository, taskLogRepo repository.TaskLogRepository) *TaskExecutor {
	return &TaskExecutor{
		taskRepo:            taskRepo,
		taskLogRepo:         taskLogRepo,
		cron:                cron.New(),
		runningTasks:        make(map[int]cron.EntryID),
		notificationService: NewNotificationService(),
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

	// 重新从数据库获取最新的任务配置（包含通知配置）
	latestTask, err := te.taskRepo.FindByID(task.ID)
	if err != nil {
		log.Printf("Failed to get latest task config for task %s: %v", task.Name, err)
		latestTask = task // 使用原任务配置作为备用
	}

	// 检查是否为HTTP请求
	if strings.HasPrefix(latestTask.Command, "http://") || strings.HasPrefix(latestTask.Command, "https://") {
		te.ExecuteHTTPRequest(latestTask)
	} else {
		// 执行系统命令
		te.ExecuteSystemCommand(latestTask)
	}
}

func (te *TaskExecutor) ExecuteSystemCommand(task *entity.Task) {
	// 记录开始时间
	startTime := time.Now()
	
	// 分割命令和参数
	parts := strings.Fields(task.Command)
	if len(parts) == 0 {
		log.Printf("Invalid command for task %s", task.Name)
		
		endTime := time.Now()
		// 记录日志到数据库
		taskLog := &entity.TaskLog{
			TaskID:    task.ID,
			TaskName:  task.Name,
			StartTime: startTime,
			EndTime:   endTime,
			Success:   false,
			Error:     "Invalid command",
		}
		if err := te.taskLogRepo.Create(taskLog); err != nil {
			log.Printf("Failed to save task log for task %s: %v", task.Name, err)
		}
		
		// 发送通知
		te.sendNotification(task, taskLog)
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	endTime := time.Now()
	
	var taskLog *entity.TaskLog
	
	if err != nil {
		log.Printf("Task %s failed: %v\nOutput: %s", task.Name, err, string(output))
		
		// 记录日志到数据库
		taskLog = &entity.TaskLog{
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
		taskLog = &entity.TaskLog{
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
	
	// 发送通知
	te.sendNotification(task, taskLog)
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
		
		endTime := time.Now()
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
		
		// 发送通知
		te.sendNotification(task, taskLog)
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
	
	var taskLog *entity.TaskLog
	
	if err != nil {
		log.Printf("HTTP request failed for task %s: %v", task.Name, err)
		
		// 记录日志到数据库
		taskLog = &entity.TaskLog{
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
	} else {
		defer resp.Body.Close()
		
		// 判断HTTP状态码是否表示成功
		success := resp.StatusCode >= 200 && resp.StatusCode < 300
		output := fmt.Sprintf("HTTP %s %s - Status: %s", method, task.Command, resp.Status)
		
		if success {
			log.Printf("HTTP request for task %s completed successfully with status %s", task.Name, resp.Status)
		} else {
			log.Printf("HTTP request for task %s failed with status %s", task.Name, resp.Status)
		}
		
		// 记录日志到数据库
		taskLog = &entity.TaskLog{
			TaskID:    task.ID,
			TaskName:  task.Name,
			StartTime: startTime,
			EndTime:   endTime,
			Success:   success,
			Output:    output,
		}
		
		if !success {
			taskLog.Error = fmt.Sprintf("HTTP request failed with status code: %d", resp.StatusCode)
		}
		
		if err := te.taskLogRepo.Create(taskLog); err != nil {
			log.Printf("Failed to save task log for task %s: %v", task.Name, err)
		}
	}
	
	// 发送通知
	te.sendNotification(task, taskLog)
}

// sendNotification 发送通知
func (te *TaskExecutor) sendNotification(task *entity.Task, taskLog *entity.TaskLog) {
	log.Printf("Checking notification for task %s: Success=%v, NotifyOnSuccess=%v, NotifyOnFailure=%v, NotificationTypes=%s", 
		task.Name, taskLog.Success, task.NotifyOnSuccess, task.NotifyOnFailure, task.NotificationTypes)
	
	// 检查是否需要发送通知
	shouldNotify := (taskLog.Success && task.NotifyOnSuccess) || (!taskLog.Success && task.NotifyOnFailure)
	if !shouldNotify {
		log.Printf("Notification not needed for task %s: shouldNotify=%v", task.Name, shouldNotify)
		return
	}
	
	if task.NotificationTypes == "" {
		log.Printf("No notification types configured for task %s", task.Name)
		return
	}

	// 解析通知类型
	var notificationTypes []string
	if err := json.Unmarshal([]byte(task.NotificationTypes), &notificationTypes); err != nil {
		log.Printf("Failed to parse notification types for task %s: %v", task.Name, err)
		return
	}

	if len(notificationTypes) == 0 {
		log.Printf("No notification types found for task %s", task.Name)
		return
	}
	
	log.Printf("Sending notification for task %s with types: %v", task.Name, notificationTypes)

	// 解析通知配置
	var notificationConfig entity.NotificationConfig
	if task.NotificationConfig != "" {
		if err := json.Unmarshal([]byte(task.NotificationConfig), &notificationConfig); err != nil {
			log.Printf("Failed to parse notification config for task %s: %v", task.Name, err)
			return
		}
	} else {
		log.Printf("No notification config found for task %s", task.Name)
		return
	}

	// 构建通知消息
	duration := taskLog.EndTime.Sub(taskLog.StartTime)
	message := &entity.NotificationMessage{
		TaskName:  task.Name,
		Success:   taskLog.Success,
		StartTime: taskLog.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:   taskLog.EndTime.Format("2006-01-02 15:04:05"),
		Duration:  duration.String(),
		Output:    taskLog.Output,
		Error:     taskLog.Error,
	}

	// 发送通知
	te.notificationService.SendNotification(&notificationConfig, message, notificationTypes)
}