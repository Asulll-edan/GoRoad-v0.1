package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/qr"
)

type qrService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewQRService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &qrService{repo: repo, logger: logger}
}

func (s *qrService) GetMyQRCard(ctx context.Context, userID uuid.UUID) (*domain.QRCard, error) {
	existing, err := s.repo.FindByUserID(ctx, userID)
	if err == nil {
		return existing, nil
	}
	code, _ := generateQRCode()
	card := &domain.QRCard{
		UserID:   userID,
		Code:     code,
		Style:    "default",
		IsActive: true,
	}
	if err := s.repo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("failed to create QR card: %w", err)
	}
	return card, nil
}

func (s *qrService) RegenerateQR(ctx context.Context, userID uuid.UUID) (*domain.QRCard, error) {
	existing, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return s.GetMyQRCard(ctx, userID)
	}
	code, _ := generateQRCode()
	existing.Code = code
	return existing, nil
}

func (s *qrService) ScanQR(ctx context.Context, code string) (*domain.QRCard, error) {
	card, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("QR code not found")
	}
	s.repo.IncrementScanCount(ctx, card.ID)
	return card, nil
}

func generateQRCode() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
