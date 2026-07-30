package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/route"
	"go-road-backend/internal/repository/redis"
)

type routeService struct {
	repo   domain.Repository
	cache  redis.CacheRepository
	logger *zap.Logger
}

func NewRouteService(repo domain.Repository, cache redis.CacheRepository, logger *zap.Logger) domain.Service {
	return &routeService{repo: repo, cache: cache, logger: logger}
}

func (s *routeService) Create(ctx context.Context, req domain.CreateRouteRequest, userID uuid.UUID) (*domain.Route, error) {
	route := &domain.Route{
		CreatedBy:   userID,
		Name:        req.Name,
		IsPublic:    req.IsPublic,
		IsActive:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if len(req.Waypoints) > 0 {
		route.OriginLat = req.Waypoints[0].Lat
		route.OriginLng = req.Waypoints[0].Lng
		route.OriginName = req.Waypoints[0].Name

		last := req.Waypoints[len(req.Waypoints)-1]
		route.DestLat = last.Lat
		route.DestLng = last.Lng
		route.DestName = last.Name
	}

	if err := s.repo.Create(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	for _, wp := range req.Waypoints {
		waypoint := &domain.Waypoint{
			RouteID:    route.ID,
			Name:       wp.Name,
			Location:   fmt.Sprintf("POINT(%f %f)", wp.Lng, wp.Lat),
			OrderIndex: wp.OrderIndex,
			WaypointType: wp.WaypointType,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if waypoint.WaypointType == "" {
			waypoint.WaypointType = "stop"
		}
		if err := s.repo.AddWaypoint(ctx, waypoint); err != nil {
			return nil, fmt.Errorf("failed to add waypoint: %w", err)
		}
	}

	return route, nil
}

func (s *routeService) GetRoute(ctx context.Context, id uuid.UUID) (*domain.Route, error) {
	route, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("route not found")
	}
	return route, nil
}

func (s *routeService) Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) error {
	route, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("route not found")
	}
	if route.CreatedBy != userID {
		return errors.New("unauthorized")
	}

	for k, v := range req {
		switch k {
		case "name":
			route.Name = v.(string)
		case "description":
			route.Description = v.(string)
		case "is_public":
			route.IsPublic = v.(bool)
		}
	}

	route.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, route); err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}

	if route.RoomID != nil {
		s.cache.Delete(ctx, fmt.Sprintf("cache:route:active:%s", route.RoomID.String()))
	}

	return nil
}

func (s *routeService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	route, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("route not found")
	}
	if route.CreatedBy != userID {
		return errors.New("unauthorized")
	}

	return s.repo.SoftDelete(ctx, id)
}

func (s *routeService) AddWaypoint(ctx context.Context, routeID uuid.UUID, wp domain.WaypointInput) (*domain.Waypoint, error) {
	waypoint := &domain.Waypoint{
		RouteID:    routeID,
		Name:       wp.Name,
		Location:   fmt.Sprintf("POINT(%f %f)", wp.Lng, wp.Lat),
		OrderIndex: wp.OrderIndex,
		WaypointType: wp.WaypointType,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if waypoint.WaypointType == "" {
		waypoint.WaypointType = "stop"
	}

	if err := s.repo.AddWaypoint(ctx, waypoint); err != nil {
		return nil, fmt.Errorf("failed to add waypoint: %w", err)
	}

	return waypoint, nil
}

func (s *routeService) ListWaypoints(ctx context.Context, routeID uuid.UUID) ([]domain.Waypoint, error) {
	waypoints, err := s.repo.ListWaypoints(ctx, routeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list waypoints: %w", err)
	}
	return waypoints, nil
}

func (s *routeService) Activate(ctx context.Context, id, roomID uuid.UUID) error {
	route, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("route not found")
	}

	route.RoomID = &roomID

	routes, _ := s.repo.FindByRoomID(ctx, roomID)
	for _, r := range routes {
		if r.IsActive && r.ID != id {
			r.IsActive = false
			s.repo.Update(ctx, &r)
		}
	}

	route.IsActive = true
	route.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, route); err != nil {
		return fmt.Errorf("failed to activate route: %w", err)
	}

	activeData, _ := json.Marshal(route)
	s.cache.Set(ctx, fmt.Sprintf("cache:route:active:%s", roomID.String()), string(activeData), 30*time.Minute)

	return nil
}

func (s *routeService) ExportGPX(ctx context.Context, id uuid.UUID) (string, error) {
	route, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return "", errors.New("route not found")
	}

	waypoints, err := s.repo.ListWaypoints(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to list waypoints: %w", err)
	}

	gpx := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="GoRoad" xmlns="http://www.topografix.com/GPX/1/1">
  <metadata>
    <name>%s</name>
    <time>%s</time>
  </metadata>
  <rte>
    <name>%s</name>`, escapeXML(route.Name), route.CreatedAt.Format(time.RFC3339), escapeXML(route.Name))

	for _, wp := range waypoints {
		var lat, lng float64
		fmt.Sscanf(wp.Location, "POINT(%f %f)", &lng, &lat)
		gpx += fmt.Sprintf(`
    <rtept lat="%f" lon="%f">
      <name>%s</name>
    </rtept>`, lat, lng, escapeXML(wp.Name))
	}

	gpx += `
  </rte>
</gpx>`

	return gpx, nil
}

func escapeXML(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '&':
			result += "&amp;"
		case '<':
			result += "&lt;"
		case '>':
			result += "&gt;"
		case '"':
			result += "&quot;"
		case '\'':
			result += "&apos;"
		default:
			result += string(c)
		}
	}
	return result
}
