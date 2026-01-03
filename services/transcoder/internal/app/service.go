package app

import (
	"context"
	"yup/Andrew-LC/transcoder/internal/ffmpeg"
	"yup/Andrew-LC/transcoder/internal/storage"
)

type TranscodingService {
	ctx context.Context
	ff  *FFmpeg
	storage *Storage
}

