package formats

import (
	"fmt"
	"strings"
	"testing"
)

func TestFormatWidgets(t *testing.T) {
	theme := coherenceTheme()
	out, err := FormatTheme(theme, "widgets")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []struct {
		name string
		line string
	}{
		{"bg", "@define-color palette_bg " + theme.Background + ";"},
		{"fg", "@define-color palette_fg " + theme.Foreground + ";"},
		{"cursor", "@define-color palette_cursor " + theme.Cursor + ";"},
	}
	for i, c := range theme.Palette {
		want = append(want, struct{ name, line string }{
			name: fmt.Sprintf("color%d", i),
			line: fmt.Sprintf("@define-color palette_color%d %s;", i, c),
		})
	}
	for _, tt := range want {
		if !strings.Contains(out, tt.line+"\n") {
			t.Errorf("%s: widgets output missing %q:\n%s", tt.name, tt.line, out)
		}
	}
	if got := strings.Count(out, "@define-color"); got != 19 {
		t.Errorf("widgets output should define 19 colors, got %d:\n%s", got, out)
	}
	if got := len(strings.Split(strings.TrimSpace(out), "\n")); got != 19 {
		t.Errorf("widgets output should be 19 lines, got %d:\n%s", got, out)
	}
}
