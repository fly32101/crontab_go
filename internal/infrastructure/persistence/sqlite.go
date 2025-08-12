package persistence

import (
	"crontab_go/internal/domain/entity"
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type SQLiteDB struct {
	Client *gorm.DB
}

func NewSQLiteDB(dsn string) (*SQLiteDB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(&entity.Task{}, &entity.SystemStats{}, &entity.TaskLog{}, &entity.User{}); err != nil {
		return nil, err
	}
	
	// 创建默认管理员用户
	if err := createDefaultAdmin(db); err != nil {
		return nil, err
	}

	return &SQLiteDB{Client: db}, nil
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin(db *gorm.DB) error {
	// 检查是否已存在管理员用户
	var count int64
	if err := db.Model(&entity.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}
	
	// 如果已存在管理员，则不创建
	if count > 0 {
		return nil
	}
	
	// 创建默认管理员用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	admin := &entity.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@example.com",
		Role:     "admin",
		IsActive: true,
	}
	
	if err := db.Create(admin).Error; err != nil {
		return err
	}
	
	log.Println("默认管理员用户已创建: username=admin, password=admin123")
	return nil
}