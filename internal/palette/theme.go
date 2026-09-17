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

func GenerateTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	if opts.Light {
		return generateLightTheme(colors, dom, opts)
	}
	return generateDarkTheme(colors, dom, opts)
}

func generateDarkTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	seedARGB := toARGB(dom)
	seedHct := seedARGB.ToHct()

	hcts := make([]matcolor.Hct, len(colors))
	for i, c := range colors {
		hcts[i] = toARGB(c).ToHct()
	}

	bgChroma := clamp(seedHct.Chroma*0.35, 6.0, 20.0)
	fgChroma := clamp(seedHct.Chroma*0.12, 2.0, 8.0)
	bg := hexLower(palettes.FromHueAndChroma(seedHct.Hue, bgChroma).Tone(6))
	fg := hexLower(palettes.FromHueAndChroma(seedHct.Hue, fgChroma).Tone(90))
	c0 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, bgChroma*1.2).Tone(12))
	c7 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, fgChroma).Tone(80))
	c8 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, bgChroma*1.2).Tone(46))
	c15 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, fgChroma*0.5).Tone(97))
	cursor := hexLower(palettes.FromHueAndChroma(seedHct.Hue, math.Max(seedHct.Chroma, 36.0)).Tone(82))

	var normalAccents [6]string
	var brightAccents [6]string

	switch opts.Mode {
	case "wal", "image":
		ansiHues := [6]float64{25, 140, 85, 250, 315, 195}
		used := make(map[int]bool)
		for i, target := range ansiHues {
			bestIdx := -1
			bestDist := 999.0
			for j, hct := range hcts {
				if used[j] {
					continue
				}
				dist := num.DifferenceDegrees(hct.Hue, target)
				if dist < bestDist && dist < 65 {
					bestDist = dist
					bestIdx = j
				}
			}
			var pal *palettes.TonalPalette
			if bestIdx >= 0 {
				used[bestIdx] = true
				matched := hcts[bestIdx]
				pal = palettes.FromHueAndChroma(matched.Hue, math.Max(matched.Chroma, 32.0))
			} else {
				base := matcolor.NewHct(target, 50, 60).ToARGB()
				harm := blend.Harmonize(base, seedARGB).ToHct()
				pal = palettes.FromHueAndChroma(harm.Hue, clamp(seedHct.Chroma*0.9, 30.0, 50.0))
			}
			normalAccents[i] = hexLower(pal.Tone(70))
			brightAccents[i] = hexLower(pal.Tone(86))
		}
		cursor = normalAccents[0]

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
		palA := palettes.FromHueAndChroma(hctA.Hue, math.Max(hctA.Chroma, 36.0))
		palB := palettes.FromHueAndChroma(hctB.Hue, math.Max(hctB.Chroma, 36.0))

		colA := hexLower(palA.Tone(70))
		colB := hexLower(palB.Tone(70))
		bcolA := hexLower(palA.Tone(86))
		bcolB := hexLower(palB.Tone(86))

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
			pal := palettes.FromHueAndChroma(hct.Hue, math.Max(hct.Chroma, 34.0))
			normalAccents[i] = hexLower(pal.Tone(70))
			brightAccents[i] = hexLower(pal.Tone(86))
		}
		cursor = normalAccents[1]

	case "material":
		scheme := dynamic.NewDynamicScheme(seedHct, dynamic.VariantVibrant, 0.0, true, dynamic.PlatformPhone, dynamic.Version2025)
		palGreen := palettes.FromHueAndChroma(blend.Harmonize(matcolor.NewHct(140, 50, 60).ToARGB(), seedARGB).ToHct().Hue, math.Max(scheme.PrimaryPalette.KeyColor.Chroma, 34.0))
		palYellow := palettes.FromHueAndChroma(blend.Harmonize(matcolor.NewHct(85, 50, 60).ToARGB(), seedARGB).ToHct().Hue, math.Max(scheme.SecondaryPalette.KeyColor.Chroma, 34.0))
		palCyan := palettes.FromHueAndChroma(blend.Harmonize(matcolor.NewHct(195, 50, 60).ToARGB(), seedARGB).ToHct().Hue, math.Max(scheme.TertiaryPalette.KeyColor.Chroma, 34.0))

		normalAccents[0] = hexLower(scheme.ErrorPalette.Tone(70))
		normalAccents[1] = hexLower(palGreen.Tone(70))
		normalAccents[2] = hexLower(palYellow.Tone(74))
		normalAccents[3] = hexLower(scheme.PrimaryPalette.Tone(72))
		normalAccents[4] = hexLower(scheme.TertiaryPalette.Tone(72))
		normalAccents[5] = hexLower(palCyan.Tone(74))

		brightAccents[0] = hexLower(scheme.ErrorPalette.Tone(86))
		brightAccents[1] = hexLower(palGreen.Tone(86))
		brightAccents[2] = hexLower(palYellow.Tone(88))
		brightAccents[3] = hexLower(scheme.PrimaryPalette.Tone(86))
		brightAccents[4] = hexLower(scheme.TertiaryPalette.Tone(86))
		brightAccents[5] = hexLower(palCyan.Tone(88))
		cursor = hexLower(scheme.PrimaryPalette.Tone(82))

	case "expressive":
		ansiHues := [6]float64{25, 140, 85, 250, 315, 195}
		usedHcts := make(map[int]bool)
		for i, targetHue := range ansiHues {
			bestIdx := -1
			bestDist := 999.0
			for j, hct := range hcts {
				if usedHcts[j] {
					continue
				}
				dist := num.DifferenceDegrees(hct.Hue, targetHue)
				if dist < bestDist && dist < 50 {
					bestDist = dist
					bestIdx = j
				}
			}

			var pal *palettes.TonalPalette
			if bestIdx >= 0 {
				usedHcts[bestIdx] = true
				matched := hcts[bestIdx]
				pal = palettes.FromHueAndChroma(matched.Hue, math.Max(matched.Chroma*1.3, 38.0))
			} else {
				base := matcolor.NewHct(targetHue, 55, 60).ToARGB()
				harm := blend.Harmonize(base, seedARGB).ToHct()
				pal = palettes.FromHueAndChroma(harm.Hue, clamp(seedHct.Chroma*1.2, 36.0, 58.0))
			}

			normalAccents[i] = hexLower(pal.Tone(72))
			brightAccents[i] = hexLower(pal.Tone(88))
		}
		cursor = hexLower(palettes.FromHueAndChroma(seedHct.Hue+120.0, math.Max(seedHct.Chroma, 38.0)).Tone(84))

	default:
		ansiHues := [6]float64{25, 140, 85, 250, 315, 195}
		usedHcts := make(map[int]bool)

		chromaMultiplier := 1.0
		toneNormal := 70.0
		toneBright := 86.0
		if opts.Mode == "vibrant" {
			chromaMultiplier = 1.30
			toneNormal = 72.0
			toneBright = 88.0
		}

		for i, targetHue := range ansiHues {
			bestIdx := -1
			bestDist := 999.0
			for j, hct := range hcts {
				if usedHcts[j] {
					continue
				}
				dist := num.DifferenceDegrees(hct.Hue, targetHue)
				if dist < bestDist && dist < 42 {
					bestDist = dist
					bestIdx = j
				}
			}

			var pal *palettes.TonalPalette
			if bestIdx >= 0 {
				usedHcts[bestIdx] = true
				matched := hcts[bestIdx]
				chroma := math.Max(matched.Chroma*chromaMultiplier, 32.0)
				pal = palettes.FromHueAndChroma(matched.Hue, chroma)
			} else {
				base := matcolor.NewHct(targetHue, 50, 60).ToARGB()
				harm := blend.Harmonize(base, seedARGB)
				harmHct := harm.ToHct()
				targetChroma := clamp(seedHct.Chroma*0.9*chromaMultiplier, 30.0, 52.0)
				pal = palettes.FromHueAndChroma(harmHct.Hue, targetChroma)
			}

			normalAccents[i] = hexLower(pal.Tone(toneNormal))
			brightAccents[i] = hexLower(pal.Tone(toneBright))
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

	hcts := make([]matcolor.Hct, len(colors))
	for i, c := range colors {
		hcts[i] = toARGB(c).ToHct()
	}

	bgChroma := clamp(seedHct.Chroma*0.25, 4.0, 14.0)
	fgChroma := clamp(seedHct.Chroma*0.12, 2.0, 8.0)
	bg := hexLower(palettes.FromHueAndChroma(seedHct.Hue, bgChroma).Tone(96))
	fg := hexLower(palettes.FromHueAndChroma(seedHct.Hue, fgChroma).Tone(12))
	c0 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, bgChroma*1.2).Tone(22))
	c7 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, fgChroma).Tone(78))
	c8 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, bgChroma*1.2).Tone(50))
	c15 := hexLower(palettes.FromHueAndChroma(seedHct.Hue, fgChroma*0.5).Tone(8))
	cursor := hexLower(palettes.FromHueAndChroma(seedHct.Hue, math.Max(seedHct.Chroma, 36.0)).Tone(40))

	var normalAccents [6]string
	var brightAccents [6]string

	switch opts.Mode {
	case "wal", "image":
		ansiHues := [6]float64{25, 140, 85, 250, 315, 195}
		used := make(map[int]bool)
		for i, target := range ansiHues {
			bestIdx := -1
			bestDist := 999.0
			for j, hct := range hcts {
				if used[j] {
					continue
				}
				dist := num.DifferenceDegrees(hct.Hue, target)
				if dist < bestDist && dist < 65 {
					bestDist = dist
					bestIdx = j
				}
			}
			var pal *palettes.TonalPalette
			if bestIdx >= 0 {
				used[bestIdx] = true
				matched := hcts[bestIdx]
				pal = palettes.FromHueAndChroma(matched.Hue, math.Max(matched.Chroma, 34.0))
			} else {
				base := matcolor.NewHct(target, 50, 60).ToARGB()
				harm := blend.Harmonize(base, seedARGB).ToHct()
				pal = palettes.FromHueAndChroma(harm.Hue, clamp(seedHct.Chroma*0.9, 30.0, 50.0))
			}
			normalAccents[i] = hexLower(pal.Tone(42))
			brightAccents[i] = hexLower(pal.Tone(32))
		}
		cursor = normalAccents[0]

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
		palA := palettes.FromHueAndChroma(hctA.Hue, math.Max(hctA.Chroma, 36.0))
		palB := palettes.FromHueAndChroma(hctB.Hue, math.Max(hctB.Chroma, 36.0))

		colA := hexLower(palA.Tone(42))
		colB := hexLower(palB.Tone(42))
		bcolA := hexLower(palA.Tone(32))
		bcolB := hexLower(palB.Tone(32))

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
			pal := palettes.FromHueAndChroma(hct.Hue, math.Max(hct.Chroma, 34.0))
			normalAccents[i] = hexLower(pal.Tone(42))
			brightAccents[i] = hexLower(pal.Tone(32))
		}
		cursor = normalAccents[1]

	case "material":
		scheme := dynamic.NewDynamicScheme(seedHct, dynamic.VariantVibrant, 0.0, false, dynamic.PlatformPhone, dynamic.Version2025)
		palGreen := palettes.FromHueAndChroma(blend.Harmonize(matcolor.NewHct(140, 50, 60).ToARGB(), seedARGB).ToHct().Hue, math.Max(scheme.PrimaryPalette.KeyColor.Chroma, 34.0))
		palYellow := palettes.FromHueAndChroma(blend.Harmonize(matcolor.NewHct(85, 50, 60).ToARGB(), seedARGB).ToHct().Hue, math.Max(scheme.SecondaryPalette.KeyColor.Chroma, 34.0))
		palCyan := palettes.FromHueAndChroma(blend.Harmonize(matcolor.NewHct(195, 50, 60).ToARGB(), seedARGB).ToHct().Hue, math.Max(scheme.TertiaryPalette.KeyColor.Chroma, 34.0))

		normalAccents[0] = hexLower(scheme.ErrorPalette.Tone(42))
		normalAccents[1] = hexLower(palGreen.Tone(42))
		normalAccents[2] = hexLower(palYellow.Tone(44))
		normalAccents[3] = hexLower(scheme.PrimaryPalette.Tone(42))
		normalAccents[4] = hexLower(scheme.TertiaryPalette.Tone(44))
		normalAccents[5] = hexLower(palCyan.Tone(44))

		brightAccents[0] = hexLower(scheme.ErrorPalette.Tone(32))
		brightAccents[1] = hexLower(palGreen.Tone(32))
		brightAccents[2] = hexLower(palYellow.Tone(34))
		brightAccents[3] = hexLower(scheme.PrimaryPalette.Tone(32))
		brightAccents[4] = hexLower(scheme.TertiaryPalette.Tone(32))
		brightAccents[5] = hexLower(palCyan.Tone(34))
		cursor = hexLower(scheme.PrimaryPalette.Tone(40))

	case "expressive":
		ansiHues := [6]float64{25, 140, 85, 250, 315, 195}
		usedHcts := make(map[int]bool)
		for i, targetHue := range ansiHues {
			bestIdx := -1
			bestDist := 999.0
			for j, hct := range hcts {
				if usedHcts[j] {
					continue
				}
				dist := num.DifferenceDegrees(hct.Hue, targetHue)
				if dist < bestDist && dist < 50 {
					bestDist = dist
					bestIdx = j
				}
			}

			var pal *palettes.TonalPalette
			if bestIdx >= 0 {
				usedHcts[bestIdx] = true
				matched := hcts[bestIdx]
				pal = palettes.FromHueAndChroma(matched.Hue, math.Max(matched.Chroma*1.3, 38.0))
			} else {
				base := matcolor.NewHct(targetHue, 55, 60).ToARGB()
				harm := blend.Harmonize(base, seedARGB).ToHct()
				pal = palettes.FromHueAndChroma(harm.Hue, clamp(seedHct.Chroma*1.2, 36.0, 58.0))
			}

			normalAccents[i] = hexLower(pal.Tone(40))
			brightAccents[i] = hexLower(pal.Tone(30))
		}
		cursor = hexLower(palettes.FromHueAndChroma(seedHct.Hue+120.0, math.Max(seedHct.Chroma, 38.0)).Tone(38))

	default:
		ansiHues := [6]float64{25, 140, 85, 250, 315, 195}
		usedHcts := make(map[int]bool)

		chromaMultiplier := 1.0
		toneNormal := 42.0
		toneBright := 32.0
		if opts.Mode == "vibrant" {
			chromaMultiplier = 1.30
			toneNormal = 40.0
			toneBright = 30.0
		}

		for i, targetHue := range ansiHues {
			bestIdx := -1
			bestDist := 999.0
			for j, hct := range hcts {
				if usedHcts[j] {
					continue
				}
				dist := num.DifferenceDegrees(hct.Hue, targetHue)
				if dist < bestDist && dist < 42 {
					bestDist = dist
					bestIdx = j
				}
			}

			var pal *palettes.TonalPalette
			if bestIdx >= 0 {
				usedHcts[bestIdx] = true
				matched := hcts[bestIdx]
				chroma := math.Max(matched.Chroma*chromaMultiplier, 32.0)
				pal = palettes.FromHueAndChroma(matched.Hue, chroma)
			} else {
				base := matcolor.NewHct(targetHue, 50, 60).ToARGB()
				harm := blend.Harmonize(base, seedARGB)
				harmHct := harm.ToHct()
				targetChroma := clamp(seedHct.Chroma*0.9*chromaMultiplier, 30.0, 52.0)
				pal = palettes.FromHueAndChroma(harmHct.Hue, targetChroma)
			}

			normalAccents[i] = hexLower(pal.Tone(toneNormal))
			brightAccents[i] = hexLower(pal.Tone(toneBright))
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
