package motor

import (
	"time"

	"github.com/google/uuid"
)

type Motor struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID           uuid.UUID  `json:"user_id" gorm:"not null;index"`
	Brand            string     `json:"brand" gorm:"not null"`
	Model            string     `json:"model" gorm:"not null"`
	Year             int        `json:"year" gorm:"not null"`
	LicensePlate     string     `json:"license_plate,omitempty"`
	VIN              []byte     `json:"-" gorm:"type:bytea"`
	InsuranceInfo    []byte     `json:"-" gorm:"type:bytea"`
	STNKNumber       []byte     `json:"-" gorm:"type:bytea"`
	EngineCC         int        `json:"engine_cc,omitempty"`
	FuelType         string     `json:"fuel_type,omitempty"`
	TankCapacity     float64    `json:"tank_capacity,omitempty"`
	FuelConsumption  float64    `json:"fuel_consumption,omitempty"`
	TirePressureFront float64   `json:"tire_pressure_front,omitempty"`
	TirePressureRear float64    `json:"tire_pressure_rear,omitempty"`
	PhotoURL         string     `json:"photo_url,omitempty"`
	IsPrimary        bool       `json:"is_primary" gorm:"default:false"`
	IsActive         bool       `json:"is_active" gorm:"default:true"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type CreateMotorRequest struct {
	Brand        string  `json:"brand" validate:"required"`
	Model        string  `json:"model" validate:"required"`
	Year         int     `json:"year" validate:"required,gte=2000,lte=2030"`
	LicensePlate string  `json:"license_plate,omitempty"`
	EngineCC     int     `json:"engine_cc,omitempty"`
	FuelType     string  `json:"fuel_type,omitempty"`
	TankCapacity float64 `json:"tank_capacity,omitempty"`
	PhotoURL     string  `json:"photo_url,omitempty"`
}
