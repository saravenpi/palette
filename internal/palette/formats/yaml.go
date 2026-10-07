package formats

import (
	"gopkg.in/yaml.v3"
)

func formatYAML(t Theme) (string, error) {
	bytes, err := yaml.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseYAML(data []byte) (Theme, error) {
	var t Theme
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return Theme{}, err
	}
	return t, nil
}