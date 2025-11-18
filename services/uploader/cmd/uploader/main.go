package main

import (
	"log"
	"net/http"
	"github.com/Andrew-LC/storage"
	"github.com/Andrew-LC/uploader/internal"
)



func main() {
    minoClient := storage.NewClient()

    mux := http.NewServeMux()
    mux.HandleFunc("POST /upload", internal.UploadHandler(minoClient))

    log.Println("Uploader running on :8080")
    http.ListenAndServe(":8080", mux)
}
