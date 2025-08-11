package entity

type Task struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Schedule    string `json:"schedule" gorm:"not null"`
	Command     string `json:"command" gorm:"not null"`
	Method      string `json:"method" gorm:"default:'GET'"` // HTTP请求方法
	Headers     string `json:"headers"`                    // HTTP请求头，JSON格式存储
	Enabled     bool   `json:"enabled" gorm:"default:true"`
	Description string `json:"description"`
}

func (Task) TableName() string {
	return "tasks"
}