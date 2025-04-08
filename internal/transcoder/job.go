package transcoder

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func RunTranscodingJob(inputPath, outputPath string, transcoder Transcoder, playlistGen PlaylistGenerator) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %w", err)
	}

	ext := filepath.Ext(inputPath)
	if ext != ".mp4" && ext != ".mkv" && ext != ".avi" {
		return fmt.Errorf("unsupported file type: %s", inputPath)
	}

	resolutions := generateResolutions()

	numWorkers := len(resolutions)
	if runtime.NumCPU() < numWorkers {
		numWorkers = runtime.NumCPU()
	}

	var wg sync.WaitGroup
	taskCh := make(chan Resolution)
	errCh := make(chan error)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for res := range taskCh {
				if err := transcoder.Transcode(inputPath, outputPath, res); err != nil {
					errCh <- err
				}
			}
		}()
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

	if err := playlistGen.Generate(outputPath, resolutions); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
