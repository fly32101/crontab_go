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

// GetAllLogs 获取所有任务执行日志
func (s *Service) GetAllLogs() ([]entity.TaskLog, error) {
	return s.taskLogRepo.GetAllLogs()
}

// GetAllLogsWithPagination 分页获取所有任务执行日志
func (s *Service) GetAllLogsWithPagination(page, pageSize int) ([]entity.TaskLog, int64, error) {
	return s.taskLogRepo.GetAllLogsWithPagination(page, pageSize)
}

// ListTasksWithPagination 分页获取任务列表
func (s *Service) ListTasksWithPagination(req *entity.PaginationRequest) (*entity.PaginationResponse, error) {
	tasks, total, err := s.taskRepo.FindWithPagination(req)
	if err != nil {
		return nil, err
	}
	
	return entity.NewPaginationResponse(req.Page, req.PageSize, total, tasks), nil
}

// GetTaskLogsWithPagination 分页获取任务执行日志
func (s *Service) GetTaskLogsWithPagination(taskID int, req *entity.PaginationRequest) (*entity.PaginationResponse, error) {
	logs, total, err := s.taskLogRepo.GetLogsByTaskIDWithPagination(taskID, req)
	if err != nil {
		return nil, err
	}
	
	return entity.NewPaginationResponse(req.Page, req.PageSize, total, logs), nil
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