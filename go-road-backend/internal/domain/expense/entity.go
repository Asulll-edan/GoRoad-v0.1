package expense

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"not null;index"`
	RoomID      *uuid.UUID `json:"room_id,omitempty"`
	Category    string     `json:"category" gorm:"not null"`
	Amount      float64    `json:"amount" gorm:"not null"`
	Description string     `json:"description,omitempty"`
	Location    string     `json:"location,omitempty" gorm:"type:geography(point,4326)"`
	ReceiptURL  string     `json:"receipt_url,omitempty"`
	IsSplitBill bool       `json:"is_split_bill" gorm:"default:false"`
	SplitWith   []string   `json:"split_with,omitempty" gorm:"type:uuid[]"`
	LoggedAt    time.Time  `json:"logged_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateExpenseRequest struct {
	Category    string   `json:"category" validate:"required"`
	Amount      float64  `json:"amount" validate:"required,gt=0"`
	Description string   `json:"description,omitempty"`
	IsSplitBill bool     `json:"is_split_bill,omitempty"`
	SplitWith   []string `json:"split_with,omitempty"`
}
