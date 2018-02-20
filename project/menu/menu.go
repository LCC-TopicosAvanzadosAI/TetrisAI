package menu

import (
	"fmt"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	"math"
	"os"
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

	//Carga de imágen de la pieza
	blockGen, err := hp.LoadSpriteSheet("./../../resources/blocks.png", 2, 8)
	if err != nil {
		panic(err)
	}

	pic, err := hp.LoadPicture("./../../resources/marco.png")
	if err != nil {
		panic(err)
	}

	//tetris audio
	f, _ := os.Open("./../../resources/korobeiniki.mp3")
	s, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(s)
	//s.Loop(-1, s)

	gameBoard := gb.NewGameBoard()

	gameBoard.AddPiece()

	fmt.Println(gameBoard)

	gameBoard.DisplayBoard(win, blockGen)

	frame := pixel.NewSprite(pic, pic.Bounds())
	frame.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	win.Clear(colornames.Black)

	go func() {
		for !win.Closed() {

			time.Sleep(time.Second * 1)
			gameBoard.Gravity()
			win.Clear(colornames.Black)
			//gameBoard.DisplayBoard(win, blockGen)
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

		if win.Pressed(pixelgl.KeyDown) && MovementDelay == 0 {

			gameBoard.MovePiece(gb.MoveToBottom)
			//gameBoard.MovePiece(gb.MoveDown)
			if moveCounter > 0 {
				MovementDelay = 0.1
			} else {
				MovementDelay = 0.5
			}
			moveCounter++
		}
		if win.JustPressed(pixelgl.KeySpace) {
			gameBoard.MoveToBottom1()
		}

		if !win.Pressed(pixelgl.KeyRight) && !win.Pressed(pixelgl.KeyLeft) {
			moveCounter = 0
			MovementDelay = 0.0
		}
		win.Update()
		win.Clear(colornames.Black)
		frame.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		gameBoard.DisplayBoard(win, blockGen)

		//fmt.Println(gameBoard)
	}

}
