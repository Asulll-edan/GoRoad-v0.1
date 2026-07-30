package request

type CreateFuelLogRequest struct {
	MotorID  string  `json:"motor_id" validate:"required"`
	Liters   float64 `json:"liters" validate:"required,min=0"`
	Price    float64 `json:"price" validate:"required,min=0"`
	Odometer float64 `json:"odometer,omitempty"`
	FuelType string  `json:"fuel_type" validate:"oneof=pertalite pertamax pertamax_turbo solar"`
	Station  string  `json:"station,omitempty"`
}

type CreateExpenseRequest struct {
	RoomID      string  `json:"room_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,min=0"`
	Category    string  `json:"category" validate:"required,oneof=fuel food toll parking accommodation other"`
	Description string  `json:"description,omitempty" validate:"max=500"`
	SplitWith   []string `json:"split_with,omitempty"`
}
