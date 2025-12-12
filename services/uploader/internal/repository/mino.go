package repository

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	models "github.com/Andrew-LC/uploader/pkg/models"
)

type UploadRepository interface {
	UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (models.FileMetaData, error)
}

type MinIORepo struct {
	Client *minio.Client
	Bucket string
}

func NewMinIORepo(endpoint, accessKey, secretKey, bucket string) (*MinIORepo, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Set to true for HTTPS/TLS
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    found, err := minioClient.BucketExists(ctx, bucket)
    if err != nil {
        return nil, fmt.Errorf("error checking bucket existence: %w", err)
    }
    if !found {
        err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
        if err != nil {
            return nil, fmt.Errorf("failed to create bucket: %w", err)
        }
    }
    
	return &MinIORepo{
		Client: minioClient,
		Bucket: bucket,
	}, nil
}

func (r *MinIORepo) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (models.FileMetaData, error) {
	info, err := r.Client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return models.FileMetaData{}, fmt.Errorf("minio put object failed: %w", err)
	}

	url := fmt.Sprintf("http://%s/%s/%s", r.Client.EndpointURL().Host, bucketName, info.Key)

	return models.FileMetaData{
		FileName: info.Key,
		FileSize: info.Size,
		URL:      url,
	}, nil
}
