package transcoder

type Resolution struct {
	Label     string
	Height    int
	Bandwidth int
}

func generateResolutions() []Resolution {
	return []Resolution{
		{"1080p", 1080, 2800000},
		{"720p", 720, 1400000},
		{"480p", 480, 800000},
		{"360p", 360, 500000},
	}
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
