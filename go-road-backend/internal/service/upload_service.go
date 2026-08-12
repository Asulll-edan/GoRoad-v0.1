package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	domain "go-road-backend/internal/domain/upload"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

var maxFileSize int64 = 10 * 1024 * 1024

type uploadService struct {
	repo   domain.Repository
	logger *zap.Logger
}

func NewUploadService(repo domain.Repository, logger *zap.Logger) domain.Service {
	return &uploadService{repo: repo, logger: logger}
}

func (s *uploadService) UploadFile(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader, category string) (*domain.Upload, error) {
	if header.Size > maxFileSize {
		return nil, fmt.Errorf("file too large, max %d bytes", maxFileSize)
	}
	mimeType := header.Header.Get("Content-Type")
	if !allowedImageTypes[mimeType] {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectKey := fmt.Sprintf("uploads/%s/%s", userID, filename)

	_ = data

	upload := &domain.Upload{
		UserID:    userID,
		FileName:  header.Filename,
		FileSize:  header.Size,
		MimeType:  mimeType,
		URL:       fmt.Sprintf("/static/%s", objectKey),
		Bucket:    "goroad",
		ObjectKey: objectKey,
		Category:  category,
		IsPublic:  true,
	}
	if err := s.repo.Create(ctx, upload); err != nil {
		return nil, fmt.Errorf("failed to save upload metadata: %w", err)
	}
	return upload, nil
}

func (s *uploadService) UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*domain.Upload, error) {
	return s.UploadFile(ctx, userID, file, header, "avatar")
}

func (s *uploadService) UploadPhoto(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*domain.Upload, error) {
	return s.UploadFile(ctx, userID, file, header, "photo")
}

func (s *uploadService) DeleteFile(ctx context.Context, id, userID uuid.UUID) error {
	upload, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("file not found")
	}
	if upload.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	_ = strings.TrimSpace(upload.ObjectKey)
	return s.repo.SoftDelete(ctx, id)
}
