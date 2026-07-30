package checklist

import (
	"time"
	"github.com/google/uuid"
)

type ChecklistTemplate struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedBy   uuid.UUID  `json:"created_by" gorm:"not null"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description,omitempty"`
	IsPublic    bool       `json:"is_public" gorm:"default:true"`
	Category    string     `json:"category" gorm:"default:general"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ChecklistItem struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TemplateID uuid.UUID `json:"template_id" gorm:"not null"`
	Label      string    `json:"label" gorm:"not null"`
	OrderIndex int       `json:"order_index" gorm:"not null"`
	IsRequired bool      `json:"is_required" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
}

type TouringChecklist struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID    uuid.UUID  `json:"room_id" gorm:"not null;uniqueIndex:idx_room_user_item"`
	UserID    uuid.UUID  `json:"user_id" gorm:"not null;uniqueIndex:idx_room_user_item"`
	ItemID    uuid.UUID  `json:"item_id" gorm:"not null;uniqueIndex:idx_room_user_item"`
	IsChecked bool       `json:"is_checked" gorm:"default:false"`
	CheckedAt *time.Time `json:"checked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
