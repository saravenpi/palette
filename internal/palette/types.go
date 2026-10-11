package palette

import "palette/internal/palette/formats"

// ParseYAML parses a YAML theme definition into a Theme struct.
func ParseYAML(data []byte) (Theme, error) {
	return formats.ParseYAML(data)
}