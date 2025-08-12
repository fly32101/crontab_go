package persistence

import (
	"crontab_go/internal/domain/entity"
	"crontab_go/internal/domain/repository"
	"gorm.io/gorm"
)

type SQLiteUserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &SQLiteUserRepository{DB: db}
}

func (r *SQLiteUserRepository) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}

func (r *SQLiteUserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteUserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteUserRepository) Update(user *entity.User) error {
	return r.DB.Save(user).Error
}

func (r *SQLiteUserRepository) Delete(id uint) error {
	return r.DB.Delete(&entity.User{}, id).Error
}

func (r *SQLiteUserRepository) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}