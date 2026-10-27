package models

import "time"

type AuditLog struct {
	ID         int            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID     *int           `json:"user_id" gorm:"column:user_id"`
	Action     string         `json:"action" gorm:"column:action;not null"`
	EntityType *string        `json:"entity_type" gorm:"column:entity_type"`
	EntityID   *string        `json:"entity_id" gorm:"column:entity_id;type:uuid"`
	Details    map[string]any `json:"details" gorm:"column:details;type:jsonb;serializer:json"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at"`
}
