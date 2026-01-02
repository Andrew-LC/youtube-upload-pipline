package app

import (
	"context"
	"github.com/Andrew-LC/transcoder/internal/ffmpeg"
	"github.com/Andrew-LC/transcoder/internal/storage"
)

type TranscodingService {
	ctx context.Context
	ff  *FFmpeg
	storage *Storage
}

