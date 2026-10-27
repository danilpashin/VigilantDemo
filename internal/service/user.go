package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

const userPasswordHashCost = 12

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{repository: repository}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), userPasswordHashCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
	return s.repository.CreateUser(ctx, user)
}

func (s *userService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.repository.ListUsers(ctx)
}

func (s *userService) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.repository.GetUser(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id int, user *models.User, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), userPasswordHashCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
	return s.repository.UpdateUser(ctx, id, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.repository.DeleteUser(ctx, id)
}

var _ UserService = (*userService)(nil)
