package repository

import (
	"user-management-api/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindAll() (*[]model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepo) FindAll() (*[]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return &users, err
}
