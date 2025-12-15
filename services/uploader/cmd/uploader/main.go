package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Andrew-LC/libs/logger"
	"github.com/Andrew-LC/libs/mq"
	"github.com/Andrew-LC/libs/storage"
	"github.com/Andrew-LC/uploader/internal/api"
	"github.com/Andrew-LC/uploader/internal/app"
	"go.uber.org/zap"
)

const (
	minioEndpoint  = "MINIO_ENDPOINT"
	minioAccessKey = "MINIO_ACCESS_KEY"
	minioSecretKey = "MINIO_SECRET_KEY"
	minioBucket    = "MINIO_BUCKET_NAME"
	servicePort    = "SERVICE_PORT"
	mqURI          = "MQ_URI"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	l, err := logger.NewLogger("upload_log", true)
	// Initialize Logger
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer l.Sync()
	zapLog := l.GetZapLogger()

	endpoint := getEnv(minioEndpoint, "localhost:9000")
	accessKey := getEnv(minioAccessKey, "minioadmin")
	secretKey := getEnv(minioSecretKey, "minioadmin")
	bucketName := getEnv(minioBucket, "uploads")
	port := getEnv(servicePort, "8080")
	mqUri := getEnv(mqURI, "amqp://guest:guest@localhost:5672/")

	if endpoint == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		zapLog.Fatal("MinIO configuration environment variables (ENDPOINT, ACCESS_KEY, SECRET_KEY, BUCKET_NAME) must be set.")
	}
	if mqUri == "" {
		zapLog.Fatal("Rabbitmq configuration requires a dialup key")
	}

	rabbitMQ, err := mq.NewRabbitMQ(mqUri)
	if err != nil {
		zapLog.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer rabbitMQ.Close()

	// Initialize MinIO Repo with secure=false (TODO: make configurable)
	repo, err := storage.NewMinIORepo(endpoint, accessKey, secretKey, bucketName, false)
	if err != nil {
		zapLog.Fatal("Failed to initialize MinIO Repository", zap.Error(err))
	}
	zapLog.Info("Successfully connected to MinIO endpoint", zap.String("endpoint", endpoint))

	svc := app.NewUploadService(repo, bucketName, rabbitMQ, l)

	handler := api.NewHandler(svc, l)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", handler.UploadFileHandler)

	listenAddr := ":" + port
	zapLog.Info("Upload Service starting...", zap.String("port", port))

	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		zapLog.Fatal("Server failed", zap.Error(err))
	}
}
