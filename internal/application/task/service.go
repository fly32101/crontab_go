package task

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"strings"

	"crontab_go/internal/domain/service"
)

type Service struct {
	taskRepo    repository.TaskRepository
	taskLogRepo repository.TaskLogRepository
}

func NewService(taskRepo repository.TaskRepository, taskLogRepo repository.TaskLogRepository) *Service {
	return &Service{taskRepo: taskRepo, taskLogRepo: taskLogRepo}
}

func (s *Service) CreateTask(task *entity.Task) error {
	return s.taskRepo.Create(task)
}

func (s *Service) UpdateTask(task *entity.Task) error {
	return s.taskRepo.Update(task)
}

func (s *Service) DeleteTask(id int) error {
	return s.taskRepo.Delete(id)
}

func (s *Service) GetTask(id int) (*entity.Task, error) {
	return s.taskRepo.FindByID(id)
}

func (s *Service) ListTasks() ([]*entity.Task, error) {
	return s.taskRepo.FindAll()
}

func (s *Service) ListEnabledTasks() ([]*entity.Task, error) {
	return s.taskRepo.FindEnabled()
}

// GetTaskLogs 获取任务执行日志
func (s *Service) GetTaskLogs(taskID int) ([]entity.TaskLog, error) {
	return s.taskLogRepo.GetLogsByTaskID(taskID)
}

// ExecuteTask 立即执行任务
func (s *Service) ExecuteTask(id int) error {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 创建TaskExecutor实例来执行任务
	taskExecutor := service.NewTaskExecutor(s.taskRepo, s.taskLogRepo)
	
	// 根据任务命令是否为URL来决定执行方式
	if strings.HasPrefix(task.Command, "http://") || strings.HasPrefix(task.Command, "https://") {
		taskExecutor.ExecuteHTTPRequest(task)
	} else {
		taskExecutor.ExecuteSystemCommand(task)
	}
	
	return nil
}