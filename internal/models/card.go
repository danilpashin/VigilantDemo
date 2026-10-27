package models

import "time"

type Card struct {
	ID         string     `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	CustomerID string     `json:"customer_id" gorm:"column:customer_id;type:uuid;not null"`
	Token      string     `json:"token" gorm:"column:token;not null;uniqueIndex"`
	Last4      string     `json:"last4" gorm:"column:last4;type:char(4);not null"`
	Status     string     `json:"status" gorm:"column:status;not null"`
	ExpiresAt  *time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

func (Card) TableName() string { return "cards" }
