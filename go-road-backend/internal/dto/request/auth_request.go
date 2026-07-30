package request

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6,max=100"`
	DisplayName string `json:"display_name,omitempty" validate:"max=100"`
	Phone       string `json:"phone,omitempty" validate:"max=20"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name,omitempty"`
	Bio         string `json:"bio,omitempty" validate:"max=500"`
	Phone       string `json:"phone,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

type RegisterDeviceRequest struct {
	Token    string `json:"token" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=ios android web"`
}
