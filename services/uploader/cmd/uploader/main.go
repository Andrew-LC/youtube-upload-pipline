package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Andrew-LC/uploader/internal/api"
	"github.com/Andrew-LC/uploader/internal/app"
	"github.com/Andrew-LC/uploader/internal/repository"
	"github.com/Andrew-LC/libs/mq"
)

const (
	minioEndpoint = "MINIO_ENDPOINT"
	minioAccessKey = "MINIO_ACCESS_KEY"
	minioSecretKey = "MINIO_SECRET_KEY"
	minioBucket = "MINIO_BUCKET_NAME"
	servicePort = "SERVICE_PORT"
	mqURI = "MQ_URI"
)

func main() {
	os.Setenv(minioEndpoint, "localhost:9000")
	os.Setenv(minioAccessKey, "minioadmin")
	os.Setenv(minioSecretKey, "minioadmin")
	os.Setenv(minioBucket, "uploads")
	os.Setenv(servicePort, "8080")
	os.Setenv(mqURI, "amqp://guest:guest@localhost:5672/")
	
	endpoint := os.Getenv(minioEndpoint)
	accessKey := os.Getenv(minioAccessKey)
	secretKey := os.Getenv(minioSecretKey)
	bucketName := os.Getenv(minioBucket)
	port := os.Getenv(servicePort)
	mqUri := os.Getenv(mqURI)

	if endpoint == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		log.Fatal("MinIO configuration environment variables (ENDPOINT, ACCESS_KEY, SECRET_KEY, BUCKET_NAME) must be set.")
	}
	if port == "" {
		port = "8080" 
	}
	if mqUri == "" {
		log.Fatal("Rabbitmq configuration requires a dialup key")
	}

	rabbitMQ, err := mq.NewRabbitMQ(mqUri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	repo, err := repository.NewMinIORepo(endpoint, accessKey, secretKey, bucketName)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO Repository: %v", err)
	}
	log.Printf("Successfully connected to MinIO endpoint: %s", endpoint)

	svc := app.NewUploadService(repo, bucketName, rabbitMQ)

	handler := api.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", handler.UploadFileHandler) 

	listenAddr := ":" + port
	log.Printf("Upload Service starting on port %s...", port)
	
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
