package main

import (
	"encoding/json"
	"log"
	"time"

	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/service"
)

func main() {
	// 创建通知服务
	notificationService := service.NewNotificationService()

	// 创建测试通知配置（钉钉）
	config := &entity.NotificationConfig{
		DingTalk: &entity.DingTalkConfig{
			WebhookURL: "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN_HERE",
			AtAll:      false,
		},
	}

	// 创建测试消息
	message := &entity.NotificationMessage{
		TaskName:  "测试任务",
		Success:   true,
		StartTime: time.Now().Add(-time.Minute).Format("2006-01-02 15:04:05"),
		EndTime:   time.Now().Format("2006-01-02 15:04:05"),
		Duration:  "1m0s",
		Output:    "测试输出内容",
	}

	// 发送通知
	log.Println("发送测试通知...")
	notificationService.SendNotification(config, message, []string{"dingtalk"})
	log.Println("测试通知发送完成")

	// 测试邮件通知（需要配置真实的SMTP信息）
	emailConfig := &entity.NotificationConfig{
		Email: &entity.EmailConfig{
			SMTPHost:  "smtp.gmail.com",
			SMTPPort:  587,
			Username:  "your-email@gmail.com",
			Password:  "your-password",
			From:      "your-email@gmail.com",
			To:        []string{"recipient@example.com"},
			EnableTLS: true,
		},
	}

	log.Println("发送邮件测试通知...")
	notificationService.SendNotification(emailConfig, message, []string{"email"})
	log.Println("邮件测试通知发送完成")
}