package gameboard

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
	_ "time"
)

const (
	BoardRows = 22
	BoardCols = 10
)

type Movement int

const (
	moveDown  Movement = -1
	moveLeft           = -1
	moveRight          = 1
	nothing            = 0
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
			Block{row: 0, col: 0},
			Block{row: 0, col: 1},
			Block{row: 0, col: 2},
			Block{row: 1, col: 2},
		},
		color: c,
		id:    IdJPiece,
	}
}
func LPiece(c BlockColor) Piece {
	return Piece{
		piece: [4]Block{
			Block{row: 0, col: 0},
			Block{row: 0, col: 1},
			Block{row: 0, col: 2},
			Block{row: 1, col: 0},
		},
		color: c,
		id:    IdLPiece,
	}
}

func randomPiece() Piece {
	p := rand.Intn(2)
	if p == 0 {
		return LPiece(1)
	} else {
		return JPiece(1)
	}
}

//TO DO
//s,t,z,i,o

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
}

//----

func NewGameBoard() Board {
	return Board{
		nextPiece: randomPiece(),
		gameOver:  false,
	}
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

func (b *Board) Gravity() {
	//blockType := b[activeShape[0].row][activeShape[0].col]
	// Erase old piece

	// Does the block collide if it moves down?
	//didCollide := b.checkCollision()
	fmt.Println(b.activePiece)
	b.drawPiece(b.activePiece, Empty)
	if !b.checkCollision(b.activePiece, moveDown, nothing) {
		b.activePiece.movePiece(moveDown, nothing)
	} else {
		b.AddPiece()
	}
	b.drawPiece(b.activePiece, b.activePiece.color)
	/*if didCollide {
		b.checkRowCompletion(activeShape)
		b.addPiece() // Replace with random piece
		return true
	}*/
	//return false
}

func (b *Board) AddPiece() {
	var offset int
	/*if p.Id == IdIPiece {
		offset = rand.Intn(7)
	} else if p.Id == OPiece {
		offset = rand.Intn(9)
	} else {*/
	offset = rand.Intn(8)
	//}
	//baseShape := getShapeFromPiece(nextPiece)
	//baseShape = moveShape(20, offset, baseShape)
	b.nextPiece.movePiece(Movement(16), Movement(offset))
	b.activePiece = b.nextPiece
	b.drawPiece(b.activePiece, b.activePiece.color)
	b.nextPiece = randomPiece()

}

func (b *Board) DisplayBoard(win *pixelgl.Window, blockGen func(int) pixel.Picture) {
	boardBlockSize := 20.0 //win.Bounds().Max.X / 10
	pic := blockGen(0)
	imgSize := pic.Bounds().Max.X
	scaleFactor := float64(boardBlockSize) / float64(imgSize)

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
	fmt.Println("c:: ", c)
	for i := 0; i < 4; i++ {
		//	b[activeShape[i].row][activeShape[i].col] = t
		//b[p.piece[p.position][i].row][p.piece[p.position][i].col] = p.color
		b.gameboard[p.piece[i].row][p.piece[i].col] = c
	}
}
