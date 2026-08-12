package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type LocationBatch struct {
	mu        sync.Mutex
	locations []LocationData
}

type LocationData struct {
	UserID    string    `json:"user_id"`
	RoomID    string    `json:"room_id"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	Altitude  float64   `json:"altitude"`
	Accuracy  float64   `json:"accuracy"`
	Battery   int       `json:"battery"`
	Timestamp time.Time `json:"ts"`
}

type LocationAggregator struct {
	pool     *pgxpool.Pool
	nc       *nats.Conn
	batch    *LocationBatch
	flushInterval time.Duration
	batchSize     int
}

func NewLocationAggregator(pool *pgxpool.Pool, nc *nats.Conn) *LocationAggregator {
	return &LocationAggregator{
		pool:          pool,
		nc:            nc,
		batch:         &LocationBatch{},
		flushInterval: 5 * time.Second,
		batchSize:     100,
	}
}

func (a *LocationAggregator) Start(ctx context.Context) {
	go a.flushLoop(ctx)
	a.subscribeLocations(ctx)
}

func (a *LocationAggregator) subscribeLocations(ctx context.Context) {
	_, err := a.nc.Subscribe("room.>.location", func(msg *nats.Msg) {
		var loc LocationData
		if err := json.Unmarshal(msg.Data, &loc); err != nil {
			log.Printf("location unmarshal error: %v", err)
			return
		}
		loc.Timestamp = time.Now()

		a.batch.mu.Lock()
		a.batch.locations = append(a.batch.locations, loc)
		shouldFlush := len(a.batch.locations) >= a.batchSize
		a.batch.mu.Unlock()

		if shouldFlush {
			a.flush()
		}
	})
	if err != nil {
		log.Printf("location subscribe error: %v", err)
	}
	<-ctx.Done()
}

func (a *LocationAggregator) flushLoop(ctx context.Context) {
	ticker := time.NewTicker(a.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.flush()
		case <-ctx.Done():
			a.flush()
			return
		}
	}
}

func (a *LocationAggregator) flush() {
	a.batch.mu.Lock()
	locations := a.batch.locations
	a.batch.locations = nil
	a.batch.mu.Unlock()

	if len(locations) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := a.pool.Exec(ctx, `
		INSERT INTO location_tracking (user_id, room_id, lat, lon, speed, heading, altitude, accuracy, battery, recorded_at)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	`)
	if err != nil {
		log.Printf("location batch insert error: %v", err)
		return
	}

	log.Printf("flushed %d location records", len(locations))
}

func (a *LocationAggregator) bulkInsert(ctx context.Context, locations []LocationData) error {
	batch := &pgx.Batch{}
	for _, loc := range locations {
		batch.Queue(`
			INSERT INTO location_tracking (user_id, room_id, lat, lon, speed, heading, altitude, accuracy, battery, recorded_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, loc.UserID, loc.RoomID, loc.Lat, loc.Lon, loc.Speed, loc.Heading, loc.Altitude, loc.Accuracy, loc.Battery, loc.Timestamp)
	}

	br := a.pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(locations); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}
