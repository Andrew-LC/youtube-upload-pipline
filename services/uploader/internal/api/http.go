package api

import (
	"encoding/json"
	"net/http"

	"github.com/Andrew-LC/libs/logger"
	"github.com/Andrew-LC/uploader/internal/app"
	"go.uber.org/zap"
)

type Handler struct {
	Service app.UploadService
	logger  *logger.Logger
}

func NewHandler(svc app.UploadService, l *logger.Logger) *Handler {
	return &Handler{Service: svc, logger: l}
}

func (h *Handler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	zapLog := h.logger.GetZapLogger()
	r.ParseMultipartForm(4 << 30)

	file, header, err := r.FormFile("file")
	if err != nil {
		zapLog.Error("Error retrieving file from form", zap.Error(err))
		http.Error(w, "Error retrieving file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	metadata, err := h.Service.ProcessUpload(r.Context(), header.Filename, file, header.Size, header.Header.Get("Content-Type"))

	if err != nil {
		zapLog.Error("Upload failed", zap.Error(err))
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
