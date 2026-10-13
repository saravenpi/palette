package formats

import (
	"strings"
	"testing"
)

func TestFormatMako(t *testing.T) {
	theme := coherenceTheme()
	out, err := FormatTheme(theme, "mako")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wants := []string{
		"background-color=" + theme.Background + makoBackgroundAlpha,
		"text-color=" + theme.Foreground,
		"border-color=" + theme.Palette[4],
		"progress-color=over " + theme.Palette[4],
		"[urgency=critical]",
		"border-color=" + theme.Palette[1],
		"progress-color=over " + theme.Palette[1],
	}
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("mako output missing %q:\n%s", want, out)
		}
	}
}

func TestWithAlpha(t *testing.T) {
	cases := map[string]string{
		"#19120c":     "#19120cf2",
		"#19120C":     "#19120cf2",
		"#19120cff":   "#19120cff",
		"not-a-color": "not-a-color",
		"#12345":      "#12345",
	}
	for in, want := range cases {
		if got := withAlpha(in, makoBackgroundAlpha); got != want {
			t.Errorf("withAlpha(%q) = %q, want %q", in, got, want)
		}
	}
}
