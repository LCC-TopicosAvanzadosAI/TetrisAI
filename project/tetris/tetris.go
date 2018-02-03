package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"time"

	gb "github.com/TetrisAI/project/gameboard"
	hp "github.com/TetrisAI/project/helper"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	windowWidth := 765.0
	windowHeight := 450.0
	cfg := pixelgl.WindowConfig{
		Title:  "Blockfall",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var blockGen func(int) pixel.Picture

	blockGen, err = hp.LoadSpriteSheet("./../../resources/blocks.png", 2, 8)

	if err != nil {
		panic(err)
	}
	gameBoard := gb.NewGameBoard()
	gameBoard.AddPiece()

	fmt.Println(gameBoard)

	//gameBoard.DisplayBoard(win, blockGen)
	for {
		win.Update()
		gameBoard.Gravity()
		win.Clear(colornames.Black)
		gameBoard.DisplayBoard(win, blockGen)
		time.Sleep(time.Second * 1)

		//fmt.Println(gameBoard)
	}

}