package service

import (
	"context"

	"uuid"

	"golang.org/x/crypto/bcrypt"

	"vigilant/internal/models"
	"vigilant/internal/repository"
)

const userPasswordHashCost = 12

type Service struct {
	users     repository.UserRepository
	auditLogs repository.AuditLogRepository
	rulesets  repository.RulesetRepository
	roles     repository.RoleRepository
	customers repository.CustomerRepository
	cards     repository.CardRepository
}

func NewService(
	users repository.UserRepository,
	auditLogs repository.AuditLogRepository,
	rulesets repository.RulesetRepository,
	roles repository.RoleRepository,
	entityRepositories ...any,
) *Service {
	s := &Service{users: users, auditLogs: auditLogs, rulesets: rulesets, roles: roles}
	for _, repo := range entityRepositories {
		if customerRepo, ok := repo.(repository.CustomerRepository); ok {
			s.customers = customerRepo
		}
		if cardRepo, ok := repo.(repository.CardRepository); ok {
			s.cards = cardRepo
		}
	}
	return s
}

func (s *Service) CreateUser(ctx context.Context, user *models.User, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), userPasswordHashCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
	return s.users.CreateUser(ctx, user)
}

func (s *Service) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.users.ListUsers(ctx)
}

func (s *Service) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.users.GetUser(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, id int, user *models.User, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), userPasswordHashCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)
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

func (s *Service) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	if customer.ID == "" {
		customer.ID = uuid.NewV4().String()
	}
	return s.customers.CreateCustomer(ctx, customer)
}

func (s *Service) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	return s.customers.ListCustomers(ctx)
}

func (s *Service) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	return s.customers.GetCustomer(ctx, id)
}

func (s *Service) UpdateCustomer(ctx context.Context, id string, customer *models.Customer) error {
	return s.customers.UpdateCustomer(ctx, id, customer)
}

func (s *Service) DeleteCustomer(ctx context.Context, id string) error {
	return s.customers.DeleteCustomer(ctx, id)
}

func (s *Service) CreateCard(ctx context.Context, card *models.Card) error {
	if card.ID == "" {
		card.ID = uuid.NewV4().String()
	}
	return s.cards.CreateCard(ctx, card)
}

func (s *Service) ListCards(ctx context.Context) ([]models.Card, error) {
	return s.cards.ListCards(ctx)
}

func (s *Service) GetCard(ctx context.Context, id string) (*models.Card, error) {
	return s.cards.GetCard(ctx, id)
}

func (s *Service) UpdateCard(ctx context.Context, id string, card *models.Card) error {
	return s.cards.UpdateCard(ctx, id, card)
}

func (s *Service) DeleteCard(ctx context.Context, id string) error {
	return s.cards.DeleteCard(ctx, id)
}

var _ UserService = (*Service)(nil)
var _ AuditLogService = (*Service)(nil)
var _ RulesetService = (*Service)(nil)
var _ RoleService = (*Service)(nil)
var _ CustomerService = (*Service)(nil)
var _ CardService = (*Service)(nil)
