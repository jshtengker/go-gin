package requests

type RegisterRequest struct {
	FullName    string  `json:"full_name" binding:"required,min=2,max=150"`
	Username    string  `json:"username" binding:"required,min=3,max=50"`
	Email       string  `json:"email" binding:"required,email,max=255"`
	Password    string  `json:"password" binding:"required,min=8,max=100"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,max=20"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Email or Username
	Password   string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
