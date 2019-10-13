package instagram

import (
	"bytes"
	"fmt"
	"image"
	"io"

	resizer "github.com/nfnt/resize"
	"github.com/oliamb/cutter"

	_ "image/jpeg"
	_ "image/png"
)

const (
	maxWidth = 1080
)

var (
	// prefer to split into less images with a taller height
	possibleSplitRatios = []float64{1.25, 1.00}
)

type instagram struct {
	img       image.Image
	extension string
}

type PreparedImage struct {
	Reader    []io.Reader
	Extension string
	Width     int
	Height    int
}

func Prepare(imageReader io.Reader) (*PreparedImage, error) {
	insta, err := decode(imageReader)
	if err != nil {
		return nil, err
	}

	imgs, err := split(insta.img, possibleSplitRatios)
	if err != nil {
		return nil, err
	}

	out := &PreparedImage{
		Reader:    []io.Reader{},
		Extension: insta.extension,
	}
	for n := range imgs {
		resized, width, height := resize(imgs[n], maxWidth)
		if err != nil {
			return nil, err
		}
		writer, err := encode(resized, insta.extension)
		if err != nil {
			return nil, err
		}
		out.Width = width
		out.Height = height
		out.Reader = append(out.Reader, writer)
	}

	return out, nil
}

func decode(imageReader io.Reader) (*instagram, error) {
	img, extension, err := image.Decode(imageReader)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}

	return &instagram{
		img:       img,
		extension: extension,
	}, nil
}

func split(img image.Image, ratios []float64) ([]image.Image, error) {
	var splitImgs []image.Image

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	// check if the image needs to be split for multi pictures
	for _, ratio := range ratios {
		splits := canSplitInto(width, height, ratio)
		if splits == 1 {
			splitImgs = append(splitImgs, img)
			break
		} else if splits > 1 {
			for n := 1; n <= splits; n++ {
				croppedImg, err := cutter.Crop(img, cutter.Config{
					Width:  width / int(splits),
					Height: height,
					Anchor: image.Point{X: width / splits * (n - 1), Y: 0},
				})
				if err != nil {
					return nil, fmt.Errorf("could not crop image: %v", err)
				}
				splitImgs = append(splitImgs, croppedImg)
			}
			break
		}
	}

	if len(splitImgs) == 0 {
		return nil, fmt.Errorf("the provided image is is not a supported aspect ratio")
	}

	return splitImgs, nil
}

func resize(img image.Image, maxWidth uint) (image.Image, int, int) {
	resizedImg := resizer.Resize(maxWidth, 0, img, resizer.Lanczos3)
	return resizedImg, resizedImg.Bounds().Dx(), resizedImg.Bounds().Dy()
}

func encode(img image.Image, extension string) (io.Reader, error) {
	encoded := &bytes.Buffer{}
	if err := encodeImage(encoded, img, extension); err != nil {
		return nil, fmt.Errorf("could not encode image: %v", err)
	}

	return encoded, nil
}

// canSplitInto will calculate if the image should be split into multiple images
// returning 0 if the width * ratio is not a multiple of the height
func canSplitInto(width, height int, ratio float64) int {
	var splits int
	if width%int(float64(height)/ratio) == 0 {
		splits = int(ratio*float64(width)) / height
	}
	return splits
}
