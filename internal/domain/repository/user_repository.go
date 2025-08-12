package repository

import "crontab_go/internal/domain/entity"

// UserRepository 用户仓库接口
type UserRepository interface {
	// Create 创建用户
	Create(user *entity.User) error
	
	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*entity.User, error)
	
	// FindByID 根据ID查找用户
	FindByID(id uint) (*entity.User, error)
	
	// Update 更新用户
	Update(user *entity.User) error
	
	// Delete 删除用户
	Delete(id uint) error
	
	// FindAll 获取所有用户
	FindAll() ([]*entity.User, error)
}