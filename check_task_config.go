package main

import (
	"encoding/json"
	"fmt"
	"log"

	"crontab_go/internal/domain/entity"
	"crontab_go/internal/infrastructure/persistence"
)

func main() {
	// 初始化数据库
	db, err := persistence.NewSQLiteDB("crontab.db")
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 创建任务仓库
	taskRepo := persistence.NewTaskRepository(db.Client)

	// 获取所有任务
	tasks, err := taskRepo.FindAll()
	if err != nil {
		log.Fatal("Failed to get tasks:", err)
	}

	fmt.Printf("=== 数据库中的任务配置 ===\n\n")
	
	for _, task := range tasks {
		fmt.Printf("任务 ID: %d\n", task.ID)
		fmt.Printf("任务名称: %s\n", task.Name)
		fmt.Printf("命令: %s\n", task.Command)
		fmt.Printf("启用状态: %v\n", task.Enabled)
		fmt.Printf("成功时通知: %v\n", task.NotifyOnSuccess)
		fmt.Printf("失败时通知: %v\n", task.NotifyOnFailure)
		fmt.Printf("通知类型 (原始): %s\n", task.NotificationTypes)
		fmt.Printf("通知配置 (原始): %s\n", task.NotificationConfig)

		// 解析通知类型
		if task.NotificationTypes != "" {
			var notificationTypes []string
			if err := json.Unmarshal([]byte(task.NotificationTypes), &notificationTypes); err != nil {
				fmt.Printf("❌ 解析通知类型失败: %v\n", err)
			} else {
				fmt.Printf("✅ 解析的通知类型: %v\n", notificationTypes)
			}
		} else {
			fmt.Printf("⚠️  通知类型为空\n")
		}

		// 解析通知配置
		if task.NotificationConfig != "" {
			var notificationConfig entity.NotificationConfig
			if err := json.Unmarshal([]byte(task.NotificationConfig), &notificationConfig); err != nil {
				fmt.Printf("❌ 解析通知配置失败: %v\n", err)
			} else {
				fmt.Printf("✅ 解析的通知配置:\n")
				if notificationConfig.Email != nil {
					fmt.Printf("  - 邮件配置: SMTP=%s:%d, From=%s, To=%v\n", 
						notificationConfig.Email.SMTPHost, 
						notificationConfig.Email.SMTPPort,
						notificationConfig.Email.From,
						notificationConfig.Email.To)
				}
				if notificationConfig.DingTalk != nil {
					fmt.Printf("  - 钉钉配置: WebhookURL=%s\n", notificationConfig.DingTalk.WebhookURL)
				}
				if notificationConfig.WeChat != nil {
					fmt.Printf("  - 企业微信配置: WebhookURL=%s\n", notificationConfig.WeChat.WebhookURL)
				}
			}
		} else {
			fmt.Printf("⚠️  通知配置为空\n")
		}

		// 判断是否应该发送通知
		fmt.Printf("\n通知判断逻辑:\n")
		fmt.Printf("- 如果任务成功且 NotifyOnSuccess=true: %v\n", task.NotifyOnSuccess)
		fmt.Printf("- 如果任务失败且 NotifyOnFailure=true: %v\n", task.NotifyOnFailure)
		fmt.Printf("- 通知类型不为空: %v\n", task.NotificationTypes != "")
		fmt.Printf("- 通知配置不为空: %v\n", task.NotificationConfig != "")

		fmt.Printf("\n" + "="*50 + "\n\n")
	}

	if len(tasks) == 0 {
		fmt.Println("数据库中没有任务")
	}
}