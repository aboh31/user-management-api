package service

import (
	"user-management-api/internal/model"
	"user-management-api/internal/repository"
	"user-management-api/pkg/middleware"

	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(username, password, email string) (*model.UserResponse, error)
	Login(username, password string) (string, error)
	GetProfile(id uint) (*model.User, error)
	GetUsers() (*[]model.UserResponse, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) Register(username, password, email string) (*model.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Uuid:      uuid.New().String(),
		Username:  username,
		Password:  string(hashed),
		Email:     email,
		CreatedAt: time.Now(),
	}
	data, err := u.repo.Create(user)
	response := data.ConvertToResponse()

	return &response, err
}

func (u *userUsecase) Login(username, password string) (string, error) {
	user, err := u.repo.FindByUsername(username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", err
	}
	return middleware.GenerateJWT(user.ID)
}

func (u *userUsecase) GetProfile(id uint) (*model.User, error) {
	return u.repo.FindByID(id)
}

func (u *userUsecase) GetUsers() (*[]model.UserResponse, error) {
	data, err := u.repo.FindAll()
	response := []model.UserResponse{}
	for _, user := range *data {
		response = append(response, user.ConvertToResponse())
	}

	return &response, err
}
