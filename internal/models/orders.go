package models

import (
	"time"

	"github.com/google/uuid"
)

type Orders struct {
	ID                uuid.UUID `gorm:"primaryKey"`
	UserID            uuid.UUID
	ShippingAddressID uuid.UUID
	OrderNumber       string
	Status            string
	TotalAmount       float64
	OrderedAt         time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	User              User `gorm:"foreignKey:UserID;references:ID"`
	// ShippingAddress   Address `gorm:"foreignKey:ShippingAddressID;references:ID"`
}
