package tetris

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	//"math"
	gb "github.com/TetrisAI/project/gameboard"
	hp "github.com/TetrisAI/project/helper"
	"math"
	"time"
)

func LoadResources() (*pixel.Batch, []pixel.Rect, pixel.Picture, pixel.Picture) {
	spritesheet, _ := hp.LoadPicture("./../../resources/blocks.png")
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	blocksFrames := hp.NewBlockFrame(spritesheet)
	background, _ := hp.LoadPicture("./../../resources/marco.png")

	return batch, blocksFrames, spritesheet, background

}

func Play(win *pixelgl.Window, cfg pixelgl.WindowConfig) string {

	frames := 0
	second := time.Tick(time.Second)
	batch, blocksFrames, spritesheet, background := LoadResources()
	gameBoard := gb.NewGameBoard(win, batch, blocksFrames, spritesheet)

	win.Clear(colornames.Black)

	gameBoard.UpdateScore()
	gameBoard.UpdateBoard()
	gameBoard.AddPiece()

	MovementDelay := 0.0
	moveCounter := 0
	last := time.Now()

	var gravityTimer float64
	//var baseSpeed float64 = 0.8
	var gravitySpeed float64 = 1.0

	frame := pixel.NewSprite(background, background.Bounds())
	frame.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		gravityTimer += dt

		if gravityTimer > gravitySpeed && !win.Pressed(pixelgl.KeyDown) {
			gravityTimer -= gravitySpeed
			quit := gameBoard.Gravity()
			if quit == "quit" {
				return "quit"
			}
		}

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

			quit := gameBoard.MovePiece(gb.MoveToBottom)
			if quit == "quit" {
				return "quit"
			}
			if moveCounter > 0 {
				MovementDelay = 0.05
			} else {
				MovementDelay = 0.1
			}
			moveCounter++
		}
		if win.JustPressed(pixelgl.KeySpace) {
			quit := gameBoard.MoveToBottom1()
			if quit == "quit" {
				return "quit"
			}
		}

		if !win.Pressed(pixelgl.KeyRight) && !win.Pressed(pixelgl.KeyLeft) && !win.Pressed(pixelgl.KeyDown) {
			moveCounter = 0
			MovementDelay = 0.0
		}

		if win.Pressed(pixelgl.KeyR) {
			return "r"

		}

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		frame.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
	return "quit"

}
