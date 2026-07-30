package ai

import "context"

type Service interface {
	Chat(ctx context.Context, userID, roomID, message string) (chan ChatResponse, error)
	GenerateItinerary(ctx context.Context, req ItineraryRequest) (*string, error)
	EstimateCost(ctx context.Context, routeID string, motorIDs []string, riderCount, durationDays int, fuelType string) (*string, error)
	AdviseRoute(ctx context.Context, origin, destination string, waypoints, preferences []string) (*string, error)
	GeneratePackingList(ctx context.Context, durationDays int, weatherCondition, touringType string) (*string, error)
	AdviseSafety(ctx context.Context, routeID, weatherCondition, skillLevel string, riderCount int) (*string, error)
	RecommendPOI(ctx context.Context, lat, lng, radiusKm float64, types []string) (*string, error)
}
