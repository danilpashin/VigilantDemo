package service

import (
	"context"

	"vigilant/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User, password string) error
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, user *models.User, password string) error
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

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *models.Customer) error
	ListCustomers(ctx context.Context) ([]models.Customer, error)
	GetCustomer(ctx context.Context, id string) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, id string, customer *models.Customer) error
	DeleteCustomer(ctx context.Context, id string) error
}

type CardService interface {
	CreateCard(ctx context.Context, card *models.Card) error
	ListCards(ctx context.Context) ([]models.Card, error)
	GetCard(ctx context.Context, id string) (*models.Card, error)
	UpdateCard(ctx context.Context, id string, card *models.Card) error
	DeleteCard(ctx context.Context, id string) error
}
