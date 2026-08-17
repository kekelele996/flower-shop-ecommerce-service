package util

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// PlaceholderPNG 生成纯色占位 PNG 图片（种子商品图）。
func PlaceholderPNG(w, h int, r, g, b uint8) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	step := w / 8
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if (x/step)%2 == 0 {
				img.Set(x, y, color.RGBA{r, g, b, 255})
			} else {
				img.Set(x, y, color.RGBA{r - 40, g - 30, b - 20, 255})
			}
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
