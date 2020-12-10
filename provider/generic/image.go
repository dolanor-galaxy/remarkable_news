package generic

import (
	"context"
	"errors"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

const (
	reWidth  = 1404
	reHeight = 1872
)

// Process the picture ...
func (p *Provider) Process(ctx context.Context, img image.Image) (image.Image, error) {

	switch p.Mode {
	case "fill":
		// scale image to remarkable width
		// imaging resize is slow for some reason, use other library
		// img = imaging.Resize(img, re_width, 0, imaging.Linear)
		img = resize.Resize(uint(reWidth), 0, img, resize.Bilinear)
		// cut off parts of image that overflow
		img = imaging.Crop(img, image.Rect(0, 0, reWidth, reHeight))
	case "center":
	default:
		return nil, errors.New("Invalid mode")
	}
	if p.Scale != 1 {
		imgWidth := float64(img.Bounds().Max.X)
		img = resize.Resize(uint(p.Scale*imgWidth), 0, img, resize.Bilinear)

	}

	// put image in center of screen
	background := imaging.New(
		reWidth,
		reHeight,
		color.RGBA{255, 255, 255, 255},
	)
	img = imaging.PasteCenter(background, img)

	return img, nil
}
