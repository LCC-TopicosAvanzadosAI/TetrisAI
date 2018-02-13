package menu

import (
	"fmt"
	//"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"time"

	gb "github.com/TetrisAI/project/gameboard"
	hp "github.com/TetrisAI/project/helper"
)

const (
	BoardRows = 22
	BoardCols = 10
)

type Menu struct {
	aux int
}

func NewMenu() Menu {
	return Menu{
		aux: 1,
	}
}

func (m *Menu) DisplayMenu(win *pixelgl.Window) {

}

func (m *Menu) Jugar(win *pixelgl.Window) {

	//Carga de imágen
	blockGen, err := hp.LoadSpriteSheet("./../../resources/blocks.png", 2, 8)

	if err != nil {
		panic(err)
	}

	gameBoard := gb.NewGameBoard()

	gameBoard.AddPiece()

	fmt.Println(gameBoard)

	gameBoard.DisplayBoard(win, blockGen)
	win.Clear(colornames.Black)

	go func() {
		for {
			time.Sleep(time.Second * 1)
			gameBoard.Gravity()
			win.Clear(colornames.Black)
			gameBoard.DisplayBoard(win, blockGen)

		}

	}()

	//movimiento de las piezas
	for {
		win.Update()
		win.Clear(colornames.Black)
		gameBoard.DisplayBoard(win, blockGen)

		if win.Pressed(pixelgl.KeyRight) {
			gameBoard.MovePiece(gb.MoveRight)
		}
		if win.Pressed(pixelgl.KeyLeft) {
			gameBoard.MovePiece(gb.MoveLeft)
		}
		if win.Pressed(pixelgl.KeyUp) {
			gameBoard.RotatePiece()
		}
		if win.Pressed(pixelgl.KeyDown) {
			gameBoard.MovePiece(gb.MoveToBottom)
		}

		//fmt.Println(gameBoard)
	}

}
