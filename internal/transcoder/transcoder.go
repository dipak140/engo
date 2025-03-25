package transcoder

import(
	"fmt"
	"os/exec"
	"path/filepath"
)

var resolutions = []struct {
	Label string
	Height int
}{
	{"1080p", 1080},
	{"720p", 720},
	{"480p", 480},
	{"360p", 360},
}

func RunTranscodingJob(inputPath string) error{
	for _, res := range resolutions{ 
		output := fmt.Sprintf("%s_%s.mp4", inputPath[:len(inputPath)-len(filepath.Ext(inputPath))], res.Label)
		cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", fmt.Sprintf("scale=-2:%d", res.Height), "-c:a", "copy", output)

		fmt.Println("Starting:", res.Label)
		outputBytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(outputBytes))
			return fmt.Errorf("error transcoding %s: %w", res.Label, err)
		}
		fmt.Println("Finished:", res.Label)
	}
	return nil
}

