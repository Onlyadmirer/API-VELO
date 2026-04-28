package service

import (
	"VELO-backend/internal/entity"
	"VELO-backend/internal/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user entity.RegisterUser) (*entity.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(user entity.RegisterUser) (*entity.User, error) {
	err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal hashing password")
	}

	newUser := entity.RegisterUser{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPass),
	}

	dataUser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return dataUser, nil

}
