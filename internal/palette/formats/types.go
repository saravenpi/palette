package formats

// Theme represents a terminal color palette with 16 ANSI colors plus background,
// foreground, and cursor colors.
type Theme struct {
	Background string     `json:"background" yaml:"background"`
	Foreground string     `json:"foreground" yaml:"foreground"`
	Cursor     string     `json:"cursor" yaml:"cursor"`
	Palette    [16]string `json:"palette" yaml:"palette"`
}

// Options configures how a theme is generated.
type Options struct {
	Mode      string
	Light     bool
	MaxDim    int
	NumColors int
}