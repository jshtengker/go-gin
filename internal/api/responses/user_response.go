package responses

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber *string   `json:"phone"`
	IsVerified  bool      `json:"is_verified"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
