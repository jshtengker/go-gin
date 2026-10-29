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
}

type OrderListItem struct {
	ID                uuid.UUID `gorm:"column:id"`
	UserID            uuid.UUID `gorm:"column:user_id"`
	UserFullName      string    `gorm:"column:user_full_name"`
	Username          string    `gorm:"column:username"`
	UserEmail         string    `gorm:"column:user_email"`
	ShippingAddressID uuid.UUID `gorm:"column:shipping_address_id"`
	Phone             string    `gorm:"column:phone"`
	AddressLine       string    `gorm:"column:address_line"`
	City              string    `gorm:"column:city"`
	Province          string    `gorm:"column:province"`
	PostalCode        string    `gorm:"column:postal_code"`
	Country           string    `gorm:"column:country"`
	OrderNumber       string    `gorm:"column:order_number"`
	Status            string    `gorm:"column:status"`
	TotalAmount       float64   `gorm:"column:total_amount"`
	OrderedAt         time.Time `gorm:"column:ordered_at"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

