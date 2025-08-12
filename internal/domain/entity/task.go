package entity

type Task struct {
	ID                int    `json:"id" gorm:"primaryKey"`
	Name              string `json:"name" gorm:"not null"`
	Schedule          string `json:"schedule" gorm:"not null"`
	Command           string `json:"command" gorm:"not null"`
	Method            string `json:"method" gorm:"default:'GET'"` // HTTP请求方法
	Headers           string `json:"headers"`                    // HTTP请求头，JSON格式存储
	Enabled           bool   `json:"enabled" gorm:"default:true"`
	Description       string `json:"description"`
	NotifyOnSuccess   bool   `json:"notify_on_success" gorm:"default:false"`   // 成功时是否通知
	NotifyOnFailure   bool   `json:"notify_on_failure" gorm:"default:true"`    // 失败时是否通知
	NotificationTypes string `json:"notification_types"`                       // 通知类型，JSON格式存储 ["email", "dingtalk", "wechat"]
	NotificationConfig string `json:"notification_config"`                     // 通知配置，JSON格式存储
}

func (Task) TableName() string {
	return "tasks"
}