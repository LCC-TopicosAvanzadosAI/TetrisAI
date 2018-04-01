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
	_ "math"
	"time"
)

func LoadResources() (*pixel.Batch, []pixel.Rect, pixel.Picture, pixel.Picture) {
	spritesheet, _ := hp.LoadPicture("./../../resources/blocks.png")
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	blocksFrames := hp.NewBlockFrame(spritesheet)
	background, _ := hp.LoadPicture("./../../resources/marco.png")

	return batch, blocksFrames, spritesheet, background

}

type BlockColor int

type Tetris struct {
	game_board   gb.Board
	win          *pixelgl.Window
	cfg          pixelgl.WindowConfig
	batch        *pixel.Batch
	blocksFrames []pixel.Rect
	spritesheet  pixel.Picture
	background   pixel.Picture
	game_over    bool
	last_time    time.Time
	last_action  string
}

func (t *Tetris) New(win *pixelgl.Window, cfg pixelgl.WindowConfig) {
	t.win = win
	t.cfg = cfg
	t.batch, t.blocksFrames, t.spritesheet, t.background = LoadResources()
	t.game_board = gb.NewGameBoard(t.win, t.batch, t.blocksFrames, t.spritesheet)
	t.game_over = false
	t.last_time = time.Now()
}

func (t *Tetris) Display() {

	t.win.Clear(colornames.Black)

	frames := 0
	second := time.Tick(time.Second)

	frame := pixel.NewSprite(t.background, t.background.Bounds())
	frame.Draw(t.win, pixel.IM.Moved(t.win.Bounds().Center()))

	//last := time.Now()

	for !t.win.Closed() {
		frames++
		select {
		case <-second:
			t.win.SetTitle(fmt.Sprintf("%s | FPS: %d", t.cfg.Title, frames))
			frames = 0
		}
		t.game_board.Mutex.Lock()
		frame.Draw(t.win, pixel.IM.Moved(t.win.Bounds().Center()))
		t.win.Update()
		t.game_board.Mutex.Unlock()

	}
	t.game_board.Game_over = true
}

func (t *Tetris) Play() string {

	go t.Display()
	go t.game_board.Start()

	for !t.game_board.Game_over {
		action := t.Get_action_player()
		t.Take_action(action)
	}
	return "quit"
}

func (t *Tetris) Take_action(action string) {
	if time.Since(t.last_time).Seconds() < 0.3 && t.last_action == action {
		return
	}
	t.last_time = time.Now()
	t.last_action = action
	t.game_board.Mutex.Lock()
	switch action {
	case "KeyDown":
		t.game_board.Gravity()
	case "KeyRight":
		t.game_board.MovePiece(gb.MoveRight)
	case "KeyLeft":
		t.game_board.MovePiece(gb.MoveLeft)
	case "KeyUp":
		t.game_board.RotatePiece()
	case "KeySpace":
		t.game_board.MoveToBottom1()
	}
	t.game_board.Mutex.Unlock()
}

func (t *Tetris) Get_action_player() string {
	if t.win.Pressed(pixelgl.KeyDown) {
		return "KeyDown"
	}
	if t.win.Pressed(pixelgl.KeyRight) {
		return "KeyRight"
	}
	if t.win.Pressed(pixelgl.KeyLeft) {
		return "KeyLeft"
	}
	if t.win.Pressed(pixelgl.KeyUp) {
		return "KeyUp"
	}
	if t.win.Pressed(pixelgl.KeySpace) {
		return "KeySpace"
	}
	return "wea"
}
