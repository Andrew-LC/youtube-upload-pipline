package repository

import (
	"io"
	"context"
	models "github.com/Andrew-LC/libs/models"
)

type UploadRepository interface {
	UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (models.FileMetaData, error)
}

