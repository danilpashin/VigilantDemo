package service

import (
	"context"

	"vigilant/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type AuditLogService interface {
	CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error
	ListAuditLogs(ctx context.Context) ([]models.AuditLog, error)
}

type RulesetService interface {
	CreateRuleset(ctx context.Context, ruleset *models.Ruleset) error
	ListRulesets(ctx context.Context) ([]models.Ruleset, error)
}

type RoleService interface {
	ListRoles(ctx context.Context) ([]models.Role, error)
}
