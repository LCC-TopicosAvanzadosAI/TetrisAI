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
	piece [4]Block
	color BlockColor
	id    IdPiece
}

//I'm still thinking about this, we could use a binary representation but that's really ugly and it will be really difficult to update in the future
//This way is really easy to understand but the code there is really ugly. So I don't know. I'll let my team decide...

func JPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 0, col: 1},
			Block{row: 0, col: 0},
			Block{row: 0, col: 2},
		},
		color: c,
		id:    IdJPiece,
	}
}
func LPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
			Block{row: 0, col: 0},
		},
		color: c,
		id:    IdLPiece,
	}
}
func SPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 0, col: 0},
			Block{row: 0, col: 1},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
		},
		color: c,
		id:    IdSPiece,
	}
}
func IPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
			Block{row: 1, col: 3},
		},
		color: c,
		id:    IdIPiece,
	}
}
func ZPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 0, col: 1},
			Block{row: 0, col: 2},
		},
		color: c,
		id:    IdZPiece,
	}
}
func TPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 1, col: 2},
			Block{row: 0, col: 1},
		},
		color: c,
		id:    IdTPiece,
	}
}
func OPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 1, col: 0},
			Block{row: 1, col: 1},
			Block{row: 0, col: 0},
			Block{row: 0, col: 1},
		},
		color: c,
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

type Board struct {
	gameboard   [BoardRows][BoardCols]BlockColor
	activePiece Piece
	activeShape int
	nextPiece   Piece
	gameOver    bool
	PieceArray  []Piece
	score       int
	Atlas       *text.Atlas
}

func NewGameBoard() Board {

	face, err := hp.LoadTTF("./../../resources/saarland.ttf", 40) //Loading font and size-font
	if err != nil {
		panic(err)
	}

	board := Board{
		gameOver: false,
		score:    0,
		Atlas:    text.NewAtlas(face, text.ASCII), //Atlas necessary for the font
	}

	board.FillArray()

	return board
}

func (b *Board) FillArray() {
	b.PieceArray = append(b.PieceArray, JPiece(1), LPiece(2), SPiece(3),
		TPiece(4), ZPiece(5), IPiece(6), OPiece(7))

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
		p.piece[i].row += int(rows)
		p.piece[i].col += int(cols)
	}
}

func (b *Board) checkCollision(p Piece, _r, _c Movement) bool {
	for i := 0; i < 4; i++ {
		r := p.piece[i].row + int(_r)
		c := p.piece[i].col + int(_c)
		if r < 0 || r > 21 || c < 0 || c > 9 || b.gameboard[r][c] != Empty {
			return true
		}
	}
	return false
}

func (p *Piece) Rotate() {

	pivot := p.piece[1]
	arr := []int{0, 2, 3}
	for _, i := range arr {
		dRow := pivot.row - p.piece[i].row
		dCol := pivot.col - p.piece[i].col
		p.piece[i].row = pivot.row + (dCol * -1)
		p.piece[i].col = pivot.col + (dRow)
	}
}

func (b *Board) RotatePiece() {

	if b.activePiece.id == IdOPiece {
		return
	}

	b.drawPiece(b.activePiece, Empty)
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
			return
		}
	}

	b.drawPiece(b.activePiece, b.activePiece.color)

}

func (b *Board) MovePiece(dir Movement) {
	b.drawPiece(b.activePiece, Empty)

	if dir == MoveToBottom {
		if !b.checkCollision(b.activePiece, MoveDown, nothing) {
			b.activePiece.movePiece(MoveDown, nothing)
		} else {
			b.Gravity()
		}
	} else if !b.checkCollision(b.activePiece, nothing, dir) {
		b.activePiece.movePiece(nothing, dir)
	}

	b.drawPiece(b.activePiece, b.activePiece.color)
}

func (b *Board) MoveToBottom1() {
	b.drawPiece(b.activePiece, Empty)
	for !b.checkCollision(b.activePiece, MoveDown, nothing) {
		b.activePiece.movePiece(MoveDown, nothing)
	}
	b.drawPiece(b.activePiece, b.activePiece.color)
}

func (b *Board) Gravity() bool {
	//We remove the piece that we're trying to move down
	b.drawPiece(b.activePiece, Empty)
	//If there are no collisions we move the piece down.
	if !b.checkCollision(b.activePiece, MoveDown, nothing) {
		b.activePiece.movePiece(MoveDown, nothing)
	} else { //If there are collision we don't move down the piece and we add a new piece
		b.drawPiece(b.activePiece, b.activePiece.color)
		b.checkRowCompletion()
		b.AddPiece()

		return true
	}
	//Draw the piece that we remove.
	b.drawPiece(b.activePiece, b.activePiece.color)
	return false
}

func (b *Board) checkRowCompletion() {
	rowWasDeleted := true
	var linesFound int
	var deleteRowCt int
	linesFound = 0
	for rowWasDeleted {
		rowWasDeleted = false

		for i := 0; i < 4; i++ {
			r := b.activePiece.piece[i].row
			emptyFound := false
			for c := 0; c < 10; c++ {
				if b.gameboard[r][c] == Empty {
					emptyFound = true
					continue
				}
			}
			if !emptyFound {
				b.deleteRow(r)
				rowWasDeleted = true
				deleteRowCt++
				linesFound += 1
				b.score += 100
			}
		}
		fmt.Println(linesFound)
		if linesFound == 4 {
			b.score += 800
			linesFound = 0
		}
	}

}

func (b *Board) deleteRow(row int) {
	for r := row; r < 21; r++ {
		for c := 0; c < 10; c++ {
			b.gameboard[r][c] = b.gameboard[r+1][c]
		}
	}
}

func (b *Board) AddPiece() {
	/*if p.Id == IdIPiece {
		offset = rand.Intn(7)
	} else if p.Id == OPiece {
		offset = rand.Intn(9)
	} else {*/
	//}
	//baseShape := getShapeFromPiece(nextPiece)
	//baseShape = moveShape(20, offset, baseShape)
	b.nextPiece.movePiece(Movement(18), Movement(4))
	b.activePiece = b.nextPiece
	b.drawPiece(b.activePiece, b.activePiece.color)

	//remove a piece or refill the PieceArray
	b.PieceArray = append(b.PieceArray[:0], b.PieceArray[1:]...) //remove a piece of PieceArray
	if len(b.PieceArray) <= 0 {
		b.FillArray()
	}
	b.nextPiece = b.PieceArray[0]
}

func (b *Board) DisplayBoard(win *pixelgl.Window, blockGen func(int) pixel.Picture) {
	boardBlockSize := 20.0 //win.Bounds().Max.X / 10
	pic := blockGen(0)
	imgSize := pic.Bounds().Max.X
	scaleFactor := float64(boardBlockSize) / float64(imgSize)

	//description for the score
	ScoreTxt := text.New(pixel.V(260, 307), b.Atlas) //here, I put the coordinates where the
	ScoreTxt.Dot.X -= ScoreTxt.BoundsOf(strconv.Itoa(b.score)).W()
	ScoreTxt.Color = colornames.Lightcyan //text color
	fmt.Fprintln(ScoreTxt, b.score)
	ScoreTxt.Draw(win, pixel.IM)

	//display next piece
	//b.drawPiece(b.nextPiece, b.nextPiece.color)

	for col := 0; col < BoardCols; col++ {
		for row := 0; row < BoardRows-2; row++ {
			val := b.gameboard[row][col]
			if val == Empty {
				continue
			}

			x := float64(col)*boardBlockSize + boardBlockSize/2
			y := float64(row)*boardBlockSize + boardBlockSize/2
			pic := blockGen(int(Gray) - 1)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor).Moved(pixel.V(x+282, y+25)))
		}
	}

	/*// Display Shadow
	pieceType := activePiece.color
	ghostShape := activeShape
	b.drawPiece(activeShape, Empty)
	for {
		if b.checkCollision(moveShapeDown(ghostShape)) {
			break
		}
		ghostShape = moveShapeDown(ghostShape)
	}
	b.drawPiece(activeShape, pieceType)

	gpic := blockGen(block2spriteIdx(Gray))
	sprite := pixel.NewSprite(gpic, gpic.Bounds())
	for i := 0; i < 4; i++ {
		if b[ghostShape[i].row][ghostShape[i].col] == Empty {
			x := float64(ghostShape[i].col)*boardBlockSize + boardBlockSize/2
			y := float64(ghostShape[i].row)*boardBlockSize + boardBlockSize/2
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scaleFactor/2).Moved(pixel.V(x+282, y+25)))
		}
	}
	*/
}

func (b *Board) drawPiece(p Piece, c BlockColor) {
	//fmt.Println("c:: ", c)
	for i := 0; i < 4; i++ {
		//	b[activeShape[i].row][activeShape[i].col] = t
		//b[p.piece[p.position][i].row][p.piece[p.position][i].col] = p.color
		b.gameboard[p.piece[i].row][p.piece[i].col] = c
	}
}
