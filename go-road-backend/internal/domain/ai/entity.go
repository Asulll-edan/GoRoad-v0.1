package ai

import "time"

type ChatRequest struct {
	UserID  string `json:"user_id"`
	RoomID  string `json:"room_id"`
	Message string `json:"message"`
}

type ChatResponse struct {
	Content string `json:"content"`
	IsFinal bool   `json:"is_final"`
}

type ItineraryRequest struct {
	RouteID      string   `json:"route_id"`
	StartLoc     string   `json:"start_location"`
	EndLoc       string   `json:"end_location"`
	DurationDays int      `json:"duration_days"`
	RiderCount   int      `json:"rider_count"`
	MotorIDs     []string `json:"motor_ids"`
	Preferences  []string `json:"preferences"`
}

type CostEstimate struct {
	RouteID      string  `json:"route_id"`
	TotalFuel    float64 `json:"total_fuel"`
	TotalFood    float64 `json:"total_food"`
	TotalHotel   float64 `json:"total_hotel"`
	TotalToll    float64 `json:"total_toll"`
	TotalParking float64 `json:"total_parking"`
	TotalAll     float64 `json:"total_all"`
	Breakdown    string  `json:"breakdown"`
}

type RouteAdvice struct {
	Alternatives []RouteOption `json:"alternatives"`
}

type RouteOption struct {
	Name        string   `json:"name"`
	DistanceKm  float64  `json:"distance_km"`
	Duration    string   `json:"duration"`
	Pros        []string `json:"pros"`
	Cons        []string `json:"cons"`
	RoadType    string   `json:"road_type"`
	Scenery     string   `json:"scenery"`
}
