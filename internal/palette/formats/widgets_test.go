package formats

import (
	"fmt"
	"strings"
	"testing"
)

// declared counts the lines that open with prefix, so the mechanisms named in
// the stylesheet's header comment never inflate a declaration count.
func declared(out, prefix string) int {
	total := 0
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			total++
		}
	}
	return total
}

// baseColors pairs the three colors that are not in the 16-color array with
// their custom property and named color stems.
func baseColors(theme Theme) []struct{ name, value string } {
	return []struct{ name, value string }{
		{"bg", theme.Background},
		{"fg", theme.Foreground},
		{"cursor", theme.Cursor},
	}
}

func TestFormatWidgets(t *testing.T) {
	theme := coherenceTheme()
	out, err := FormatTheme(theme, "widgets")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertCustomProperties(t, "widgets", theme, out)
	assertNamedColors(t, "widgets", theme, out)
}

// assertCustomProperties checks the :root block GTK 4.16 and later — and any
// web consumer — reads with var().
func assertCustomProperties(t *testing.T, format string, theme Theme, out string) {
	t.Helper()
	if !strings.Contains(out, ":root {") {
		t.Errorf("%s output has no :root block:\n%s", format, out)
	}
	for i, c := range theme.Palette {
		want := fmt.Sprintf("--palette-color%d: %s;", i, c)
		if !strings.Contains(out, want) {
			t.Errorf("%s output missing custom property %q:\n%s", format, want, out)
		}
	}
	for _, p := range baseColors(theme) {
		want := fmt.Sprintf("--palette-%s: %s;", p.name, p.value)
		if !strings.Contains(out, want) {
			t.Errorf("%s output missing custom property %q:\n%s", format, want, out)
		}
	}
	if got := declared(out, "--palette-"); got != 19 {
		t.Errorf("%s output should define 19 custom properties, got %d:\n%s", format, got, out)
	}
}

// assertNamedColors checks the @define-color fallback that GTK 4.10 through
// 4.15 need, and that it comes after the :root block.
func assertNamedColors(t *testing.T, format string, theme Theme, out string) {
	t.Helper()
	for i, c := range theme.Palette {
		want := fmt.Sprintf("@define-color palette_color%d %s;", i, c)
		if !strings.Contains(out, want) {
			t.Errorf("%s output missing fallback %q:\n%s", format, want, out)
		}
	}
	for _, p := range baseColors(theme) {
		want := fmt.Sprintf("@define-color palette_%s %s;", p.name, p.value)
		if !strings.Contains(out, want) {
			t.Errorf("%s output missing fallback %q:\n%s", format, want, out)
		}
	}
	if got := declared(out, "@define-color"); got != 19 {
		t.Errorf("%s output should define 19 fallback colors, got %d:\n%s", format, got, out)
	}
	root, fallback := strings.Index(out, ":root {"), strings.Index(out, "@define-color")
	if root == -1 || root > fallback {
		t.Errorf("%s output should declare :root before the fallback:\n%s", format, out)
	}
}
