package main

import (
	"encoding/json"
	"fmt"
	"log"

	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/service"
	"crontab_go/internal/infrastructure/persistence"
)

func main() {
	// 初始化数据库
	db, err := persistence.NewSQLiteDB("crontab.db")
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 创建仓库
	taskRepo := persistence.NewTaskRepository(db.Client)
	taskLogRepo := persistence.NewTaskLogRepository(db.Client)

	// 获取所有任务
	tasks, err := taskRepo.FindAll()
	if err != nil {
		log.Fatal("Failed to get tasks:", err)
	}

	fmt.Printf("Found %d tasks:\n", len(tasks))
	for _, task := range tasks {
		fmt.Printf("\nTask ID: %d\n", task.ID)
		fmt.Printf("Name: %s\n", task.Name)
		fmt.Printf("NotifyOnSuccess: %v\n", task.NotifyOnSuccess)
		fmt.Printf("NotifyOnFailure: %v\n", task.NotifyOnFailure)
		fmt.Printf("NotificationTypes: %s\n", task.NotificationTypes)
		fmt.Printf("NotificationConfig: %s\n", task.NotificationConfig)

		// 解析通知类型
		if task.NotificationTypes != "" {
			var notificationTypes []string
			if err := json.Unmarshal([]byte(task.NotificationTypes), &notificationTypes); err != nil {
				fmt.Printf("Error parsing notification types: %v\n", err)
			} else {
				fmt.Printf("Parsed notification types: %v\n", notificationTypes)
			}
		}

		// 解析通知配置
		if task.NotificationConfig != "" {
			var notificationConfig entity.NotificationConfig
			if err := json.Unmarshal([]byte(task.NotificationConfig), &notificationConfig); err != nil {
				fmt.Printf("Error parsing notification config: %v\n", err)
			} else {
				fmt.Printf("Parsed notification config: %+v\n", notificationConfig)
			}
		}

		// 测试通知发送
		if task.NotifyOnSuccess && task.NotificationTypes != "" {
			fmt.Printf("Testing notification for task %s...\n", task.Name)
			
			// 创建测试任务日志
			taskLog := &entity.TaskLog{
				TaskID:   task.ID,
				TaskName: task.Name,
				Success:  true,
				Output:   "Test output",
			}

			// 创建任务执行器并发送通知
			executor := service.NewTaskExecutor(taskRepo, taskLogRepo)
			// 这里我们需要使用反射或者公开方法来测试通知
			fmt.Printf("Would send notification for task %s\n", task.Name)
		}
	}
}