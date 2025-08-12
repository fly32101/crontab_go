package entity

// NotificationConfig 通知配置
type NotificationConfig struct {
	Email    *EmailConfig    `json:"email,omitempty"`
	DingTalk *DingTalkConfig `json:"dingtalk,omitempty"`
	WeChat   *WeChatConfig   `json:"wechat,omitempty"`
}

// EmailConfig 邮件通知配置
type EmailConfig struct {
	SMTPHost     string   `json:"smtp_host"`
	SMTPPort     int      `json:"smtp_port"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	From         string   `json:"from"`
	To           []string `json:"to"`
	Subject      string   `json:"subject,omitempty"`
	EnableTLS    bool     `json:"enable_tls"`
}

// DingTalkConfig 钉钉通知配置
type DingTalkConfig struct {
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret,omitempty"`
	AtMobiles  []string `json:"at_mobiles,omitempty"`
	AtAll      bool   `json:"at_all,omitempty"`
}

// WeChatConfig 企业微信通知配置
type WeChatConfig struct {
	WebhookURL string   `json:"webhook_url"`
	AtUserIds  []string `json:"at_user_ids,omitempty"`
	AtAll      bool     `json:"at_all,omitempty"`
}

// NotificationMessage 通知消息
type NotificationMessage struct {
	TaskName    string `json:"task_name"`
	Success     bool   `json:"success"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Duration    string `json:"duration"`
	Output      string `json:"output,omitempty"`
	Error       string `json:"error,omitempty"`
}