package fuel

import (
	"time"

	"github.com/google/uuid"
)

type FuelLog struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID  `json:"user_id" gorm:"not null;index"`
	MotorID       *uuid.UUID `json:"motor_id,omitempty"`
	RoomID        *uuid.UUID `json:"room_id,omitempty"`
	FuelType      string     `json:"fuel_type" gorm:"not null"`
	AmountLiters  float64    `json:"amount_liters" gorm:"not null"`
	PricePerLiter float64    `json:"price_per_liter" gorm:"not null"`
	TotalCost     float64    `json:"total_cost" gorm:"not null"`
	StationName   string     `json:"station_name,omitempty"`
	Location      string     `json:"location,omitempty" gorm:"type:geography(point,4326)"`
	OdometerKm    float64    `json:"odometer_km,omitempty"`
	IsFullTank    bool       `json:"is_full_tank" gorm:"default:false"`
	LoggedAt      time.Time  `json:"logged_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateFuelLogRequest struct {
	MotorID       string  `json:"motor_id" validate:"required"`
	FuelType      string  `json:"fuel_type" validate:"required"`
	AmountLiters  float64 `json:"amount_liters" validate:"required,gt=0"`
	PricePerLiter float64 `json:"price_per_liter" validate:"required,gt=0"`
	StationName   string  `json:"station_name,omitempty"`
	OdometerKm    float64 `json:"odometer_km,omitempty"`
	IsFullTank    bool    `json:"is_full_tank,omitempty"`
}
