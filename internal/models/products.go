package models

import (
	"time"

	"github.com/google/uuid"
)

type Products struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	CategoryID  uuid.UUID `gorm:"primaryKey"`
	Sku         string
	Name        string
	Description string
	Price       string
	Stock       int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
