package route

import (
	"time"

	"github.com/google/uuid"
)

type Route struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID           *uuid.UUID `json:"room_id,omitempty"`
	CreatedBy        uuid.UUID  `json:"created_by" gorm:"not null"`
	Name             string     `json:"name" gorm:"not null"`
	Description      string     `json:"description,omitempty"`
	RouteGeom        string     `json:"route_geom,omitempty" gorm:"type:geography(linestring,4326)"`
	DistanceKm       float64    `json:"distance_km,omitempty"`
	EstimatedDuration string    `json:"estimated_duration,omitempty"`
	ElevationGain    int        `json:"elevation_gain,omitempty"`
	ElevationLoss    int        `json:"elevation_loss,omitempty"`
	MaxElevation     int        `json:"max_elevation,omitempty"`
	MinElevation     int        `json:"min_elevation,omitempty"`
	OriginLat        float64    `json:"origin_lat,omitempty"`
	OriginLng        float64    `json:"origin_lng,omitempty"`
	OriginName       string     `json:"origin_name,omitempty"`
	DestLat          float64    `json:"dest_lat,omitempty"`
	DestLng          float64    `json:"dest_lng,omitempty"`
	DestName         string     `json:"dest_name,omitempty"`
	Polyline         string     `json:"polyline,omitempty"`
	IsActive         bool       `json:"is_active" gorm:"default:false"`
	IsPublic         bool       `json:"is_public" gorm:"default:true"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type Waypoint struct {
	ID                 uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RouteID            uuid.UUID  `json:"route_id" gorm:"not null"`
	Name               string     `json:"name,omitempty"`
	Description        string     `json:"description,omitempty"`
	Location           string     `json:"location" gorm:"type:geography(point,4326);not null"`
	OrderIndex         int        `json:"order_index" gorm:"not null"`
	WaypointType       string     `json:"waypoint_type" gorm:"default:stop"`
	EstimatedArrival   *time.Time `json:"estimated_arrival,omitempty"`
	EstimatedDeparture *time.Time `json:"estimated_departure,omitempty"`
	IsPOI              bool       `json:"is_poi" gorm:"default:false"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type CreateRouteRequest struct {
	Name       string           `json:"name" validate:"required"`
	Waypoints  []WaypointInput  `json:"waypoints" validate:"required,min=2"`
	IsPublic   bool             `json:"is_public,omitempty"`
}

type WaypointInput struct {
	Name         string  `json:"name"`
	Lat          float64 `json:"lat" validate:"required"`
	Lng          float64 `json:"lng" validate:"required"`
	WaypointType string  `json:"waypoint_type,omitempty"`
	OrderIndex   int     `json:"order_index"`
}
