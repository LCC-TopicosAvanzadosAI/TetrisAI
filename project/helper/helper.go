package helper

import (
	"github.com/faiface/pixel"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"io/ioutil"
	"os"
)

func NewBlockFrame(spritesheet pixel.Picture) []pixel.Rect {

	var blocksFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 40 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 40 {
			blocksFrames = append(blocksFrames, pixel.R(x, y, x+40, y+40))
		}
	}
	return blocksFrames
}

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
