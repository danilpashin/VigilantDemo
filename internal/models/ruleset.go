package models

import "time"

type Ruleset struct {
	Version         int            `json:"version" gorm:"column:version;primaryKey;autoIncrement:false"`
	ContentChecksum string         `json:"content_checksum" gorm:"column:content_checksum;not null"`
	ConfigSnapshot  map[string]any `json:"config_snapshot" gorm:"column:config_snapshot;type:jsonb;serializer:json"`
	ActivatedAt     *time.Time     `json:"activated_at" gorm:"column:activated_at"`
	DeactivatedAt   *time.Time     `json:"deactivated_at" gorm:"column:deactivated_at"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at"`
}
