package transcoder

type Transcoder interface {
	Transcode(inputPath, outputPath string, res Resolution) error
}

type PlaylistGenerator interface {
	Generate(outputPath string, resolutions []Resolution) error
}
