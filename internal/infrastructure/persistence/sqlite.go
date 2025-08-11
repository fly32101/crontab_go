package persistence

import (
	"crontab_go/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	if err := db.AutoMigrate(&entity.Task{}, &entity.SystemStats{}, &entity.TaskLog{}); err != nil {
		return nil, err
	}

	return &SQLiteDB{Client: db}, nil
}