package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string         `gorm:"size:255;unique;not null"`
	Username     string         `gorm:"size:50;unique;not null"`
	FullName     string         `gorm:"column:full_name;size:150;not null"`
	PasswordHash string         `gorm:"column:password_hash;not null"`
	PhoneNumber  *string        `gorm:"size:20"`
	IsVerified   bool           `gorm:"column:is_verified;default:false"`
	IsActive     bool           `gorm:"column:is_active;default:true"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}
