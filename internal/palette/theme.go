package palette

import (
	"math"
	"sort"
	"strings"

	"github.com/Nadim147c/material/v3/blend"
	matcolor "github.com/Nadim147c/material/v3/color"
	"github.com/Nadim147c/material/v3/dynamic"
	"github.com/Nadim147c/material/v3/num"
	"github.com/Nadim147c/material/v3/palettes"
	"github.com/lucasb-eyer/go-colorful"
)

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func toARGB(c colorful.Color) matcolor.ARGB {
	return matcolor.NewARGB(255,
		uint8(clamp(c.R*255.0, 0, 255)),
		uint8(clamp(c.G*255.0, 0, 255)),
		uint8(clamp(c.B*255.0, 0, 255)),
	)
}

func hexLower(c matcolor.ARGB) string {
	return strings.ToLower(c.HexRGB())
}

func getVariant(mode string) dynamic.Variant {
	switch mode {
	case "material", "vibrant":
		return dynamic.VariantVibrant
	case "expressive":
		return dynamic.VariantExpressive
	case "tonal":
		return dynamic.VariantTonalSpot
	case "content":
		return dynamic.VariantContent
	case "rainbow":
		return dynamic.VariantRainbow
	case "fruit_salad":
		return dynamic.VariantFruitSalad
	default:
		return dynamic.VariantTonalSpot
	}
}

func buildAnsiPalettes(scored []matcolor.ARGB, seedARGB matcolor.ARGB) [6]*palettes.TonalPalette {
	ansiHues := [6]float64{25, 140, 85, 255, 315, 195}
	var res [6]*palettes.TonalPalette

	for i, targetHue := range ansiHues {
		bestDist := 999.0
		var bestHct matcolor.Hct
		found := false

		for _, sc := range scored {
			hct := sc.ToHct()
			if hct.Chroma < 14 {
				continue
			}
			dist := num.DifferenceDegrees(hct.Hue, targetHue)
			if dist < bestDist && dist < 35 {
				bestDist = dist
				bestHct = hct
				found = true
			}
		}

		if found {
			chroma := math.Max(bestHct.Chroma, 40)
			res[i] = palettes.FromHueAndChroma(bestHct.Hue, chroma)
		} else {
			base := matcolor.NewHct(targetHue, 55, 60).ToARGB()
			harm := blend.Harmonize(base, seedARGB)
			harmHct := harm.ToHct()
			res[i] = palettes.FromHueAndChroma(harmHct.Hue, math.Max(harmHct.Chroma, 44))
		}
	}
	return res
}

func GenerateTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	if opts.Light {
		return generateLightTheme(colors, dom, opts)
	}
	return generateDarkTheme(colors, dom, opts)
}

func generateDarkTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	seedARGB := toARGB(dom)
	seedHct := seedARGB.ToHct()

	scored := make([]matcolor.ARGB, len(colors))
	for i, c := range colors {
		scored[i] = toARGB(c)
	}

	variant := getVariant(opts.Mode)
	scheme := dynamic.NewDynamicScheme(seedHct, variant, 0.0, true, dynamic.PlatformPhone, dynamic.Version2025)

	bg := hexLower(scheme.NeutralPalette.Tone(6))
	fg := hexLower(scheme.NeutralPalette.Tone(90))
	c0 := hexLower(scheme.NeutralPalette.Tone(12))
	c7 := hexLower(scheme.NeutralPalette.Tone(80))
	c8 := hexLower(scheme.NeutralVariantPalette.Tone(35))
	c15 := hexLower(scheme.NeutralPalette.Tone(96))
	cursor := hexLower(scheme.PrimaryPalette.Tone(80))

	var normalAccents [6]string
	var brightAccents [6]string

	switch opts.Mode {
	case "duo":
		bestA := dom
		bestB := dom
		maxDist := 0.0
		if len(colors) >= 2 {
			for i := 0; i < len(colors); i++ {
				for j := i + 1; j < len(colors); j++ {
					d := colors[i].DistanceLab(colors[j])
					if d > maxDist {
						maxDist = d
						bestA = colors[i]
						bestB = colors[j]
					}
				}
			}
		} else {
			h, s, l := dom.Hsl()
			bestB = colorful.Hsl(math.Mod(h+180, 360), clamp(s, 0.4, 0.9), clamp(l, 0.4, 0.7))
		}

		hctA := toARGB(bestA).ToHct()
		hctB := toARGB(bestB).ToHct()
		palA := palettes.FromHueAndChroma(hctA.Hue, math.Max(hctA.Chroma, 40))
		palB := palettes.FromHueAndChroma(hctB.Hue, math.Max(hctB.Chroma, 40))

		colA := hexLower(palA.Tone(70))
		colB := hexLower(palB.Tone(70))
		bcolA := hexLower(palA.Tone(85))
		bcolB := hexLower(palB.Tone(85))

		for i := 0; i < 3; i++ {
			normalAccents[i] = colA
			brightAccents[i] = bcolA
		}
		for i := 3; i < 6; i++ {
			normalAccents[i] = colB
			brightAccents[i] = bcolB
		}
		cursor = colB

	case "dominant":
		sorted := make([]colorful.Color, len(colors))
		copy(sorted, colors)
		sort.Slice(sorted, func(i, j int) bool {
			return toARGB(sorted[i]).ToHct().Tone < toARGB(sorted[j]).ToHct().Tone
		})
		if len(sorted) == 0 {
			sorted = append(sorted, dom)
		}
		for i := 0; i < 6; i++ {
			c := sorted[i%len(sorted)]
			hct := toARGB(c).ToHct()
			pal := palettes.FromHueAndChroma(hct.Hue, math.Max(hct.Chroma, 40))
			normalAccents[i] = hexLower(pal.Tone(70))
			brightAccents[i] = hexLower(pal.Tone(85))
		}
		cursor = normalAccents[1]

	default:
		palettesList := buildAnsiPalettes(scored, seedARGB)
		for i := 0; i < 6; i++ {
			normalAccents[i] = hexLower(palettesList[i].Tone(70))
			brightAccents[i] = hexLower(palettesList[i].Tone(85))
		}
	}

	return Theme{
		Background: bg,
		Foreground: fg,
		Cursor:     cursor,
		Palette: [16]string{
			c0,
			normalAccents[0],
			normalAccents[1],
			normalAccents[2],
			normalAccents[3],
			normalAccents[4],
			normalAccents[5],
			c7,
			c8,
			brightAccents[0],
			brightAccents[1],
			brightAccents[2],
			brightAccents[3],
			brightAccents[4],
			brightAccents[5],
			c15,
		},
	}
}

func generateLightTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	seedARGB := toARGB(dom)
	seedHct := seedARGB.ToHct()

	scored := make([]matcolor.ARGB, len(colors))
	for i, c := range colors {
		scored[i] = toARGB(c)
	}

	variant := getVariant(opts.Mode)
	scheme := dynamic.NewDynamicScheme(seedHct, variant, 0.0, false, dynamic.PlatformPhone, dynamic.Version2025)

	bg := hexLower(scheme.NeutralPalette.Tone(96))
	fg := hexLower(scheme.NeutralPalette.Tone(12))
	c0 := hexLower(scheme.NeutralPalette.Tone(20))
	c7 := hexLower(scheme.NeutralPalette.Tone(80))
	c8 := hexLower(scheme.NeutralVariantPalette.Tone(55))
	c15 := hexLower(scheme.NeutralPalette.Tone(99))
	cursor := hexLower(scheme.PrimaryPalette.Tone(40))

	var normalAccents [6]string
	var brightAccents [6]string

	switch opts.Mode {
	case "duo":
		bestA := dom
		bestB := dom
		maxDist := 0.0
		if len(colors) >= 2 {
			for i := 0; i < len(colors); i++ {
				for j := i + 1; j < len(colors); j++ {
					d := colors[i].DistanceLab(colors[j])
					if d > maxDist {
						maxDist = d
						bestA = colors[i]
						bestB = colors[j]
					}
				}
			}
		} else {
			h, s, l := dom.Hsl()
			bestB = colorful.Hsl(math.Mod(h+180, 360), clamp(s, 0.4, 0.9), clamp(l, 0.4, 0.7))
		}

		hctA := toARGB(bestA).ToHct()
		hctB := toARGB(bestB).ToHct()
		palA := palettes.FromHueAndChroma(hctA.Hue, math.Max(hctA.Chroma, 40))
		palB := palettes.FromHueAndChroma(hctB.Hue, math.Max(hctB.Chroma, 40))

		colA := hexLower(palA.Tone(45))
		colB := hexLower(palB.Tone(45))
		bcolA := hexLower(palA.Tone(35))
		bcolB := hexLower(palB.Tone(35))

		for i := 0; i < 3; i++ {
			normalAccents[i] = colA
			brightAccents[i] = bcolA
		}
		for i := 3; i < 6; i++ {
			normalAccents[i] = colB
			brightAccents[i] = bcolB
		}
		cursor = colB

	case "dominant":
		sorted := make([]colorful.Color, len(colors))
		copy(sorted, colors)
		sort.Slice(sorted, func(i, j int) bool {
			return toARGB(sorted[i]).ToHct().Tone < toARGB(sorted[j]).ToHct().Tone
		})
		if len(sorted) == 0 {
			sorted = append(sorted, dom)
		}
		for i := 0; i < 6; i++ {
			c := sorted[i%len(sorted)]
			hct := toARGB(c).ToHct()
			pal := palettes.FromHueAndChroma(hct.Hue, math.Max(hct.Chroma, 40))
			normalAccents[i] = hexLower(pal.Tone(45))
			brightAccents[i] = hexLower(pal.Tone(35))
		}
		cursor = normalAccents[1]

	default:
		palettesList := buildAnsiPalettes(scored, seedARGB)
		for i := 0; i < 6; i++ {
			normalAccents[i] = hexLower(palettesList[i].Tone(45))
			brightAccents[i] = hexLower(palettesList[i].Tone(35))
		}
	}

	return Theme{
		Background: bg,
		Foreground: fg,
		Cursor:     cursor,
		Palette: [16]string{
			c0,
			normalAccents[0],
			normalAccents[1],
			normalAccents[2],
			normalAccents[3],
			normalAccents[4],
			normalAccents[5],
			c7,
			c8,
			brightAccents[0],
			brightAccents[1],
			brightAccents[2],
			brightAccents[3],
			brightAccents[4],
			brightAccents[5],
			c15,
		},
	}
}
