package tetris

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/png"
	//"math"
	"encoding/gob"
	"github.com/TetrisAI/project/agent"
	gb "github.com/TetrisAI/project/gameboard"
	hp "github.com/TetrisAI/project/helper"
	_ "math"
	"os"
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
	time_learn   int
}

func (t *Tetris) New(win *pixelgl.Window, cfg pixelgl.WindowConfig) {
	//fmt.Println("antes de eso")
	//fmt.Println("Win que recivo: ", win)
	t.win = win
	t.cfg = cfg
	t.batch, t.blocksFrames, t.spritesheet, t.background = LoadResources()
	t.game_board = gb.NewGameBoard(t.win, t.batch, t.blocksFrames, t.spritesheet)
	t.game_board.Draw = false
	t.game_over = false
	t.last_time = time.Now()
}

func (t *Tetris) Display() {

	//fmt.Println("hola soy display")
	t.game_board.Draw = true
	t.win.Clear(colornames.Black)

	//fmt.Println("3")
	frames := 0
	//second := time.Tick(time.Second)
	//fmt.Println("5")
	frame := pixel.NewSprite(t.background, t.background.Bounds())
	//fmt.Println("maldito win")
	//fmt.Println("win: ", t.win)
	frame.Draw(t.win, pixel.IM.Moved(t.win.Bounds().Center()))
	//fmt.Println("4")

	last := time.Now()
	//fmt.Println("entrare al for del display")
	for !t.win.Closed() {
		//fmt.Println("Display")
		if time.Since(last).Seconds() > 2 {
			last = time.Now()
			t.win.SetTitle(fmt.Sprintf("%s | FPS: %d", t.cfg.Title, frames))
			frames = 0
		}
		frames++

		t.game_board.Mutex.Lock()
		frame.Draw(t.win, pixel.IM.Moved(t.win.Bounds().Center()))
		t.win.Update()
		t.game_board.Mutex.Unlock()

	}
	//fmt.Println("bye display")
	t.game_board.Game_over = true
}

func (t *Tetris) Learn(win *pixelgl.Window, cfg pixelgl.WindowConfig) {
	fmt.Println("Learning :D")
	t.New(win, cfg)
	go t.Display()
	a := agent.NewState()
	a.SetBoard(t.game_board)
	read_value_function("value.gob", a.GetValueFunction())
	go a.StartTime()
	//steps := 0
	for i := 0; i < 100000; i++ {
		steps := a.Start(1)

		//fmt.Println(float64(steps) / float64(i+1))
		if i%1 == 0 {
			fmt.Println(steps)
			//fmt.Println("Guardando...")
			save_value_function("value.gob", a.GetValueFunction())
		}
	}

}
func read_value_function(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func save_value_function(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}
func (t *Tetris) Play() string {

	fmt.Println("llamare al display")
	go t.Display()
	fmt.Println("ya lo llame")
	go t.game_board.Start()

	for !t.game_board.Game_over {
		//fmt.Println(!t.game_board.Game_over)
		//time.Sleep(100 * time.Millisecond)
		//fmt.Println("prueba 1")
		action := t.Get_action_player()
		fmt.Println(action)
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
	//fmt.Println("amm")
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
