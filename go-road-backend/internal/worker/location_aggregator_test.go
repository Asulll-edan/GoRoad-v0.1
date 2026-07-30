package worker

import (
	"testing"
	"time"
)

func TestNewLocationAggregator(t *testing.T) {
	ag := &LocationAggregator{
		flushInterval: 5 * time.Second,
		batchSize:     100,
		batch:         &LocationBatch{},
	}
	if ag == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestLocationBatchAdd(t *testing.T) {
	batch := &LocationBatch{}
	loc := LocationData{
		UserID: "user1",
		RoomID: "room1",
		Lat:    -6.2,
		Lon:    106.8,
		Speed:  60,
	}

	batch.mu.Lock()
	batch.locations = append(batch.locations, loc)
	batch.mu.Unlock()

	if len(batch.locations) != 1 {
		t.Errorf("expected 1 location, got %d", len(batch.locations))
	}
}
