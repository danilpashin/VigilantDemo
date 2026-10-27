package models

import "time"

type User struct {
	ID           int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"column:name;not null"`
	Login        string    `json:"login" gorm:"column:login;not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	RoleID       int       `json:"role_id" gorm:"column:role_id;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}
