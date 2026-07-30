package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email           string     `json:"email" gorm:"uniqueIndex;not null"`
	Phone           string     `json:"phone,omitempty" gorm:"uniqueIndex"`
	Username        string     `json:"username" gorm:"uniqueIndex;not null"`
	FullName        string     `json:"full_name" gorm:"not null"`
	PasswordHash    string     `json:"-" gorm:"not null"`
	PhotoURL        string     `json:"photo_url,omitempty"`
	Bio             string     `json:"bio,omitempty"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	Gender          string     `json:"gender,omitempty"`

	RidingSkill        string  `json:"riding_skill" gorm:"default:beginner"`
	RidingRole         string  `json:"riding_role" gorm:"default:member"`
	TotalPoints        int64   `json:"total_points" gorm:"default:0"`
	ExperienceYears    float64 `json:"experience_years" gorm:"default:0"`
	EmergencyContactName string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string `json:"emergency_contact_phone,omitempty"`
	MedicalNotes       string  `json:"medical_notes,omitempty"`
	BloodType          string  `json:"blood_type,omitempty"`

	Preferences      map[string]interface{} `json:"preferences" gorm:"type:jsonb;default:'{}'"`
	PrivacySettings  map[string]interface{} `json:"privacy_settings" gorm:"type:jsonb;default:'{\"show_location\":true,\"show_phone\":false,\"show_email\":false}'"`

	IsVerified  bool       `json:"is_verified" gorm:"default:false"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	IsBanned    bool       `json:"is_banned" gorm:"default:false"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type DeviceToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"not null"`
	Platform  string    `json:"platform" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"not null"`
	TokenHash string     `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	IsRevoked bool       `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	FullName string `json:"full_name" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Phone    string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	RidingSkill string    `json:"riding_skill"`
	RidingRole  string    `json:"riding_role"`
	IsVerified  bool      `json:"is_verified"`
}
