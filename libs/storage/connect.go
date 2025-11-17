package storage

import (
	"log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient() *minio.Client {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false, // local → http
	})
	if err != nil {
		log.Fatalf("failed to create MinIO client: %v", err)
	}
	return client
}
