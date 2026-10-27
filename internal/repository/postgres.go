package repository

import (
	"context"

	"vigilant/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *Repository) GetUser(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, user *models.User) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":          user.Name,
		"login":         user.Login,
		"password_hash": user.PasswordHash,
		"role_id":       user.RoleID,
		"updated_at":    gorm.Expr("now()"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) CreateAuditLog(ctx context.Context, auditLog *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(auditLog).Error
}

func (r *Repository) ListAuditLogs(ctx context.Context) ([]models.AuditLog, error) {
	var auditLogs []models.AuditLog
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&auditLogs).Error
	return auditLogs, err
}

func (r *Repository) CreateRuleset(ctx context.Context, ruleset *models.Ruleset) error {
	return r.db.WithContext(ctx).Create(ruleset).Error
}

func (r *Repository) ListRulesets(ctx context.Context) ([]models.Ruleset, error) {
	var rulesets []models.Ruleset
	err := r.db.WithContext(ctx).Order("version DESC").Find(&rulesets).Error
	return rulesets, err
}

func (r *Repository) ListRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Order("id").Find(&roles).Error
	return roles, err
}

func (r *Repository) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *Repository) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	var customers []models.Customer
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&customers).Error
	return customers, err
}

func (r *Repository) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *Repository) UpdateCustomer(ctx context.Context, id string, customer *models.Customer) error {
	result := r.db.WithContext(ctx).Model(&models.Customer{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       customer.Status,
		"home_country": customer.HomeCountry,
		"home_city":    customer.HomeCity,
		"updated_at":   gorm.Expr("now()"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) DeleteCustomer(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Customer{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) CreateCard(ctx context.Context, card *models.Card) error {
	return r.db.WithContext(ctx).Create(card).Error
}

func (r *Repository) ListCards(ctx context.Context) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&cards).Error
	return cards, err
}

func (r *Repository) GetCard(ctx context.Context, id string) (*models.Card, error) {
	var card models.Card
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *Repository) UpdateCard(ctx context.Context, id string, card *models.Card) error {
	result := r.db.WithContext(ctx).Model(&models.Card{}).Where("id = ?", id).Updates(map[string]interface{}{
		"customer_id": card.CustomerID,
		"token":       card.Token,
		"last4":       card.Last4,
		"status":      card.Status,
		"expires_at":  card.ExpiresAt,
		"updated_at":  gorm.Expr("now()"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) DeleteCard(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Card{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

var _ UserRepository = (*Repository)(nil)
var _ AuditLogRepository = (*Repository)(nil)
var _ RulesetRepository = (*Repository)(nil)
var _ RoleRepository = (*Repository)(nil)
var _ CustomerRepository = (*Repository)(nil)
var _ CardRepository = (*Repository)(nil)
