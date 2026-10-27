package models

import "time"

type Customer struct {
	ID          string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Status      string    `json:"status" gorm:"column:status;not null"`
	HomeCountry *string   `json:"home_country" gorm:"column:home_country;type:char(2)"`
	HomeCity    *string   `json:"home_city" gorm:"column:home_city"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Customer) TableName() string { return "customers" }
