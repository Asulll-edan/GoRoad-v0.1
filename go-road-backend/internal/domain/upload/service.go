package upload

import (
	"context"
	"mime/multipart"
	"github.com/google/uuid"
)

type Service interface {
	UploadFile(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader, category string) (*Upload, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*Upload, error)
	UploadPhoto(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*Upload, error)
	DeleteFile(ctx context.Context, id, userID uuid.UUID) error
}
