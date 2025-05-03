package transcoder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type FFmpegTranscoder struct{}

func (f FFmpegTranscoder) Transcode(inputPath, outputPath string, res Resolution) error {
	dir := filepath.Join(outputPath, res.Label)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	segmentPattern := filepath.Join(dir, "segment_%03d.ts")
	playlistPath := filepath.Join(dir, "index.m3u8")

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=-2:%d", res.Height),
		"-c:a", "aac", "-b:a", "128k",
		"-c:v", "h264", "-b:v", fmt.Sprintf("%dk", res.Bandwidth/1000),
		"-hls_time", "4",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentPattern,
		"-f", "hls", playlistPath,
	)

	fmt.Println("Starting HLS:", res.Label)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return fmt.Errorf("HLS failed for %s: %w", res.Label, err)
	}
	fmt.Println("Finished HLS:", res.Label)
	return nil
}
