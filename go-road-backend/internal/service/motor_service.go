package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/motor"
)

type motorService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewMotorService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &motorService{repo: repo, logger: logger}
}

func (s *motorService) Create(ctx context.Context, req domain.CreateMotorRequest, userID uuid.UUID) (*domain.Motor, error) {
	motor := &domain.Motor{
		UserID:       userID,
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         req.Year,
		LicensePlate: req.LicensePlate,
		EngineCC:     req.EngineCC,
		FuelType:     req.FuelType,
		TankCapacity: req.TankCapacity,
		PhotoURL:     req.PhotoURL,
		IsActive:     true,
	}

	motors, _ := s.repo.FindByUserID(ctx, userID)
	motor.IsPrimary = len(motors) == 0

	if err := s.repo.Create(ctx, motor); err != nil {
		return nil, fmt.Errorf("failed to create motor: %w", err)
	}

	return motor, nil
}

func (s *motorService) Get(ctx context.Context, id uuid.UUID) (*domain.Motor, error) {
	motor, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("motor not found")
	}
	return motor, nil
}

func (s *motorService) List(ctx context.Context, userID uuid.UUID) ([]domain.Motor, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *motorService) Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) (*domain.Motor, error) {
	motor, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("motor not found")
	}
	if motor.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	for k, v := range req {
		switch k {
		case "brand":
			motor.Brand = v.(string)
		case "model":
			motor.Model = v.(string)
		case "year":
			motor.Year = int(v.(float64))
		case "license_plate":
			motor.LicensePlate = v.(string)
		case "engine_cc":
			motor.EngineCC = int(v.(float64))
		case "fuel_type":
			motor.FuelType = v.(string)
		case "photo_url":
			motor.PhotoURL = v.(string)
		}
	}

	if err := s.repo.Update(ctx, motor); err != nil {
		return nil, fmt.Errorf("failed to update motor: %w", err)
	}

	return motor, nil
}

func (s *motorService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	motor, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("motor not found")
	}
	if motor.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *motorService) SetPrimary(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	motor, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("motor not found")
	}
	if motor.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.SetPrimary(ctx, userID, id)
}
