package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/fuel"
)

type fuelService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewFuelService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &fuelService{repo: repo, logger: logger}
}

func (s *fuelService) Create(ctx context.Context, req domain.CreateFuelLogRequest, userID uuid.UUID) (*domain.FuelLog, error) {
	log := &domain.FuelLog{
		UserID:        userID,
		FuelType:      req.FuelType,
		AmountLiters:  req.AmountLiters,
		PricePerLiter: req.PricePerLiter,
		TotalCost:     req.AmountLiters * req.PricePerLiter,
		StationName:   req.StationName,
		OdometerKm:    req.OdometerKm,
		IsFullTank:    req.IsFullTank,
		LoggedAt:      time.Now(),
	}

	if req.MotorID != "" {
		mid, _ := uuid.Parse(req.MotorID)
		log.MotorID = &mid
	}

	if err := s.repo.Create(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create fuel log: %w", err)
	}

	return log, nil
}

func (s *fuelService) Get(ctx context.Context, id uuid.UUID) (*domain.FuelLog, error) {
	log, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("fuel log not found")
	}
	return log, nil
}

func (s *fuelService) List(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.FuelLog, string, bool, error) {
	return s.repo.FindByUserID(ctx, userID, cursor, limit)
}

func (s *fuelService) Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) (*domain.FuelLog, error) {
	log, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("fuel log not found")
	}
	if log.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	for k, v := range req {
		switch k {
		case "fuel_type":
			log.FuelType = v.(string)
		case "amount_liters":
			log.AmountLiters = v.(float64)
			log.TotalCost = log.AmountLiters * log.PricePerLiter
		case "price_per_liter":
			log.PricePerLiter = v.(float64)
			log.TotalCost = log.AmountLiters * log.PricePerLiter
		case "station_name":
			log.StationName = v.(string)
		case "odometer_km":
			log.OdometerKm = v.(float64)
		}
	}

	if err := s.repo.Update(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to update fuel log: %w", err)
	}

	return log, nil
}

func (s *fuelService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	log, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("fuel log not found")
	}
	if log.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.SoftDelete(ctx, id)
}
