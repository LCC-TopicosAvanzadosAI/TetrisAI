package gameboard

import (
	"fmt"
	hp "github.com/TetrisAI/project/helper"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	BoardRows = 22
	BoardCols = 10
)

type Movement int

const (
	MoveDown     Movement = -1
	MoveLeft              = -1
	MoveRight             = 1
	nothing               = 0
	MoveToBottom          = -2
)

//This is the most simple unit. It means that we have a 'Block' on row x and col y
type Block struct {
	row int
	col int
}

func (b *Block) GetRow() int {
	return b.row
}
func (b *Block) GetCol() int {
	return b.col
}
func (p *Piece) GetId() int {
	return int(p.id)
}

//We can make a shape made of Blocks, luckily in Tetris all Shapes have 4 blocks, so we can represent it as an array of 4 Points.
type Shape [4]Block

type IdPiece int

const (
	IdJPiece IdPiece = iota
	IdLPiece
	IdSPiece
	IdTPiece
	IdZPiece
	IdIPiece
	IdOPiece
)

//And each Shape has 4 different positions on the board, that is when you rotate it, so a Piece is made of 4 different shapes.
type Piece struct {
	Piece_ [4]Block
	color  BlockColor
	id     IdPiece
}

//I'm still thinking about this, we could use a binary representation but that's really ugly and it will be really difficult to update in the future
//This way is really easy to understand but the code there is really ugly. So I don't know. I'll let my team decide...

func JPiece(c BlockColor) Piece {
	return Piece{
		Piece_: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 0, col: 1},
			Block{row: 0, col: 0},
			Block{row: 0, col: 2},
		},
		color: 5,
		id:    IdJPiece,
	}
}
func LPiece(c BlockColor) Piece {
	return Piece{
		Piece_: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
			Block{row: 0, col: 0},
		},
		color: 13,
		id:    IdLPiece,
	}
}
func SPiece(c BlockColor) Piece {
	return Piece{
		Piece_: [4]Block{
			Block{row: 0, col: 0},
			Block{row: 0, col: 1},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
		},
		color: 15,
		id:    IdSPiece,
	}
}
func IPiece(c BlockColor) Piece {
	return Piece{
		Piece_: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
			Block{row: 1, col: 3},
		},
		color: 3,
		id:    IdIPiece,
	}
}
func ZPiece(c BlockColor) Piece {
	return Piece{
		Piece_: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 0, col: 1},
			Block{row: 0, col: 2},
		},
		color: 11,
		id:    IdZPiece,
	}
}
func TPiece(c BlockColor) Piece {
	return Piece{
		Piece_: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
			Block{row: 0, col: 1},
		},
		color: 9,
		id:    IdTPiece,
	}
}
func OPiece(c BlockColor) Piece {
	return Piece{
		Piece_: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 0, col: 0},
			Block{row: 0, col: 1},
		},
		color: 7,
		id:    IdOPiece,
	}
}

type BlockColor int

//type Board [BoardRows][BoardCols]BlockColor

// Different values a point on the grid can hold
const (
	Empty BlockColor = iota
	Goluboy
	Siniy
	Pink
	Purple
	Red
	Yellow
	Green
	Gray
	GoluboySpecial
	SiniySpecial
	PinkSpecial
	PurpleSpecial
	RedSpecial
	YellowSpecial
	GreenSpecial
	GraySpecial
)

func (b *Board) GetGameBoard() [BoardRows][BoardCols]BlockColor {
	return b.gameboard
}

type Board struct {
	gameboard   [BoardRows][BoardCols]BlockColor
	activePiece Piece
	activeShape int
	nextPiece   Piece
	PieceArray  []Piece
	score       int
	ScoreTxt    *text.Text
	win         *pixelgl.Window
	Batch       *pixel.Batch
	//blocks       []*pixel.Sprite
	//matrices     []pixel.Matrix
	blocksFrames []pixel.Rect
	BoardCols    int
	BoardRows    int
	spritesheet  pixel.Picture
	Game_over    bool
	time_gravity int
	Mutex        *sync.Mutex
	Draw         bool
}

func (b *Board) GetWin() *pixelgl.Window {
	return b.win
}

func (b *Board) GetActivePiece() Piece {
	return b.activePiece
}

func (b *Board) Score() int {
	return b.score
}

func (b *Board) Start() {
	for !b.Game_over {
		//fmt.Println("nani?")
		b.time_gravity++
		if b.time_gravity == 9 {

			b.time_gravity = 0

			b.Mutex.Lock()
			b.Gravity()
			b.Mutex.Unlock()

		}
		time.Sleep(100 * time.Millisecond)
	}
}

func NewGameBoard(_win *pixelgl.Window, _batch *pixel.Batch, _blocksFrames []pixel.Rect, _spritesheet pixel.Picture) Board {

	face, err := hp.LoadTTF("./../../resources/saarland.ttf", 40) //Loading font and size-font
	if err != nil {
		panic(err)
	}

	board := Board{
		Game_over: false,
		score:     0,
	}
	//board.Flag_board = true
	//board.Flag_tetris = true
	board.Mutex = &sync.Mutex{}
	board.time_gravity = 0
	board.Game_over = false
	board.win = _win
	board.Draw = true
	board.Batch = _batch
	//board.blocks = _blocks
	//board.matrices = _matrices
	board.blocksFrames = _blocksFrames
	board.spritesheet = _spritesheet
	board.BoardCols = 10
	board.BoardRows = 20

	board.FillArray()

	Atlas := text.NewAtlas(face, text.ASCII)            //Atlas necessary for the font
	board.ScoreTxt = text.New(pixel.V(200, 307), Atlas) //here, I put the coordinates where the
	board.ScoreTxt.Dot.X -= board.ScoreTxt.BoundsOf(strconv.Itoa(board.score)).W()
	board.ScoreTxt.Color = colornames.Lightcyan //text color
	fmt.Fprint(board.ScoreTxt, board.score)
	board.ScoreTxt.Draw(board.win, pixel.IM)

	return board
}

func (b *Board) FillArray() {
	b.PieceArray = append(b.PieceArray, JPiece(1), LPiece(2), SPiece(3), TPiece(4), ZPiece(5), IPiece(6), OPiece(7))
	//b.PieceArray = append(b.PieceArray, OPiece(1), OPiece(1))
	dest := make([]Piece, len(b.PieceArray))
	rand.Seed(time.Now().Unix())
	perm := rand.Perm(len(b.PieceArray))
	for i, v := range perm {
		dest[v] = b.PieceArray[i]
	}

	for i := range perm {
		b.PieceArray[i] = dest[i]
	}
	b.nextPiece = b.PieceArray[0]
}

func (p *Piece) movePiece(rows, cols Movement) {

	for i := 0; i < 4; i++ {
		p.Piece_[i].row += int(rows)
		p.Piece_[i].col += int(cols)
	}
}

func (b *Board) checkCollision(p Piece, _r, _c Movement) bool {
	for i := 0; i < 4; i++ {
		r := p.Piece_[i].row + int(_r)
		c := p.Piece_[i].col + int(_c)
		if r < 0 || r > 21 || c < 0 || c > 9 || b.gameboard[r][c] != Empty {
			return true
		}
	}
	return false
}

func (p *Piece) Rotate() {

	pivot := p.Piece_[1]
	arr := []int{0, 2, 3}
	for _, i := range arr {
		dRow := pivot.row - p.Piece_[i].row
		dCol := pivot.col - p.Piece_[i].col
		p.Piece_[i].row = pivot.row + (dCol * -1)
		p.Piece_[i].col = pivot.col + (dRow)
	}
}

func (b *Board) RotatePiece() {
	if b.activePiece.id == IdOPiece {
		return
	}

	b.drawPiece(b.activePiece, Empty)
	b.patchActivePiece(Empty)

	copy := b.activePiece
	b.activePiece.Rotate()

	if b.checkCollision(b.activePiece, nothing, nothing) {
		if !b.checkCollision(b.activePiece, nothing, MoveRight) {
			b.activePiece.movePiece(nothing, MoveRight)
		} else if !b.checkCollision(b.activePiece, nothing, MoveLeft) {
			b.activePiece.movePiece(nothing, MoveLeft)
		} else if !b.checkCollision(b.activePiece, MoveDown, nothing) {
			b.activePiece.movePiece(MoveDown, nothing)
		} else {
			b.activePiece = copy
			b.drawPiece(b.activePiece, b.activePiece.color)
			b.patchActivePiece(b.activePiece.color)
			return
		}
	}

	b.drawPiece(b.activePiece, b.activePiece.color)
	b.patchActivePiece(b.activePiece.color)
}

func (b *Board) MovePiece(dir Movement) string {
	b.drawPiece(b.activePiece, Empty)
	b.patchActivePiece(Empty)

	if dir == MoveToBottom {
		if !b.checkCollision(b.activePiece, MoveDown, nothing) {
			b.activePiece.movePiece(MoveDown, nothing)
		} else {
			return b.Gravity()
		}
	} else if !b.checkCollision(b.activePiece, nothing, dir) {
		b.activePiece.movePiece(nothing, dir)
	}

	b.drawPiece(b.activePiece, b.activePiece.color)
	b.patchActivePiece(b.activePiece.color)
	return ""
}

func (b *Board) MoveToBottom1() string {
	//b.score += 1
	b.UpdateScore()
	b.drawPiece(b.activePiece, Empty)
	b.patchActivePiece(Empty)
	for !b.checkCollision(b.activePiece, MoveDown, nothing) {
		b.activePiece.movePiece(MoveDown, nothing)
	}
	b.drawPiece(b.activePiece, b.activePiece.color)
	b.patchActivePiece(b.activePiece.color)
	return b.Gravity()
}

func (b *Board) Gravity() string {

	b.drawPiece(b.activePiece, Empty)
	b.patchActivePiece(Empty)

	if !b.checkCollision(b.activePiece, MoveDown, nothing) {

		b.activePiece.movePiece(MoveDown, nothing)
		b.drawPiece(b.activePiece, b.activePiece.color)
		b.patchActivePiece(b.activePiece.color)
	} else {

		b.drawPiece(b.activePiece, b.activePiece.color)
		b.patchActivePiece(b.activePiece.color)
		b.checkRowCompletion()

		return b.AddPiece()
	}
	return ""
}

func (b *Board) patchActivePiece(c BlockColor) {
	if b.Draw {
		//fmt.Println("aqui")
		boardBlockSize := 20.0 //win.Bounds().Max.X / 10
		//pic := b.blocksFrames[0] //blockGen(0)
		imgSize := 40 //pic.Bounds().Max.X
		scaleFactor := float64(boardBlockSize) / float64(imgSize)

		for i := 0; i < len(b.activePiece.Piece_); i++ {
			x := float64(b.activePiece.Piece_[i].col)*boardBlockSize + boardBlockSize/2
			y := float64(b.activePiece.Piece_[i].row)*boardBlockSize + boardBlockSize/2

			block := pixel.NewSprite(b.spritesheet, b.blocksFrames[c])
			//b.blocks = append(b.blocks, block)
			block.Draw(b.Batch, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+282, y+25)))
		}
		b.Batch.Draw(b.win)
	}
	//b.win.Update()
}

func (b *Board) checkRowCompletion() {
	rowWasDeleted := true
	var linesFound int
	var deleteRowCt int
	linesFound = 0
	for rowWasDeleted {
		rowWasDeleted = false

		for i := 0; i < 4; i++ {
			r := b.activePiece.Piece_[i].row
			emptyFound := false
			for c := 0; c < 10; c++ {
				if b.gameboard[r][c] == Empty {
					emptyFound = true
					continue
				}
			}
			if !emptyFound {
				fmt.Println("HIZO UNA LINEA WE")
				b.deleteRow(r)
				rowWasDeleted = true
				deleteRowCt++
				linesFound += 1
				b.score += 1
				b.UpdateScore()
			}
		}
		if linesFound == 4 {
			b.score += 4
			b.UpdateScore()
			linesFound = 0
		}
	}

	//fmt.Println("Por tu jugada ", negval)
	//b.UpdateScore()

}

func (b *Board) deleteRow(row int) {

	for r := row; r < 21; r++ {
		for c := 0; c < 10; c++ {
			b.gameboard[r][c] = b.gameboard[r+1][c]
		}
	}
	if b.Draw {
		b.UpdateBoard()
	}
}

func (b *Board) UpdateBoard() {
	if b.Draw {
		b.Batch.Clear()
		boardBlockSize := 20.0 //win.Bounds().Max.X / 10
		imgSize := 40          //pic.Bounds().Max.X
		scaleFactor := float64(boardBlockSize) / float64(imgSize)

		for col := 0; col < b.BoardCols; col++ {
			for row := 0; row < b.BoardRows-2; row++ {
				val := b.gameboard[row][col]
				if val == Empty {
					continue
				}
				x := float64(col)*boardBlockSize + boardBlockSize/2
				y := float64(row)*boardBlockSize + boardBlockSize/2

				block := pixel.NewSprite(b.spritesheet, b.blocksFrames[val])
				//b.blocks = append(b.blocks, block)
				block.Draw(b.Batch, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+282, y+25)))
			}
		}
		b.Batch.Draw(b.win)
		b.win.Update()
	}
}

func (b *Board) drawNextPiece(c BlockColor) {
	if b.Draw {

		boardBlockSize := 20.0
		imgSize := 40 //pic.Bounds().Max.X
		scaleFactor := float64(boardBlockSize) / float64(imgSize)

		for i := 0; i < len(b.nextPiece.Piece_); i++ {
			x := float64(b.nextPiece.Piece_[i].col)*boardBlockSize + boardBlockSize/2
			y := float64(b.nextPiece.Piece_[i].row)*boardBlockSize + boardBlockSize/2

			block := pixel.NewSprite(b.spritesheet, b.blocksFrames[c])
			//b.blocks = append(b.blocks, block)
			block.Draw(b.Batch, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+500, y+285)))
		}

		b.Batch.Draw(b.win)
	}
}

func (b *Board) UpdateScore() {
	b.win.Clear(colornames.Black)
	b.ScoreTxt.Clear()
	fmt.Fprint(b.ScoreTxt, b.score)
	b.ScoreTxt.Draw(b.win, pixel.IM)
}

func (b *Board) SetScore(score int) {
	b.score = score
}

func (b *Board) ResetBoard() {
	for i := range b.gameboard {
		for j := range b.gameboard[i] {
			b.gameboard[i][j] = 0
		}
	}
	b.PieceArray = make([]Piece, 0)
	b.FillArray()

	b.UpdateBoard()

}

func (b *Board) AddPiece() string {
	//b.score++
	b.drawNextPiece(Empty)
	b.nextPiece.movePiece(Movement(18), Movement(4))
	b.activePiece = b.nextPiece

	if b.checkCollision(b.activePiece, 0, 0) {
		b.Game_over = true
		b.score -= 1
		//b.win.Clear(colornames.Black)
		//b.UpdateBoard()
	}
	b.drawPiece(b.activePiece, b.activePiece.color)
	b.patchActivePiece(b.activePiece.color)
	//remove a piece or refill the PieceArray
	b.PieceArray = append(b.PieceArray[:0], b.PieceArray[1:]...) //remove a piece of PieceArray
	if len(b.PieceArray) <= 0 {
		b.FillArray()
	}
	b.nextPiece = b.PieceArray[0]
	b.drawNextPiece(b.nextPiece.color)

	return ""
}

func (b *Board) drawPiece(p Piece, c BlockColor) {
	for i := 0; i < 4; i++ {
		b.gameboard[p.Piece_[i].row][p.Piece_[i].col] = c
	}
}
