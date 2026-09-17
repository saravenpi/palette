package palette

type Theme struct {
	Background string     `json:"background" yaml:"background"`
	Foreground string     `json:"foreground" yaml:"foreground"`
	Cursor     string     `json:"cursor" yaml:"cursor"`
	Palette    [16]string `json:"palette" yaml:"palette"`
}

type Options struct {
	Mode      string
	Light     bool
	MaxDim    int
	NumColors int
}
