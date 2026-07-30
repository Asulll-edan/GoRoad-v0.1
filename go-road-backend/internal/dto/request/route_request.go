package request

type CreateRouteRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=200"`
	Description string  `json:"description,omitempty" validate:"max=2000"`
	DistanceKm  float64 `json:"distance_km" validate:"min=0"`
	DurationH   float64 `json:"duration_hours,omitempty"`
	ElevationG  float64 `json:"elevation_gain,omitempty"`
	Polyline    string  `json:"polyline,omitempty"`
	OriginLat   string  `json:"origin_lat,omitempty"`
	OriginLon   string  `json:"origin_lon,omitempty"`
	DestLat     string  `json:"dest_lat,omitempty"`
	DestLon     string  `json:"dest_lon,omitempty"`
	Waypoints   []WaypointRequest `json:"waypoints,omitempty"`
}

type WaypointRequest struct {
	Lat        float64 `json:"lat" validate:"required"`
	Lon        float64 `json:"lon" validate:"required"`
	Name       string  `json:"name,omitempty"`
	OrderIndex int     `json:"order_index"`
	StopDurationMin int `json:"stop_duration_min,omitempty"`
}

type UpdateRouteRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status,omitempty" validate:"oneof=draft active completed cancelled"`
}

type GPXImportRequest struct {
	GPXContent string `json:"gpx_content" validate:"required"`
}

type AddWaypointRequest struct {
	Lat float64 `json:"lat" validate:"required"`
	Lon float64 `json:"lon" validate:"required"`
	Name string `json:"name,omitempty"`
}
