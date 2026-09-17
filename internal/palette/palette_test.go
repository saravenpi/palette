package palette

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestFormatTheme(t *testing.T) {
	theme := Theme{
		Background: "#06050c",
		Foreground: "#dedde4",
		Cursor:     "#9281cf",
		Palette: [16]string{
			"#1b1a1f", "#24a688", "#24a688", "#24a688",
			"#9770fa", "#9770fa", "#9770fa", "#e4e4e7",
			"#403f44", "#1fefc4", "#1fefc4", "#1fefc4",
			"#d0bfff", "#d0bfff", "#d0bfff", "#f8f8fa",
		},
	}

	ghostty, err := FormatTheme(theme, "ghostty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(ghostty, "palette = 0=#1b1a1f") {
		t.Errorf("ghostty output missing color 0: %s", ghostty)
	}
	if !strings.Contains(ghostty, "background = #06050c") {
		t.Errorf("ghostty output missing background: %s", ghostty)
	}

	kitty, err := FormatTheme(theme, "kitty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(kitty, "color0 #1b1a1f") {
		t.Errorf("kitty output missing color0: %s", kitty)
	}

	alacritty, err := FormatTheme(theme, "alacritty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(alacritty, "background = \"#06050c\"") {
		t.Errorf("alacritty output missing background: %s", alacritty)
	}

	jsonOut, err := FormatTheme(theme, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(jsonOut, "\"background\": \"#06050c\"") {
		t.Errorf("json output missing background: %s", jsonOut)
	}

	tmuxOut, err := FormatTheme(theme, "tmux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(tmuxOut, "set -g status-style") {
		t.Errorf("tmux output missing status-style: %s", tmuxOut)
	}

	nvimOut, err := FormatTheme(theme, "nvim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(nvimOut, "M.colors") {
		t.Errorf("nvim output missing M.colors: %s", nvimOut)
	}

	yamlOut, err := FormatTheme(theme, "yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, err := ParseYAML([]byte(yamlOut))
	if err != nil {
		t.Fatalf("failed to parse yaml: %v", err)
	}
	if parsed.Background != theme.Background {
		t.Errorf("expected bg %s, got %s", theme.Background, parsed.Background)
	}

	hexOut, err := FormatTheme(theme, "hex")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(hexOut), "\n")
	if len(lines) != 16 {
		t.Errorf("expected 16 hex lines, got %d", len(lines))
	}
}

func TestDownscale(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1000, 500))
	for y := 0; y < 500; y++ {
		for x := 0; x < 1000; x++ {
			src.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}

	scaled := Downscale(src, 200)
	bounds := scaled.Bounds()
	if bounds.Dx() != 200 {
		t.Errorf("expected width 200, got %d", bounds.Dx())
	}
	if bounds.Dy() != 100 {
		t.Errorf("expected height 100, got %d", bounds.Dy())
	}
}

func TestGenerateTheme(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if x < 50 {
				src.Set(x, y, color.RGBA{R: 36, G: 166, B: 136, A: 255})
			} else {
				src.Set(x, y, color.RGBA{R: 151, G: 112, B: 250, A: 255})
			}
		}
	}

	colors, dom, err := ExtractColors(src, 16)
	if err != nil {
		t.Fatalf("failed to extract colors: %v", err)
	}

	themeDark := GenerateTheme(colors, dom, Options{Mode: "duo", Light: false})
	if themeDark.Background == "" || themeDark.Foreground == "" {
		t.Errorf("invalid dark theme: %+v", themeDark)
	}

	themeLight := GenerateTheme(colors, dom, Options{Mode: "duo", Light: true})
	if themeLight.Background == "" || themeLight.Foreground == "" {
		t.Errorf("invalid light theme: %+v", themeLight)
	}

	themeAnsi := GenerateTheme(colors, dom, Options{Mode: "ansi", Light: false})
	if themeAnsi.Background == "" || themeAnsi.Foreground == "" {
		t.Errorf("invalid ansi theme: %+v", themeAnsi)
	}

	modes := []string{"material", "vibrant", "expressive", "tonal", "dominant"}
	for _, m := range modes {
		for _, light := range []bool{false, true} {
			th := GenerateTheme(colors, dom, Options{Mode: m, Light: light})
			if th.Background == "" || th.Foreground == "" || th.Cursor == "" {
				t.Fatalf("mode %s (light=%v) returned empty theme: %+v", m, light, th)
			}
			for i, c := range th.Palette {
				if len(c) != 7 || !strings.HasPrefix(c, "#") {
					t.Fatalf("mode %s (light=%v) invalid palette color %d: %s", m, light, i, c)
				}
			}
		}
	}
}
