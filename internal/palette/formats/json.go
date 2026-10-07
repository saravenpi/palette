package formats

import (
	"encoding/json"
)

func formatJSON(t Theme) (string, error) {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes) + "\n", nil
}