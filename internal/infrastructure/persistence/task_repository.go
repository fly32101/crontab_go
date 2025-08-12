package persistence

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"gorm.io/gorm"
)

type SQLiteTaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &SQLiteTaskRepository{DB: db}
}

func (r *SQLiteTaskRepository) Create(task *entity.Task) error {
	return r.DB.Create(task).Error
}

func (r *SQLiteTaskRepository) Update(task *entity.Task) error {
	return r.DB.Save(task).Error
}

func (r *SQLiteTaskRepository) Delete(id int) error {
	return r.DB.Delete(&entity.Task{}, id).Error
}

func (r *SQLiteTaskRepository) FindByID(id int) (*entity.Task, error) {
	var task entity.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *SQLiteTaskRepository) FindAll() ([]*entity.Task, error) {
	var tasks []*entity.Task
	if err := r.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *SQLiteTaskRepository) FindEnabled() ([]*entity.Task, error) {
	var tasks []*entity.Task
	if err := r.DB.Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *SQLiteTaskRepository) FindWithPagination(req *entity.PaginationRequest) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64
	
	// 获取总数
	if err := r.DB.Model(&entity.Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	if err := r.DB.Offset(req.GetOffset()).Limit(req.PageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	
	return tasks, total, nil
}