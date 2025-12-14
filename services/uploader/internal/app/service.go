package app

import (
	"context"
	"fmt"
	"io"
	"time"
	"github.com/Andrew-LC/uploader/internal/repository"
	"github.com/Andrew-LC/uploader/pkg/models"
	"github.com/Andrew-LC/libs/mq"
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
	repo       repository.UploadRepository
	BucketName string
	mq         *mq.RabbitMQ
}

func NewUploadService(r repository.UploadRepository, bucket string, mq *mq.RabbitMQ) *uploadServiceImpl {
	return &uploadServiceImpl{repo: r, BucketName: bucket, mq: mq}
}

func (s *uploadServiceImpl) ProcessUpload(ctx context.Context, originalFileName string, fileStream io.Reader, size int64, contentType string) (models.FileMetaData, error) {
	uniqueName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), originalFileName)

	if size > 4*GiB {
		return models.FileMetaData{}, fmt.Errorf("file size exceeds 4GB limit")
	}

	metadata, err := s.repo.UploadFile(ctx, s.BucketName, uniqueName, fileStream, size, contentType)
	if err != nil {
		return models.FileMetaData{}, fmt.Errorf("failed to store file: %w", err)
	}

	// Define the exchange
	_ = s.mq.DeclareExchange("upload_events", "direct")

	// Publish event to RabbitMQ
	err = s.mq.PublishJSON(ctx, "upload_events", "upload.created", metadata)
	if err != nil {
		fmt.Printf("Failed to publish upload event: %v\n", err)
	}

	return metadata, nil
}
