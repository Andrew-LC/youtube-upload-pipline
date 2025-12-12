package api

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/Andrew-LC/uploader/internal/app"
)

type Handler struct {
	Service app.UploadService
}

func NewHandler(svc app.UploadService) *Handler {
	return &Handler{Service: svc}
}

func (h *Handler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(4 << 30) 

	file, header, err := r.FormFile("file") 
	if err != nil {
		http.Error(w, "Error retrieving file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	metadata, err := h.Service.ProcessUpload(r.Context(), header.Filename, file, header.Size, header.Header.Get("Content-Type"))

	if err != nil {
		log.Printf("Upload failed: %v", err)
        if err.Error() == "file size exceeds 4GB limit" {
            http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
            return
        }
		http.Error(w, "Internal error processing upload", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metadata)
}
