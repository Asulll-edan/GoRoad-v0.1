package poi

import (
	"time"
	"github.com/google/uuid"
)

type POI struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description,omitempty"`
	Lat         float64   `json:"lat" gorm:"not null"`
	Lng         float64   `json:"lng" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null"`
	SubCategory string    `json:"sub_category,omitempty"`
	Address     string    `json:"address,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Website     string    `json:"website,omitempty"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	Rating      float64   `json:"rating,omitempty"`
	PriceLevel  int       `json:"price_level,omitempty"`
	OpenHours   string    `json:"open_hours,omitempty"`
	Source      string    `json:"source" gorm:"default:internal"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type POIReport struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `json:"user_id" gorm:"not null"`
	POIID       uuid.UUID `json:"poi_id" gorm:"not null"`
	ReportType  string    `json:"report_type" gorm:"not null"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
