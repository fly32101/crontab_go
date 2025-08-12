package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"crontab_go/internal/domain/entity"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// SendNotification 发送通知
func (ns *NotificationService) SendNotification(config *entity.NotificationConfig, message *entity.NotificationMessage, notificationTypes []string) {
	log.Printf("SendNotification called for task %s with types: %v", message.TaskName, notificationTypes)
	
	for _, notificationType := range notificationTypes {
		log.Printf("Processing notification type: %s", notificationType)
		switch notificationType {
		case "email":
			if config.Email != nil {
				log.Printf("Sending email notification for task %s", message.TaskName)
				ns.sendEmailNotification(config.Email, message)
			} else {
				log.Printf("Email config is nil for task %s", message.TaskName)
			}
		case "dingtalk":
			if config.DingTalk != nil {
				log.Printf("Sending DingTalk notification for task %s", message.TaskName)
				ns.sendDingTalkNotification(config.DingTalk, message)
			} else {
				log.Printf("DingTalk config is nil for task %s", message.TaskName)
			}
		case "wechat":
			if config.WeChat != nil {
				log.Printf("Sending WeChat notification for task %s", message.TaskName)
				ns.sendWeChatNotification(config.WeChat, message)
			} else {
				log.Printf("WeChat config is nil for task %s", message.TaskName)
			}
		default:
			log.Printf("Unknown notification type: %s", notificationType)
		}
	}
}

// sendEmailNotification 发送邮件通知
func (ns *NotificationService) sendEmailNotification(config *entity.EmailConfig, message *entity.NotificationMessage) {
	subject := config.Subject
	if subject == "" {
		if message.Success {
			subject = fmt.Sprintf("任务执行成功 - %s", message.TaskName)
		} else {
			subject = fmt.Sprintf("任务执行失败 - %s", message.TaskName)
		}
	}

	body := ns.buildEmailBody(message)

	// 构建邮件内容
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		config.From,
		strings.Join(config.To, ","),
		subject,
		body)

	// 发送邮件
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)
	addr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	var err error
	if config.EnableTLS {
		err = ns.sendMailTLS(addr, auth, config.From, config.To, []byte(msg))
	} else {
		err = smtp.SendMail(addr, auth, config.From, config.To, []byte(msg))
	}

	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
	} else {
		log.Printf("Email notification sent successfully for task: %s", message.TaskName)
	}
}

// sendMailTLS 使用TLS发送邮件
func (ns *NotificationService) sendMailTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.StartTLS(&tls.Config{ServerName: strings.Split(addr, ":")[0]}); err != nil {
		return err
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

// buildEmailBody 构建邮件正文
func (ns *NotificationService) buildEmailBody(message *entity.NotificationMessage) string {
	status := "成功"
	statusColor := "#28a745"
	if !message.Success {
		status = "失败"
		statusColor = "#dc3545"
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>任务执行通知</title>
</head>
<body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        <div style="background-color: %s; color: white; padding: 20px; border-radius: 8px 8px 0 0;">
            <h2 style="margin: 0;">任务执行通知</h2>
        </div>
        <div style="padding: 20px;">
            <table style="width: 100%%; border-collapse: collapse;">
                <tr>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold; width: 120px;">任务名称:</td>
                    <td style="padding: 10px; border-bottom: 1px solid #eee;">%s</td>
                </tr>
                <tr>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold;">执行状态:</td>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; color: %s; font-weight: bold;">%s</td>
                </tr>
                <tr>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold;">开始时间:</td>
                    <td style="padding: 10px; border-bottom: 1px solid #eee;">%s</td>
                </tr>
                <tr>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold;">结束时间:</td>
                    <td style="padding: 10px; border-bottom: 1px solid #eee;">%s</td>
                </tr>
                <tr>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold;">执行时长:</td>
                    <td style="padding: 10px; border-bottom: 1px solid #eee;">%s</td>
                </tr>`,
		statusColor, message.TaskName, statusColor, status, message.StartTime, message.EndTime, message.Duration)

	if message.Output != "" {
		body += fmt.Sprintf(`
                <tr>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold; vertical-align: top;">执行输出:</td>
                    <td style="padding: 10px; border-bottom: 1px solid #eee;"><pre style="background-color: #f8f9fa; padding: 10px; border-radius: 4px; overflow-x: auto; white-space: pre-wrap;">%s</pre></td>
                </tr>`, message.Output)
	}

	if message.Error != "" {
		body += fmt.Sprintf(`
                <tr>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; font-weight: bold; vertical-align: top;">错误信息:</td>
                    <td style="padding: 10px; border-bottom: 1px solid #eee; color: #dc3545;"><pre style="background-color: #f8f9fa; padding: 10px; border-radius: 4px; overflow-x: auto; white-space: pre-wrap;">%s</pre></td>
                </tr>`, message.Error)
	}

	body += `
            </table>
        </div>
        <div style="padding: 20px; background-color: #f8f9fa; border-radius: 0 0 8px 8px; text-align: center; color: #6c757d; font-size: 12px;">
            此邮件由 Crontab 管理系统自动发送，请勿回复。
        </div>
    </div>
</body>
</html>`

	return body
}

// sendDingTalkNotification 发送钉钉通知
func (ns *NotificationService) sendDingTalkNotification(config *entity.DingTalkConfig, message *entity.NotificationMessage) {
	status := "✅ 成功"
	if !message.Success {
		status = "❌ 失败"
	}

	text := fmt.Sprintf("## 任务执行通知\n\n**任务名称:** %s\n\n**执行状态:** %s\n\n**开始时间:** %s\n\n**结束时间:** %s\n\n**执行时长:** %s",
		message.TaskName, status, message.StartTime, message.EndTime, message.Duration)

	if message.Output != "" {
		text += fmt.Sprintf("\n\n**执行输出:**\n```\n%s\n```", message.Output)
	}

	if message.Error != "" {
		text += fmt.Sprintf("\n\n**错误信息:**\n```\n%s\n```", message.Error)
	}

	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": fmt.Sprintf("任务执行通知 - %s", message.TaskName),
			"text":  text,
		},
	}

	// 添加@功能
	if len(config.AtMobiles) > 0 || config.AtAll {
		payload["at"] = map[string]interface{}{
			"atMobiles": config.AtMobiles,
			"isAtAll":   config.AtAll,
		}
	}

	// 如果有签名，添加时间戳和签名
	if config.Secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, config.Secret)
		h := hmac.New(sha256.New, []byte(config.Secret))
		h.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		webhookURL := fmt.Sprintf("%s&timestamp=%d&sign=%s", config.WebhookURL, timestamp, signature)
		ns.sendWebhookRequest(webhookURL, payload, "钉钉", message.TaskName)
	} else {
		ns.sendWebhookRequest(config.WebhookURL, payload, "钉钉", message.TaskName)
	}
}

// sendWeChatNotification 发送企业微信通知
func (ns *NotificationService) sendWeChatNotification(config *entity.WeChatConfig, message *entity.NotificationMessage) {
	status := "成功"
	statusColor := "info"
	if !message.Success {
		status = "失败"
		statusColor = "warning"
	}

	content := fmt.Sprintf("任务名称: %s\n执行状态: %s\n开始时间: %s\n结束时间: %s\n执行时长: %s",
		message.TaskName, status, message.StartTime, message.EndTime, message.Duration)

	if message.Output != "" {
		content += fmt.Sprintf("\n执行输出: %s", message.Output)
	}

	if message.Error != "" {
		content += fmt.Sprintf("\n错误信息: %s", message.Error)
	}

	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": fmt.Sprintf("## 任务执行通知\n\n<font color=\"%s\">%s</font>\n\n%s",
				statusColor, fmt.Sprintf("任务 %s 执行%s", message.TaskName, status), content),
		},
	}

	// 添加@功能
	if len(config.AtUserIds) > 0 || config.AtAll {
		mentioned := make([]string, 0)
		if config.AtAll {
			mentioned = append(mentioned, "@all")
		} else {
			for _, userId := range config.AtUserIds {
				mentioned = append(mentioned, fmt.Sprintf("<@%s>", userId))
			}
		}
		if len(mentioned) > 0 {
			payload["markdown"].(map[string]interface{})["content"] = fmt.Sprintf("%s\n\n%s",
				payload["markdown"].(map[string]interface{})["content"], strings.Join(mentioned, " "))
		}
	}

	ns.sendWebhookRequest(config.WebhookURL, payload, "企业微信", message.TaskName)
}

// sendWebhookRequest 发送webhook请求
func (ns *NotificationService) sendWebhookRequest(webhookURL string, payload map[string]interface{}, platform, taskName string) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal %s notification payload: %v", platform, err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send %s notification: %v", platform, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("%s notification sent successfully for task: %s", platform, taskName)
	} else {
		log.Printf("Failed to send %s notification, status code: %d", platform, resp.StatusCode)
	}
}