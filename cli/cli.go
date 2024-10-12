package cli

type cli struct {
	iconType iconType
}

func NewCLI(iconType iconType) *cli {
	return &cli{iconType: iconType}
}
