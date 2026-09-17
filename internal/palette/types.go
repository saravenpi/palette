package palette

type Theme struct {
	Background string     `json:"background"`
	Foreground string     `json:"foreground"`
	Cursor     string     `json:"cursor"`
	Palette    [16]string `json:"palette"`
}

type Options struct {
	Mode      string
	Light     bool
	MaxDim    int
	NumColors int
}
