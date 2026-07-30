package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/ai"
	"go-road-backend/internal/repository/redis"
)

type aiService struct {
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewAIService(cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &aiService{cache: cache, logger: logger}
}

func (s *aiService) Chat(ctx context.Context, userID, roomID, message string) (chan domain.ChatResponse, error) {
	allowed, err := s.checkRateLimit(ctx, userID)
	if err != nil || !allowed {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	ch := make(chan domain.ChatResponse, 10)
	go func() {
		defer close(ch)
		ch <- domain.ChatResponse{
			Content: "AI chat akan diimplementasikan via Python gRPC service",
			IsFinal: true,
		}
	}()
	return ch, nil
}

func (s *aiService) GenerateItinerary(ctx context.Context, req domain.ItineraryRequest) (*string, error) {
	hash := s.hashParams(req)
	cacheKey := fmt.Sprintf("cache:ai:itinerary:%s", hash)
	var cached string
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}
	result := `{"days":[],"total_distance_km":0,"total_duration_hours":0,"recommendations":[]}`
	s.cache.SetJSON(ctx, cacheKey, result, 6*time.Hour)
	return &result, nil
}

func (s *aiService) EstimateCost(ctx context.Context, routeID string, motorIDs []string, riderCount, durationDays int, fuelType string) (*string, error) {
	hash := s.hashString(routeID + fmt.Sprintf("%v", motorIDs) + fmt.Sprintf("%d", durationDays))
	cacheKey := fmt.Sprintf("cache:ai:cost:%s", hash)
	var cached string
	if err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}
	result := `{"fuel_cost":0,"food_cost":0,"accommodation_cost":0,"total_cost":0,"per_person_cost":0}`
	s.cache.SetJSON(ctx, cacheKey, result, 6*time.Hour)
	return &result, nil
}

func (s *aiService) AdviseRoute(ctx context.Context, origin, destination string, waypoints, preferences []string) (*string, error) {
	result := `{"alternatives":[],"recommended_index":0,"tips":[]}`
	return &result, nil
}

func (s *aiService) GeneratePackingList(ctx context.Context, durationDays int, weatherCondition, touringType string) (*string, error) {
	result := `{"categories":[{"name":"Dokumen","items":["SIM","STNK","KTP"]},{"name":"Perlengkapan","items":["Helm","Jaket","Sarung Tangan"]}]}`
	return &result, nil
}

func (s *aiService) AdviseSafety(ctx context.Context, routeID, weatherCondition, skillLevel string, riderCount int) (*string, error) {
	result := `{"advice":["Periksa kondisi motor sebelum berangkat","Gunakan perlengkapan safety lengkap","Jaga jarak aman"]}`
	return &result, nil
}

func (s *aiService) RecommendPOI(ctx context.Context, lat, lng, radiusKm float64, types []string) (*string, error) {
	result := `{"pois":[]}`
	return &result, nil
}

func (s *aiService) checkRateLimit(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("rate:ai:%s", userID)
	count, err := s.cache.IncrWithExpiry(ctx, key, 1*time.Hour)
	if err != nil {
		return true, nil
	}
	return count <= 20, nil
}

func (s *aiService) hashParams(v interface{}) string {
	data, _ := json.Marshal(v)
	return s.hashString(string(data))
}

func (s *aiService) hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:16])
}
