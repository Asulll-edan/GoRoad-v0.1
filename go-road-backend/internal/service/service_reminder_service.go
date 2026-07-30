package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/service_reminder"
)

type serviceReminderService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewServiceReminderService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &serviceReminderService{repo: repo, logger: logger}
}

func (s *serviceReminderService) Create(ctx context.Context, req domain.CreateReminderRequest, userID uuid.UUID) (*domain.ServiceReminder, error) {
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, errors.New("invalid due_date format, use YYYY-MM-DD")
	}

	motorID, _ := uuid.Parse(req.MotorID)
	r := &domain.ServiceReminder{
		UserID:                userID,
		MotorID:               motorID,
		ServiceType:           req.ServiceType,
		Title:                 req.Title,
		Description:           req.Description,
		DueDate:               dueDate,
		DueOdometer:           req.DueOdometer,
		IsRecurring:           req.IsRecurring,
		RecurringIntervalDays: req.RecurringIntervalDays,
	}
	if err := s.repo.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}
	return r, nil
}

func (s *serviceReminderService) List(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]domain.ServiceReminder, string, bool, error) {
	return s.repo.FindByUserID(ctx, userID, cursor, limit)
}

func (s *serviceReminderService) Get(ctx context.Context, id uuid.UUID) (*domain.ServiceReminder, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *serviceReminderService) Update(ctx context.Context, id uuid.UUID, req map[string]interface{}, userID uuid.UUID) error {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("reminder not found")
	}
	if r.UserID != userID {
		return errors.New("unauthorized")
	}
	for k, v := range req {
		switch k {
		case "title":
			r.Title = v.(string)
		case "description":
			r.Description = v.(string)
		case "due_date":
			r.DueDate, _ = time.Parse("2006-01-02", v.(string))
		case "due_odometer":
			r.DueOdometer = v.(float64)
		case "service_type":
			r.ServiceType = v.(string)
		}
	}
	return s.repo.Update(ctx, r)
}

func (s *serviceReminderService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("reminder not found")
	}
	if r.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *serviceReminderService) Complete(ctx context.Context, id, userID uuid.UUID) error {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("reminder not found")
	}
	if r.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Complete(ctx, id)
}
