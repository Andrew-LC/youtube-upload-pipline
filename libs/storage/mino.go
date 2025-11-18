package storage

import (
	"log"
	"io"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinoClient struct {
	client *minio.Client
	ctx context.Context 
}

func NewClient() *MinoClient {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false, // local → http
	})
	if err != nil {
		log.Fatalf("failed to create MinIO client: %v", err)
	}
	return &MinoClient{
		client: client,
		ctx: context.Background(),
	} 
}

func (c *MinoClient) UploadStream(bucket, object string, r io.Reader, size int64) error {
	_, err := c.client.PutObject(c.ctx, bucket, object, r, size, minio.PutObjectOptions{
        ContentType: "video/mp4",
    })
    return err
}

func (c *MinoClient) EnsureBucket(bucketName string) error {
    exists, err := c.client.BucketExists(c.ctx, bucketName)
    if err != nil {
        return err
    }

    if !exists {
        return c.client.MakeBucket(c.ctx, bucketName, minio.MakeBucketOptions{})
    }

    return nil
}

func (c *MinoClient) UploadFile(bucket, filePath, objectName string) error {
    _, err := c.client.FPutObject(
	c.ctx,
        bucket,
        objectName,
        filePath,
        minio.PutObjectOptions{ContentType: "video/mp4"},
    )
    return err
}
