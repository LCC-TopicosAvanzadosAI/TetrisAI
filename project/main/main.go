package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	mn "github.com/TetrisAI/project/menu"
	"github.com/TetrisAI/project/sound"
	"github.com/TetrisAI/project/tetris"
)

func main() {
	pixelgl.Run(run) //para correr la wea gráfica
}

func run() {
	//Window creation
	windowWidth := 765.0
	windowHeight := 450.0

	cfg := pixelgl.WindowConfig{
		Title:  "Tetris",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Black)

	switch mn.DisplayMenu(win, windowWidth, windowHeight) {
	case "Play":
		sound.Play()
		tetris.Play(win, cfg)
	}

}
