package palette

import (
	"math"
	"sort"

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

func hueDist(h1, h2 float64) float64 {
	d := math.Abs(h1 - h2)
	if d > 180 {
		return 360 - d
	}
	return d
}

func filterAccents(colors []colorful.Color) []colorful.Color {
	var accents []colorful.Color
	for _, c := range colors {
		_, s, l := c.Hsl()
		if s >= 0.12 && l >= 0.15 && l <= 0.85 {
			accents = append(accents, c)
		}
	}
	if len(accents) == 0 {
		return colors
	}
	return accents
}

func GenerateTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	if opts.Light {
		return generateLightTheme(colors, dom, opts)
	}
	return generateDarkTheme(colors, dom, opts)
}

func generateDarkTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	dh, ds, _ := dom.Hsl()

	bg := colorful.Hsl(dh, clamp(ds*0.25, 0.04, 0.22), 0.045)
	fg := colorful.Hsl(dh, clamp(ds*0.12, 0.02, 0.14), 0.88)
	c0 := colorful.Hsl(dh, clamp(ds*0.20, 0.04, 0.18), 0.11)
	c8 := colorful.Hsl(dh, clamp(ds*0.18, 0.04, 0.16), 0.28)
	c7 := colorful.Hsl(dh, clamp(ds*0.12, 0.02, 0.12), 0.84)
	c15 := colorful.Hsl(dh, clamp(ds*0.08, 0.01, 0.08), 0.97)

	accents := filterAccents(colors)
	if len(accents) < 2 {
		accents = append(accents, colorful.Hsl(dh, 0.6, 0.5), colorful.Hsl(math.Mod(dh+180, 360), 0.6, 0.5))
	}

	var normalAccents [6]colorful.Color
	var brightAccents [6]colorful.Color
	var cursor colorful.Color

	switch opts.Mode {
	case "ansi", "spectrum":
		ansiHues := []float64{0, 120, 55, 225, 300, 180}
		for i, targetHue := range ansiHues {
			bestIdx := 0
			bestDist := 999.0
			for j, acc := range accents {
				h, _, _ := acc.Hsl()
				dist := hueDist(h, targetHue)
				if dist < bestDist {
					bestDist = dist
					bestIdx = j
				}
			}
			acc := accents[bestIdx]
			ah, as, al := acc.Hsl()
			if bestDist > 50 {
				ah = targetHue
			}
			as = clamp(as*1.1, 0.45, 0.85)
			al = clamp(al, 0.45, 0.65)
			normalAccents[i] = colorful.Hsl(ah, as, al)
			brightAccents[i] = colorful.Hsl(ah, clamp(as*1.05, 0.45, 0.95), clamp(al+0.18, 0.65, 0.88))
		}
		cursor = brightAccents[3]

	case "dominant", "wal":
		sorted := make([]colorful.Color, len(accents))
		copy(sorted, accents)
		sort.Slice(sorted, func(i, j int) bool {
			_, _, li := sorted[i].Hsl()
			_, _, lj := sorted[j].Hsl()
			return li < lj
		})
		for i := 0; i < 6; i++ {
			c := sorted[i%len(sorted)]
			h, s, l := c.Hsl()
			normalAccents[i] = colorful.Hsl(h, clamp(s*1.1, 0.4, 0.85), clamp(l, 0.45, 0.65))
			brightAccents[i] = colorful.Hsl(h, clamp(s*1.05, 0.4, 0.95), clamp(l+0.18, 0.65, 0.88))
		}
		cursor = normalAccents[1]

	default:
		bestA := accents[0]
		bestB := accents[1]
		maxDist := 0.0
		for i := 0; i < len(accents); i++ {
			for j := i + 1; j < len(accents); j++ {
				d := accents[i].DistanceLab(accents[j])
				if d > maxDist {
					maxDist = d
					bestA = accents[i]
					bestB = accents[j]
				}
			}
		}

		ha, sa, la := bestA.Hsl()
		hb, sb, lb := bestB.Hsl()

		acc1 := colorful.Hsl(ha, clamp(sa*1.15, 0.45, 0.85), clamp(la, 0.45, 0.62))
		acc2 := colorful.Hsl(hb, clamp(sb*1.15, 0.45, 0.85), clamp(lb, 0.45, 0.62))

		bacc1 := colorful.Hsl(ha, clamp(sa*1.05, 0.45, 0.95), clamp(la+0.20, 0.65, 0.90))
		bacc2 := colorful.Hsl(hb, clamp(sb*1.05, 0.45, 0.95), clamp(lb+0.20, 0.65, 0.90))

		normalAccents[0] = acc1
		normalAccents[1] = acc1
		normalAccents[2] = acc1
		normalAccents[3] = acc2
		normalAccents[4] = acc2
		normalAccents[5] = acc2

		brightAccents[0] = bacc1
		brightAccents[1] = bacc1
		brightAccents[2] = bacc1
		brightAccents[3] = bacc2
		brightAccents[4] = bacc2
		brightAccents[5] = bacc2

		cursor = acc2
	}

	return Theme{
		Background: bg.Hex(),
		Foreground: fg.Hex(),
		Cursor:     cursor.Hex(),
		Palette: [16]string{
			c0.Hex(),
			normalAccents[0].Hex(),
			normalAccents[1].Hex(),
			normalAccents[2].Hex(),
			normalAccents[3].Hex(),
			normalAccents[4].Hex(),
			normalAccents[5].Hex(),
			c7.Hex(),
			c8.Hex(),
			brightAccents[0].Hex(),
			brightAccents[1].Hex(),
			brightAccents[2].Hex(),
			brightAccents[3].Hex(),
			brightAccents[4].Hex(),
			brightAccents[5].Hex(),
			c15.Hex(),
		},
	}
}

func generateLightTheme(colors []colorful.Color, dom colorful.Color, opts Options) Theme {
	dh, ds, _ := dom.Hsl()

	bg := colorful.Hsl(dh, clamp(ds*0.15, 0.02, 0.12), 0.96)
	fg := colorful.Hsl(dh, clamp(ds*0.25, 0.04, 0.20), 0.15)
	c0 := colorful.Hsl(dh, clamp(ds*0.25, 0.04, 0.20), 0.20)
	c8 := colorful.Hsl(dh, clamp(ds*0.20, 0.04, 0.16), 0.42)
	c7 := colorful.Hsl(dh, clamp(ds*0.15, 0.02, 0.12), 0.82)
	c15 := colorful.Hsl(dh, clamp(ds*0.08, 0.01, 0.08), 0.08)

	accents := filterAccents(colors)
	if len(accents) < 2 {
		accents = append(accents, colorful.Hsl(dh, 0.6, 0.5), colorful.Hsl(math.Mod(dh+180, 360), 0.6, 0.5))
	}

	var normalAccents [6]colorful.Color
	var brightAccents [6]colorful.Color
	var cursor colorful.Color

	switch opts.Mode {
	case "ansi", "spectrum":
		ansiHues := []float64{0, 120, 55, 225, 300, 180}
		for i, targetHue := range ansiHues {
			bestIdx := 0
			bestDist := 999.0
			for j, acc := range accents {
				h, _, _ := acc.Hsl()
				dist := hueDist(h, targetHue)
				if dist < bestDist {
					bestDist = dist
					bestIdx = j
				}
			}
			acc := accents[bestIdx]
			ah, as, al := acc.Hsl()
			if bestDist > 50 {
				ah = targetHue
			}
			as = clamp(as*1.1, 0.5, 0.9)
			al = clamp(al, 0.35, 0.50)
			normalAccents[i] = colorful.Hsl(ah, as, al)
			brightAccents[i] = colorful.Hsl(ah, clamp(as*0.9, 0.4, 0.8), clamp(al+0.15, 0.45, 0.65))
		}
		cursor = normalAccents[3]

	case "dominant", "wal":
		sorted := make([]colorful.Color, len(accents))
		copy(sorted, accents)
		sort.Slice(sorted, func(i, j int) bool {
			_, _, li := sorted[i].Hsl()
			_, _, lj := sorted[j].Hsl()
			return li < lj
		})
		for i := 0; i < 6; i++ {
			c := sorted[i%len(sorted)]
			h, s, l := c.Hsl()
			normalAccents[i] = colorful.Hsl(h, clamp(s*1.1, 0.5, 0.9), clamp(l, 0.35, 0.50))
			brightAccents[i] = colorful.Hsl(h, clamp(s*0.9, 0.4, 0.8), clamp(l+0.15, 0.45, 0.65))
		}
		cursor = normalAccents[1]

	default:
		bestA := accents[0]
		bestB := accents[1]
		maxDist := 0.0
		for i := 0; i < len(accents); i++ {
			for j := i + 1; j < len(accents); j++ {
				d := accents[i].DistanceLab(accents[j])
				if d > maxDist {
					maxDist = d
					bestA = accents[i]
					bestB = accents[j]
				}
			}
		}

		ha, sa, _ := bestA.Hsl()
		hb, sb, _ := bestB.Hsl()

		acc1 := colorful.Hsl(ha, clamp(sa*1.15, 0.5, 0.9), 0.40)
		acc2 := colorful.Hsl(hb, clamp(sb*1.15, 0.5, 0.9), 0.40)

		bacc1 := colorful.Hsl(ha, clamp(sa*0.95, 0.4, 0.8), 0.55)
		bacc2 := colorful.Hsl(hb, clamp(sb*0.95, 0.4, 0.8), 0.55)

		normalAccents[0] = acc1
		normalAccents[1] = acc1
		normalAccents[2] = acc1
		normalAccents[3] = acc2
		normalAccents[4] = acc2
		normalAccents[5] = acc2

		brightAccents[0] = bacc1
		brightAccents[1] = bacc1
		brightAccents[2] = bacc1
		brightAccents[3] = bacc2
		brightAccents[4] = bacc2
		brightAccents[5] = bacc2

		cursor = acc2
	}

	return Theme{
		Background: bg.Hex(),
		Foreground: fg.Hex(),
		Cursor:     cursor.Hex(),
		Palette: [16]string{
			c0.Hex(),
			normalAccents[0].Hex(),
			normalAccents[1].Hex(),
			normalAccents[2].Hex(),
			normalAccents[3].Hex(),
			normalAccents[4].Hex(),
			normalAccents[5].Hex(),
			c7.Hex(),
			c8.Hex(),
			brightAccents[0].Hex(),
			brightAccents[1].Hex(),
			brightAccents[2].Hex(),
			brightAccents[3].Hex(),
			brightAccents[4].Hex(),
			brightAccents[5].Hex(),
			c15.Hex(),
		},
	}
}
