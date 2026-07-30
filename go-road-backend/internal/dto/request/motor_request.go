package request

type CreateMotorRequest struct {
	Name         string  `json:"name" validate:"required,max=100"`
	Brand        string  `json:"brand" validate:"required,max=50"`
	Model        string  `json:"model,omitempty"`
	Year         int     `json:"year" validate:"min=1990,max=2030"`
	PlateNumber  string  `json:"plate_number,omitempty"`
	Color        string  `json:"color,omitempty"`
	EngineCc     float64 `json:"engine_cc,omitempty"`
	FuelType     string  `json:"fuel_type,omitempty" validate:"oneof=pertalite pertamax solar listrik"`
	TankCapacity float64 `json:"tank_capacity,omitempty"`
	IsPrimary    bool    `json:"is_primary"`
}

type UpdateMotorRequest struct {
	Name         string  `json:"name,omitempty"`
	PlateNumber  string  `json:"plate_number,omitempty"`
	Color        string  `json:"color,omitempty"`
	IsPrimary    *bool   `json:"is_primary,omitempty"`
}
