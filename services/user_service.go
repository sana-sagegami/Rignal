package services

import (
	"auto-zen-backend/models"
	"auto-zen-backend/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Signup(username, password string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Signup(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	return s.repo.Create(user)
}
