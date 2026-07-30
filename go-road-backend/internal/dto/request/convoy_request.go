package request

type CreateFormationRequest struct {
	Name           string `json:"name" validate:"required,max=100"`
	FormationType  string `json:"formation_type" validate:"required,oneof=staggered single_column double_column"`
	MaxDistanceM   int    `json:"max_distance_m" validate:"min=10,max=500"`
	MaxSpeedKmh    int    `json:"max_speed_kmh" validate:"min=20,max=200"`
	LeadUserID     string `json:"lead_user_id" validate:"required"`
	SweepUserID    string `json:"sweep_user_id,omitempty"`
}

type UpdateFormationRequest struct {
	Name          string `json:"name,omitempty" validate:"max=100"`
	FormationType string `json:"formation_type,omitempty" validate:"oneof=staggered single_column double_column"`
	MaxDistanceM  int    `json:"max_distance_m,omitempty"`
	MaxSpeedKmh   int    `json:"max_speed_kmh,omitempty"`
}

type LocationUpdateRequest struct {
	Lat      float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lon      float64 `json:"lon" validate:"required,min=-180,max=180"`
	Speed    float64 `json:"speed,omitempty"`
	Heading  float64 `json:"heading,omitempty" validate:"min=0,max=360"`
	Altitude float64 `json:"altitude,omitempty"`
	Accuracy float64 `json:"accuracy,omitempty"`
	Battery  int     `json:"battery,omitempty" validate:"min=0,max=100"`
}
