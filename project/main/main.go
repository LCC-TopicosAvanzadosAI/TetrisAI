package main

import (
	"fmt"
	mn "github.com/TetrisAI/project/menu"
	"github.com/TetrisAI/project/sound"
	"github.com/TetrisAI/project/tetris"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

var condicion = "gurrola"

func main() {
	for {
		pixelgl.Run(run) //para run la wea gráfica
		if condicion != "quit" {
			break
		}
		condicion = "fran"
	}

}

func run() {
	//Window creation
	windowWidth := 765.0
	windowHeight := 450.0
	rand.Seed(time.Now().UTC().UnixNano())
	cfg := pixelgl.WindowConfig{
		Title:  "Tetris",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Black)
	tetris := tetris.Tetris{}
	switch mn.DisplayMenu(win, windowWidth, windowHeight) {
	case "Learn":
		fmt.Println("learning")
		tetris.Learn(win, cfg)
	case "Play":
		tetris.New(win, cfg)
		sound.Play()
		for {
			opcion := tetris.Play()
			if opcion == "quit" {
				sound.Stop()
				condicion = "quit"
				break
			}
			sound.Play()

		}
	}
}
