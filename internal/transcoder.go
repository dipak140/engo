package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

type Resolution struct {
	Label     string
	Height    int
	Bandwidth int // in bits per second
}

func generateResolutions() []Resolution {
	resolutions := []Resolution{
		{"1080p", 1080, 2800000},
		{"720p", 720, 1400000},
		{"480p", 480, 800000},
		{"360p", 360, 500000},
	}
	return resolutions
}

func GenerateMasterPlaylist(outputPath string, resolutions []Resolution) error {
	masterPath := filepath.Join(outputPath, "master.m3u8")
	file, err := os.Create(masterPath)
	if err != nil {
		return fmt.Errorf("error creating master playlist: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString("#EXTM3U\n")
	if err != nil {
		return err
	}
	for _, res := range resolutions {
		width := getWidthFromHeight(res.Height) // estimate common width
		resolution := fmt.Sprintf("%dx%d", width, res.Height)

		entry := fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n%s/index.m3u8\n",
			res.Bandwidth,
			resolution,
			res.Label,
		)

		_, err := file.WriteString(entry)
		if err != nil {
			return err
		}
	}

	fmt.Println("Generated master.m3u8")
	return nil
}

func getWidthFromHeight(height int) int {
	switch height {
	case 360:
		return 640
	case 480:
		return 854
	case 720:
		return 1280
	case 1080:
		return 1920
	default:
		return 1280
	}
}

func RunTranscodingJob(inputPath, outputPath string) error {
	// Check if the input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %w", err)
	}

	// Check if the input file is a video file
	ext := filepath.Ext(inputPath)
	if ext != ".mp4" && ext != ".mkv" && ext != ".avi" {
		return fmt.Errorf("input file is not a video file: %s", inputPath)
	}

	resolutions := generateResolutions()

	// Get the number of CPU cores
	numWorkers := len(resolutions)
	// Limit the number of workers to the number of CPU cores
	if runtime.NumCPU() < numWorkers {
		numWorkers = runtime.NumCPU()
	}

	// create waitgroups
	var wg sync.WaitGroup

	taskCh := make(chan Resolution)
	errCh := make(chan error)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go hlsWorker(inputPath, outputPath, taskCh, errCh, &wg)
	}

	for _, res := range resolutions {
		taskCh <- res
	}

	close(taskCh)

	wg.Wait()
	close(errCh)
	var firstErr error
	for err := range errCh {
		if firstErr == nil {
			firstErr = err
		}
	}
	GenerateMasterPlaylist(outputPath, generateResolutions())
	return firstErr
}

func hlsWorker(inputPath, outputPath string, resolutionChannel <-chan Resolution, errorChannel chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for res := range resolutionChannel {
		outputPath := filepath.Join(outputPath, res.Label)
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			errorChannel <- fmt.Errorf("error creating output directory: %w", err)
			continue
		}
		segmentPattern := filepath.Join(outputPath, "segment_%03d.ts")
		playlistPath := filepath.Join(outputPath, "index.m3u8")

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
		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(outputBytes))
			errorChannel <- fmt.Errorf("HLS failed for %s: %w", res.Label, err)
			continue
		}
		fmt.Println("Finished HLS:", res.Label)
	}
}
