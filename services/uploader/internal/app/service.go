package app

import (
	"context"
	"fmt"
	"io"

	"yup/Andrew-LC/libs/logger"
	"yup/Andrew-LC/libs/models"
	"yup/Andrew-LC/libs/mq"
	"github.com/Andrew-LC/uploader/internal/repository"
	"go.uber.org/zap"
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
	logger     *logger.Logger
}

func NewUploadService(r repository.UploadRepository, bucket string, mq *mq.RabbitMQ, l *logger.Logger) *uploadServiceImpl {
	return &uploadServiceImpl{repo: r, BucketName: bucket, mq: mq, logger: l}
}

func (s *uploadServiceImpl) ProcessUpload(ctx context.Context, originalFileName string, fileStream io.Reader, size int64, contentType string) (models.FileMetaData, error) {
	zapLog := s.logger.GetZapLogger()

	if size > 4*GiB {
		return models.FileMetaData{}, fmt.Errorf("file size exceeds 4GB limit")
	}

	metadata, err := s.repo.UploadFile(ctx, s.BucketName, originalFileName, fileStream, size, contentType)
	if err != nil {
		zapLog.Error("Failed to store file", zap.Error(err))
		return models.FileMetaData{}, fmt.Errorf("failed to store file: %w", err)
	}

	err = s.mq.DeclareExchange("upload_events", "direct")
	if err != nil {
		zapLog.Warn("Failed to declare exchange", zap.Error(err))
	}

	// Publish event to RabbitMQ
	err = s.mq.PublishJSON(ctx, "upload_events", "upload.created", metadata)
	if err != nil {
		zapLog.Error("Failed to publish upload event", zap.Error(err))
	} else {
		zapLog.Info("Upload event published", zap.String("file", originalFileName))
	}

	return metadata, nil
}
