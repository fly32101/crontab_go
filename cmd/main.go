package main

import (
	"crontab_go/internal/application/system"
	"crontab_go/internal/application/template"
	"crontab_go/internal/domain/service"
	"crontab_go/internal/infrastructure/persistence"
	"crontab_go/internal/interfaces/http"
	"log"
	"time"
)

func main() {
	// 初始化数据库
	db, err := persistence.NewSQLiteDB("crontab.db")
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 初始化任务执行器
	taskRepo := persistence.NewTaskRepository(db.Client)
	taskLogRepo := persistence.NewTaskLogRepository(db.Client)
	executor := service.NewTaskExecutor(taskRepo, taskLogRepo)
	executor.Start()
	defer executor.Stop()

	// 初始化系统服务
	systemRepo := persistence.NewSystemRepository(db.Client)
	systemService := system.NewService(systemRepo)

	// 初始化模板服务并创建默认数据
	templateRepo := persistence.NewTaskTemplateRepository(db.Client)
	categoryRepo := persistence.NewTaskTemplateCategoryRepository(db.Client)
	templateService := template.NewService(templateRepo, categoryRepo, taskRepo)

	// 初始化默认分类和模板
	if err := templateService.InitializeDefaultCategories(); err != nil {
		log.Printf("Failed to initialize default categories: %v", err)
	}
	if err := templateService.CreateDefaultTemplates(); err != nil {
		log.Printf("Failed to create default templates: %v", err)
	}

	// 启动系统监控数据收集
	go func() {
		for {
			// 每10秒收集一次系统监控数据
			if err := systemService.CollectAndSaveStats(); err != nil {
				log.Printf("Failed to collect system stats: %v", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// 启动数据清理任务，每分钟清理一次旧数据
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			if err := systemService.CleanOldStats(); err != nil {
				log.Printf("Failed to clean old stats: %v", err)
			}
		}
	}()

	// 启动HTTP服务器
	server := http.NewServer(db.Client)
	server.Start()
}
