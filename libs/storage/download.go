package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
)

func UploadFile(client *minio.Client, bucket, filePath, objectName string) error {
    _, err := client.FPutObject(
        context.Background(),
        bucket,
        objectName,
        filePath,
        minio.PutObjectOptions{ContentType: "video/mp4"},
    )
    return err
}
