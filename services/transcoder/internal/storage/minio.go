package storage

import (
	"io"
	"context"
	"yup/Andrew-LC/libs/models"
)

type Storage interface {
	UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (models.FileMetaData, error)
	GetVideoObject(ctx context.Context, bucketName string, objectName string) (models.VideoObjectData, error)
}

