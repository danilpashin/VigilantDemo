package domain

import "time"

type UserDTO struct {
	Name         string `json:"name" binding:"required"`
	Login        string `json:"login" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	RoleID       int    `json:"role_id" binding:"required,gt=0"`
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
