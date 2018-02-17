package helper

import (
	"errors"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"io/ioutil"
	"os"
)

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func LoadSpriteSheet(path string, row, col int) (func(int) pixel.Picture, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Load Image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Check if tile is square
	b := img.Bounds()
	if b.Max.X/col != b.Max.Y/row {
		fmt.Println("width/col = ", b.Max.X, ", height/row = ", b.Max.Y)
		return nil, errors.New(fmt.Sprintf("Invalid dimensions (%d, %d) for sprite sheet %s\n", row, col, path))
	}

	tileSize := b.Max.X / col

	return func(i int) pixel.Picture {
		if i < 0 || i >= row*col {
			panic("Index out of bounds for sprite sheet")
		}
		r := i / col
		c := i % col

		subImage := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(c*tileSize, r*tileSize, (c+1)*tileSize, (r+1)*tileSize))
		return pixel.PictureDataFromImage(subImage)
	}, nil
}

//Necessary function to load fonts
func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
