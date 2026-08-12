package worker

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	// Jakarta to Bandung great-circle distance ~118km
	jakartaLat, jakartaLon := -6.2, 106.8
	bandungLat, bandungLon := -6.9, 107.6

	dist := haversineDistance(jakartaLat, jakartaLon, bandungLat, bandungLon)
	if math.Abs(dist-118) > 5 {
		t.Errorf("expected ~118km, got %.2f", dist)
	}
}

func TestHaversineDistanceZero(t *testing.T) {
	dist := haversineDistance(-6.2, 106.8, -6.2, 106.8)
	if dist != 0 {
		t.Errorf("expected 0, got %.2f", dist)
	}
}

func TestPointToPolylineDistance(t *testing.T) {
	polyline := []struct{ Lat, Lon float64 }{
		{-6.2, 106.8},
		{-6.3, 106.9},
		{-6.4, 107.0},
	}

	dist := pointToPolylineDistance(-6.25, 106.85, polyline)
	if dist < 0 {
		t.Errorf("expected non-negative distance, got %.2f", dist)
	}
}
