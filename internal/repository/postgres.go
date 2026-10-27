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

var _ UserRepository = (*Repository)(nil)
var _ AuditLogRepository = (*Repository)(nil)
var _ RulesetRepository = (*Repository)(nil)
var _ RoleRepository = (*Repository)(nil)
