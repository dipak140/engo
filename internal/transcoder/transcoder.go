package transcoder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

type Resolution struct {
	Label  string
	Height int
}

func generateResolutions() []Resolution {
	resolutions := []Resolution{
		{"1080p", 1080},
		{"720p", 720},
		{"480p", 480},
		{"360p", 360},
	}
	return resolutions
}

func RunTranscodingJob(inputPath string) error {
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
		go newFunction(inputPath, taskCh, errCh, &wg)
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
	return firstErr
}

func newFunction(inputPath string, resolutionChannel <-chan Resolution, errorChannel chan<- error, wg *sync.WaitGroup) (bool, error) {
	defer wg.Done()
	for res := range resolutionChannel {
		output := fmt.Sprintf("%s_%s.mp4", inputPath[:len(inputPath)-len(filepath.Ext(inputPath))], res.Label)
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", fmt.Sprintf("scale=-2:%d", res.Height), "-c:a", "copy", output)

		fmt.Println("Starting:", res.Label)
		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(outputBytes))
			errorChannel <- fmt.Errorf("error transcoding %s: %w", res.Label, err)
			return true, fmt.Errorf("error transcoding %s: %w", res.Label, err)
		}
		fmt.Println("Finished:", res.Label)
	}

	return false, nil
}
