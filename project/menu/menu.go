package menu

import (
	"fmt"
	//"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"math"
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
		for !win.Closed() {
			time.Sleep(time.Second * 1)
			gameBoard.Gravity()
			win.Clear(colornames.Black)
			gameBoard.DisplayBoard(win, blockGen)

		}

	}()

	//movimiento de las piezas
	MovementDelay := 0.0
	moveCounter := 0
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if MovementDelay > 0.0 {
			MovementDelay = math.Max(MovementDelay-dt, 0.0)
		}

		if win.Pressed(pixelgl.KeyRight) && MovementDelay == 0 {
			gameBoard.MovePiece(gb.MoveRight)
			if moveCounter > 0 {
				MovementDelay = 0.1
			} else {
				MovementDelay = 0.5
			}
			moveCounter++
		}
		if win.Pressed(pixelgl.KeyLeft) && MovementDelay == 0 {
			gameBoard.MovePiece(gb.MoveLeft)
			if moveCounter > 0 {
				MovementDelay = 0.1
			} else {
				MovementDelay = 0.5
			}
			moveCounter++
		}
		if win.JustPressed(pixelgl.KeyUp) {
			gameBoard.RotatePiece()
		}
		if win.JustPressed(pixelgl.KeyDown) {
			gameBoard.MovePiece(gb.MoveToBottom)
		}

		if !win.Pressed(pixelgl.KeyRight) && !win.Pressed(pixelgl.KeyLeft) {
			moveCounter = 0
			MovementDelay = 0.0
		}

		win.Update()
		win.Clear(colornames.Black)
		gameBoard.DisplayBoard(win, blockGen)

		//fmt.Println(gameBoard)
	}

}