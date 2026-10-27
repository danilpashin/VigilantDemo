package service

import (
	"context"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

type Service struct {
	users     repository.UserRepository
	auditLogs repository.AuditLogRepository
	rulesets  repository.RulesetRepository
	roles     repository.RoleRepository
}

func NewService(
	users repository.UserRepository,
	auditLogs repository.AuditLogRepository,
	rulesets repository.RulesetRepository,
	roles repository.RoleRepository,
) *Service {
	return &Service{users: users, auditLogs: auditLogs, rulesets: rulesets, roles: roles}
}

func (s *Service) CreateUser(ctx context.Context, user *models.User) error {
	return s.users.CreateUser(ctx, user)
}

func (s *Service) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.users.ListUsers(ctx)
}

func (s *Service) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.users.GetUser(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, id int, user *models.User) error {
	return s.users.UpdateUser(ctx, id, user)
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.users.DeleteUser(ctx, id)
}

func (s *Service) CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error {
	return s.auditLogs.CreateAuditLog(ctx, auditLog)
}

func (s *Service) ListAuditLogs(ctx context.Context) ([]models.AuditLog, error) {
	return s.auditLogs.ListAuditLogs(ctx)
}

func (s *Service) CreateRuleset(ctx context.Context, ruleset *models.Ruleset) error {
	return s.rulesets.CreateRuleset(ctx, ruleset)
}

func (s *Service) ListRulesets(ctx context.Context) ([]models.Ruleset, error) {
	return s.rulesets.ListRulesets(ctx)
}

func (s *Service) ListRoles(ctx context.Context) ([]models.Role, error) {
	return s.roles.ListRoles(ctx)
}

var _ UserService = (*Service)(nil)
var _ AuditLogService = (*Service)(nil)
var _ RulesetService = (*Service)(nil)
var _ RoleService = (*Service)(nil)
