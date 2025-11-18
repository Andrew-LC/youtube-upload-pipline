package internal

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/Andrew-LC/storage"
)

func UploadHandler(store *storage.MinoClient) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, 100<<20) 

        err := r.ParseMultipartForm(300 << 20)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, "Invalid upload:", err)
            return
        }

        file, header, err := r.FormFile("video")
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, "Missing 'video' file:", err)
            return
        }
        defer file.Close()

        // Generate unique video ID
        videoID := uuid.New().String()
        objectName := videoID + "-" + header.Filename

        // Upload to MinIO
        err = store.UploadStream("raw-videos", objectName, file, header.Size)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintln(w, "Failed to store file:", err)
            return
        }

        fmt.Fprintf(w, "Uploaded successfully as %s\n", objectName)
    }
}
