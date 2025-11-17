package storage

import (
	"io"
	"context"
	"github.com/minio/minio-go/v7"
)

// func UploadFile(client *minio.Client, bucket, filePath, objectName string) error {
//     _, err := client.FPutObject(
//         context.Background(),
//         bucket,
//         objectName,
//         filePath,
//         minio.PutObjectOptions{ContentType: "video/mp4"},
//     )
//     return err
// }

func UploadStream(client *minio.Client, bucket, object string, r io.Reader, size int64) error {
	_, err := client.PutObject(context.Background(), bucket, object, r, size, minio.PutObjectOptions{
        ContentType: "video/mp4",
    })
    return err
}
