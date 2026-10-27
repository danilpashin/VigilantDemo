package domain

import "time"

type UserDTO struct {
	Name     string `json:"name" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleID   int    `json:"role_id" binding:"required,gt=0"`
}

type AuditLogDTO struct {
	UserID     *int           `json:"user_id"`
	Action     string         `json:"action" binding:"required"`
	EntityType *string        `json:"entity_type"`
	EntityID   *string        `json:"entity_id"`
	Details    map[string]any `json:"details"`
}

type RulesetDTO struct {
	Version         int            `json:"version" binding:"required,gt=0"`
	ContentChecksum string         `json:"content_checksum" binding:"required"`
	ConfigSnapshot  map[string]any `json:"config_snapshot" binding:"required"`
	ActivatedAt     *time.Time     `json:"activated_at"`
	DeactivatedAt   *time.Time     `json:"deactivated_at"`
}

type CustomerDTO struct {
	ID          string  `json:"id"`
	Status      string  `json:"status" binding:"required,oneof=ACTIVE RESTRICTED CLOSED"`
	HomeCountry *string `json:"home_country" binding:"omitempty,len=2"`
	HomeCity    *string `json:"home_city"`
}

type CardDTO struct {
	ID         string     `json:"id"`
	CustomerID string     `json:"customer_id" binding:"required"`
	Token      string     `json:"token" binding:"required"`
	Last4      string     `json:"last4" binding:"required,len=4,numeric"`
	Status     string     `json:"status" binding:"required,oneof=ACTIVE BLOCKED EXPIRED CLOSED"`
	ExpiresAt  *time.Time `json:"expires_at"`
}
