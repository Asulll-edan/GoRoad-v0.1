package request

type CreateEmergencyRequest struct {
	Type        string  `json:"type" validate:"required,oneof=accident breakdown medical weather lost other"`
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
	Description string  `json:"description,omitempty" validate:"max=500"`
	Severity    string  `json:"severity" validate:"required,oneof=low medium high critical"`
}

type SOSRequest struct {
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
	Message     string  `json:"message,omitempty"`
}

type UpdateEmergencyRequest struct {
	Status      string `json:"status" validate:"oneof=active acknowledged resolved false_alarm"`
	Description string `json:"description,omitempty"`
}
