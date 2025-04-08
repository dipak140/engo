package main

import (
	"fmt"
	"os"

	"github.com/dipak140/engo/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: engo <input-video-path>")
		return
	}

	input := os.Args[1]
	err := internal.RunTranscodingJob(input, "output")
	if err != nil {
		fmt.Printf("Transcoding wailed: %v\n", err)
	}
}
