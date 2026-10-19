package formats

import (
	"fmt"
	"strings"
	"testing"
)

func TestFormatYazi(t *testing.T) {
	theme := coherenceTheme()
	out, err := FormatTheme(theme, "yazi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertYaziLines(t, out, theme)
	assertYaziShape(t, out)
}

func assertYaziLines(t *testing.T, out string, theme Theme) {
	t.Helper()
	wants := []string{
		fmt.Sprintf("overall = { bg = %q }", theme.Background),
		fmt.Sprintf("cwd = { fg = %q, bold = true }", theme.Cursor),
		fmt.Sprintf("border_style = { fg = %q }", theme.Palette[8]),
		fmt.Sprintf("active = { fg = %q, bg = %q, bold = true }", theme.Background, theme.Palette[4]),
		fmt.Sprintf("select_main = { fg = %q, bg = %q, bold = true }", theme.Background, theme.Palette[1]),
		fmt.Sprintf("progress_normal = { fg = %q, bg = %q }", theme.Palette[2], theme.Palette[0]),
		fmt.Sprintf("cand = { fg = %q }", theme.Palette[14]),
		fmt.Sprintf("desc = { fg = %q }", theme.Palette[13]),
		fmt.Sprintf("chord = { fg = %q }", theme.Palette[6]),
		`rules = [`,
		fmt.Sprintf(`{ url = "*/", fg = %q },`, theme.Palette[4]),
	}
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("yazi output missing %q:\n%s", want, out)
		}
	}
}

func assertYaziShape(t *testing.T, out string) {
	t.Helper()
	if strings.Contains(out, "[manager]") {
		t.Errorf("yazi output uses the old [manager] table name:\n%s", out)
	}
	if strings.Count(out, "rules = [") != 1 {
		t.Errorf("yazi output must declare the filetype rules once:\n%s", out)
	}
	if strings.Count(out, "[") != strings.Count(out, "]") {
		t.Errorf("yazi output has unbalanced tables or arrays:\n%s", out)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		if strings.Count(line, "{") < strings.Count(line, "}") {
			t.Errorf("yazi output has an unterminated inline table: %q", line)
		}
	}
}
