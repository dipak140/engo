package transcoder

import (
	"fmt"
	"os"
	"path/filepath"
)

type HLSMasterPlaylist struct{}

func (h HLSMasterPlaylist) Generate(outputPath string, resolutions []Resolution) error {
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
		width := getWidthFromHeight(res.Height)
		line := fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n%s/index.m3u8\n",
			res.Bandwidth, width, res.Height, res.Label,
		)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	fmt.Println("Generated master.m3u8")
	return nil
}
