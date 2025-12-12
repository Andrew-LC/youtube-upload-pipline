package app

import (
	"context"
	"fmt"
	"io"
	"time"
	"github.com/Andrew-LC/uploader/internal/repository"
	"github.com/Andrew-LC/uploader/pkg/models"
)

const (
    KiB int64 = 1024
    MiB       = 1024 * KiB
    GiB       = 1024 * MiB
    TiB       = 1024 * GiB
)

type UploadService interface {
	ProcessUpload(ctx context.Context, fileName string, fileStream io.Reader, size int64, contentType string) (models.FileMetaData, error)
}

type uploadServiceImpl struct {
	repo repository.UploadRepository
	BucketName string
}

func NewUploadService(r repository.UploadRepository, bucket string) *uploadServiceImpl {
	return &uploadServiceImpl{repo: r, BucketName: bucket}
}

func (s *uploadServiceImpl) ProcessUpload(ctx context.Context, originalFileName string, fileStream io.Reader, size int64, contentType string) (models.FileMetaData, error) {
	uniqueName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), originalFileName)

	if size > 4 * GiB { 
		return models.FileMetaData{}, fmt.Errorf("file size exceeds 4GB limit")
	}

	metadata, err := s.repo.UploadFile(ctx, s.BucketName, uniqueName, fileStream, size, contentType)
	if err != nil {
		return models.FileMetaData{}, fmt.Errorf("failed to store file: %w", err)
	}

	return metadata, nil
}
