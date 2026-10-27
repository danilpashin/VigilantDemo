package repository

import (
	"context"

	"vigilant/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type AuditLogRepository interface {
	CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error
	ListAuditLogs(ctx context.Context) ([]models.AuditLog, error)
}

type RulesetRepository interface {
	CreateRuleset(ctx context.Context, ruleset *models.Ruleset) error
	ListRulesets(ctx context.Context) ([]models.Ruleset, error)
}

type RoleRepository interface {
	ListRoles(ctx context.Context) ([]models.Role, error)
}

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *models.Customer) error
	ListCustomers(ctx context.Context) ([]models.Customer, error)
	GetCustomer(ctx context.Context, id string) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, id string, customer *models.Customer) error
	DeleteCustomer(ctx context.Context, id string) error
}

type CardRepository interface {
	CreateCard(ctx context.Context, card *models.Card) error
	ListCards(ctx context.Context) ([]models.Card, error)
	GetCard(ctx context.Context, id string) (*models.Card, error)
	UpdateCard(ctx context.Context, id string, card *models.Card) error
	DeleteCard(ctx context.Context, id string) error
}
