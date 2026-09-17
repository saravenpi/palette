package palette

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"sort"

	matcolor "github.com/Nadim147c/material/v3/color"
	"github.com/Nadim147c/material/v3/num"
	"github.com/Nadim147c/material/v3/quantizer"
	"github.com/Nadim147c/material/v3/score"
	"github.com/lucasb-eyer/go-colorful"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}
	return img, nil
}

func Downscale(img image.Image, maxDim int) image.Image {
	if maxDim <= 0 {
		return img
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w <= maxDim && h <= maxDim {
		return img
	}

	var nw, nh int
	if w > h {
		nw = maxDim
		nh = (h * maxDim) / w
	} else {
		nh = maxDim
		nw = (w * maxDim) / h
	}

	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	for y := 0; y < nh; y++ {
		sy := bounds.Min.Y + (y*h)/nh
		for x := 0; x < nw; x++ {
			sx := bounds.Min.X + (x*w)/nw
			dst.Set(x, y, img.At(sx, sy))
		}
	}
	return dst
}

func ExtractColors(img image.Image, count int) ([]colorful.Color, colorful.Color, error) {
	if count < 8 {
		count = 16
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w == 0 || h == 0 {
		return nil, colorful.Color{}, fmt.Errorf("empty image bounds")
	}

	scaled := Downscale(img, 400)
	sb := scaled.Bounds()
	sw := sb.Dx()
	sh := sb.Dy()

	pixels := make([]matcolor.ARGB, 0, sw*sh)
	for y := sb.Min.Y; y < sb.Max.Y; y++ {
		for x := sb.Min.X; x < sb.Max.X; x++ {
			c := scaled.At(x, y)
			_, _, _, a := c.RGBA()
			if a < 32768 {
				continue
			}
			pixels = append(pixels, matcolor.ARGBFromInterface(c))
		}
	}
	if len(pixels) == 0 {
		for y := sb.Min.Y; y < sb.Max.Y; y++ {
			for x := sb.Min.X; x < sb.Max.X; x++ {
				pixels = append(pixels, matcolor.ARGBFromInterface(scaled.At(x, y)))
			}
		}
	}

	quantized := quantizer.QuantizeCelebi(pixels, 128)
	scored := score.Score(quantized, score.WithLimit(count))

	type candidate struct {
		argb  matcolor.ARGB
		hct   matcolor.Hct
		score float64
	}

	totalPop := 0
	for _, p := range quantized {
		totalPop += p
	}

	var cands []candidate
	for argb, pop := range quantized {
		hct := argb.ToHct()
		prop := float64(pop) / float64(totalPop)
		sc := prop*50.0 + hct.Chroma*1.5
		cands = append(cands, candidate{argb, hct, sc})
	}
	sort.Slice(cands, func(i, j int) bool {
		return cands[i].score > cands[j].score
	})

	var diverse []candidate
	for _, c := range cands {
		tooClose := false
		for _, d := range diverse {
			if num.DifferenceDegrees(c.hct.Hue, d.hct.Hue) < 22 && math.Abs(c.hct.Tone-d.hct.Tone) < 18 {
				tooClose = true
				break
			}
		}
		if !tooClose {
			diverse = append(diverse, c)
		}
	}

	var colors []colorful.Color
	var dom colorful.Color

	if len(scored) > 0 {
		domARGB := scored[0]
		dom = colorful.Color{
			R: float64(domARGB.Red()) / 255.0,
			G: float64(domARGB.Green()) / 255.0,
			B: float64(domARGB.Blue()) / 255.0,
		}
	} else if len(diverse) > 0 {
		domARGB := diverse[0].argb
		dom = colorful.Color{
			R: float64(domARGB.Red()) / 255.0,
			G: float64(domARGB.Green()) / 255.0,
			B: float64(domARGB.Blue()) / 255.0,
		}
	} else {
		dom = colorful.Color{R: 0.26, G: 0.52, B: 0.96}
	}

	for _, d := range diverse {
		colors = append(colors, colorful.Color{
			R: float64(d.argb.Red()) / 255.0,
			G: float64(d.argb.Green()) / 255.0,
			B: float64(d.argb.Blue()) / 255.0,
		})
	}

	for _, sc := range scored {
		if len(colors) >= count {
			break
		}
		c := colorful.Color{
			R: float64(sc.Red()) / 255.0,
			G: float64(sc.Green()) / 255.0,
			B: float64(sc.Blue()) / 255.0,
		}
		exists := false
		for _, existing := range colors {
			if existing.DistanceLab(c) < 0.05 {
				exists = true
				break
			}
		}
		if !exists {
			colors = append(colors, c)
		}
	}

	if len(colors) == 0 {
		colors = append(colors, dom)
	}

	return colors, dom, nil
}
