package voting

import (
	"time"
	"github.com/google/uuid"
)

type Voting struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID      uuid.UUID  `json:"room_id" gorm:"not null"`
	CreatedBy   uuid.UUID  `json:"created_by" gorm:"not null"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description,omitempty"`
	VotingType  string     `json:"voting_type" gorm:"default:single"`
	Status      string     `json:"status" gorm:"default:active"`
	StartsAt    time.Time  `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
	IsAnonymous bool       `json:"is_anonymous" gorm:"default:false"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type VotingAnswer struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	VotingID   uuid.UUID `json:"voting_id" gorm:"not null"`
	Label      string    `json:"label" gorm:"not null"`
	OrderIndex int       `json:"order_index" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

type Vote struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	VotingID uuid.UUID `json:"voting_id" gorm:"not null"`
	AnswerID uuid.UUID `json:"answer_id" gorm:"not null"`
	UserID   uuid.UUID `json:"user_id" gorm:"not null"`
	Rank     int       `json:"rank" gorm:"default:1"`
	VotedAt  time.Time `json:"voted_at"`
}
