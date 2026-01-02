package ffmpeg

import (
	"io"
	"fmt"
	"context"
	"os/exec"
)

type FFmpeg struct {
	ctx context.Context
}

func New(ctx context.Context) *FFmpeg {
	return &FFmpeg{ctx: ctx}
}

func (f *FFmpeg) TranscodeStream(
	input io.Reader,
	outputPath string,
	width int,
	height int,
	br string,
) error {
	cmd := exec.CommandContext(
		f.ctx,
		"ffmpeg",
		"-i", "pipe:0",
		"-vf", fmt.Sprintf("scale=%d:%d", width, height),
		"-c:v", "libx264",
		"-b:v", br,
		"-c:a", "aac",
		outputPath,
	)

	cmd.Stdin = input
	return cmd.Run()
}
