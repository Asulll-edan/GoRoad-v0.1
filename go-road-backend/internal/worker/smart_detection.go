package worker

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/nats-io/nats.go"
	"go-road-backend/internal/repository/redis"
)

type SmartDetection struct {
	nc          *nats.Conn
	cache       redis.CacheRepository
	ctx         context.Context
}

type RiderPosition struct {
	UserID  string  `json:"user_id"`
	RoomID  string  `json:"room_id"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Speed   float64 `json:"speed"`
	Battery int     `json:"battery"`
	TS      int64   `json:"ts"`
}

func NewSmartDetection(nc *nats.Conn, cache redis.CacheRepository) *SmartDetection {
	return &SmartDetection{
		nc:    nc,
		cache: cache,
		ctx:   context.Background(),
	}
}

func (d *SmartDetection) Start(ctx context.Context) {
	d.ctx = ctx
	_, err := d.nc.Subscribe("room.>.location", func(msg *nats.Msg) {
		var pos RiderPosition
		if err := json.Unmarshal(msg.Data, &pos); err != nil {
			return
		}
		d.evaluatePosition(pos)
	})
	if err != nil {
		log.Printf("smart detection subscribe error: %v", err)
	}
	<-ctx.Done()
}

func (d *SmartDetection) evaluatePosition(pos RiderPosition) {
	now := time.Now().Unix()

	cooldownKey := "detect:cooldown:" + pos.RoomID + ":" + pos.UserID
	onCooldown, err := d.cache.Exists(d.ctx, cooldownKey)
	if err == nil && onCooldown {
		return
	}

	alerts := make([]map[string]interface{}, 0)

	if pos.Speed > 130 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "speed_limit",
			"user_id": pos.UserID,
			"room_id": pos.RoomID,
			"speed":   pos.Speed,
		})
		d.cache.Set(d.ctx, cooldownKey, "1", 60*time.Second)
	}

	if pos.Battery > 0 && pos.Battery < 15 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "battery_low",
			"user_id": pos.UserID,
			"room_id": pos.RoomID,
			"battery": pos.Battery,
		})
		d.cache.Set(d.ctx, cooldownKey, "1", 300*time.Second)
	}

	age := now - pos.TS
	if age > 300 && pos.Speed < 1 {
		alerts = append(alerts, map[string]interface{}{
			"type":         "stopped_long",
			"user_id":      pos.UserID,
			"room_id":      pos.RoomID,
			"duration_sec": int(age),
		})
		d.cache.Set(d.ctx, cooldownKey, "1", 300*time.Second)
	}

	if age > 120 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "offline",
			"user_id": pos.UserID,
			"room_id": pos.RoomID,
		})
		d.cache.Set(d.ctx, cooldownKey, "1", 60*time.Second)
	}

	for _, alert := range alerts {
		alertData, _ := json.Marshal(alert)
		d.nc.Publish("detection.alert", alertData)
		d.cache.Publish(d.ctx, "detection.alert", alertData)
	}
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c * 6371
}

func pointToPolylineDistance(lat, lon float64, polyline []struct{ Lat, Lon float64 }) float64 {
	minDist := math.MaxFloat64
	for i := 0; i < len(polyline)-1; i++ {
		dist := distanceToSegment(lat, lon, polyline[i].Lat, polyline[i].Lon, polyline[i+1].Lat, polyline[i+1].Lon)
		if dist < minDist {
			minDist = dist
		}
	}
	return minDist
}

func distanceToSegment(px, py, ax, ay, bx, by float64) float64 {
	dx := bx - ax
	dy := by - ay
	if dx == 0 && dy == 0 {
		return haversineDistance(px, py, ax, ay)
	}
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	if t < 0 {
		return haversineDistance(px, py, ax, ay)
	}
	if t > 1 {
		return haversineDistance(px, py, bx, by)
	}
	projX := ax + t*dx
	projY := ay + t*dy
	return haversineDistance(px, py, projX, projY)
}
