package service

import (
	"context"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

type roleService struct {
	repository repository.RoleRepository
}

func NewRoleService(repository repository.RoleRepository) RoleService {
	return &roleService{repository: repository}
}

func (s *roleService) ListRoles(ctx context.Context) ([]models.Role, error) {
	return s.repository.ListRoles(ctx)
}

var _ RoleService = (*roleService)(nil)
