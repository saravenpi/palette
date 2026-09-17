package palette

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/cascax/colorthief-go"
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

	pal, err := colorthief.GetPalette(img, count)
	if err != nil {
		return nil, colorful.Color{}, fmt.Errorf("failed to extract palette: %w", err)
	}

	var colors []colorful.Color
	for _, c := range pal {
		col, ok := colorful.MakeColor(c)
		if ok {
			colors = append(colors, col)
		}
	}

	domCol, err := colorthief.GetColor(img)
	var dom colorful.Color
	if err == nil {
		dom, _ = colorful.MakeColor(domCol)
	} else if len(colors) > 0 {
		dom = colors[0]
	} else {
		dom = colorful.Color{R: 0.5, G: 0.5, B: 0.5}
	}

	return colors, dom, nil
}
