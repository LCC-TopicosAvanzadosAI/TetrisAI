package agent

import (
	"github.com/TetrisAI/project/gameboard"
	"github.com/sarsa"
	"strconv"
	//"strings"
	//"fmt"
	"hash/fnv"
	"time"
)

var Time = 0

type State struct {
	board      gameboard.Board
	v          *sarsa.ValueFunction
	hash_table map[string]int
	max_size   int
}

func (s State) GetValueFunction() *sarsa.ValueFunction {
	return s.v
}
func (s State) GetRandomFirstPosition() sarsa.State {
	//board := gameboard.Board{}
	//board.FillArray()
	//board.Draw = true

	//s.board = gameboard.NewGameBoard()
	s.board.SetScore(0)
	Time = 0
	return s
}

func getPosActivePiece(p gameboard.Piece) string {
	pos := make([]int, 0)

	for i := range p.Piece_ {
		pos = append(pos, p.Piece_[i].GetRow())
		pos = append(pos, p.Piece_[i].GetCol())
	}

	//return strings.Join(pos[:], ",")

	cad := ""
	for i := 0; i < len(pos); i++ {
		cad += strconv.Itoa(pos[i])
	}

	return cad
}

func getInfoActivePiece(p gameboard.Piece) string {
	return strconv.Itoa(p.GetId())
}

func (s State) GetActiveTiles(action string) [][]int {

	activeTiles := make([][]int, s.v.Features)

	for idx := range activeTiles {
		activeTiles[idx] = make([]int, 0)
	}
	//fmt.Println("FEATURES: ", s.v.Tilings)
	//for feature := 0; feature < s.v.Features; feature++ {
	for tile := 0; tile < s.v.Tilings; tile++ {
		key := strconv.Itoa(tile) + action
		activeTiles[0] = append(activeTiles[0], s.Idx(getYs(s.board)+key))
		activeTiles[1] = append(activeTiles[1], s.Idx(getPosActivePiece(s.board.GetActivePiece())+key))
		activeTiles[2] = append(activeTiles[2], s.Idx(getInfoActivePiece(s.board.GetActivePiece())+key))

	}
	//}
	//fmt.Println("MIS ACTIVE TILES PENDEJO: ", activeTiles)
	return activeTiles
	/*FEATURES
	-getY: 72357635+key
	getPosActivePiece: [x1,y1,...,x4,y4]+ key


	º*/

	/*//h := getYs(s.Gameboard)
	//pos := getPosActivePiece(s.board.ActivePiece)
	tiles := make([]int, 0)

	for tile := 0; tile < s.v.Tilings; tile++ {

		st := strconv.Itoa(tile) + action
		tiles = append(tiles, s.Idx(action))
		tiles = append(tiles, s.Idx(getInfoBoard()+st))
		tiles = append(tiles, s.Idx(getYs(s.Gameboard)+st))
		tiles = append(tiles, s.Idx(getInfoActivePiece(s.board.ActivePiece)+st))
		//tiles = append(tiles, s.Idx(action+h+pos+strconv.Itoa(tile)))
	}

	return tiles*/
}

func getYs(board gameboard.Board) string {
	Ys := make([]int, 0)

	for i := 0; i < len(board.GetGameBoard()[0]); i++ {
		j := 0
		for board.GetGameBoard()[i][j] != 0 {
			j++
		}
		Ys = append(Ys, j)
	}

	cad := ""
	for i := 0; i < len(Ys); i++ {
		cad += strconv.Itoa(Ys[i])
	}

	return cad
}

func (s State) GetActions() []string {
	actions := make([]string, 0)
	actions = append(actions, "up")
	actions = append(actions, "down")
	actions = append(actions, "left")
	actions = append(actions, "right")
	actions = append(actions, "space")
	return actions
}

func (s *State) StartTime() {
	for {
		time.Sleep(1 * time.Millisecond)
		Time++
		if Time == 9 {
			Time = 0
		}
	}
}

func (s State) InGoalState() bool {
	return false
}

func (s *State) Start() int {
	steps := sarsa.SemiGradientSarsa(s, s.v)
	return steps
}

func (s State) TakeAction(action string) (sarsa.State, float64) {
	switch action {
	case "down":
		//fmt.Println("down")
		s.board.Gravity()
	case "right":
		//fmt.Println("right")
		s.board.MovePiece(gameboard.MoveRight)
	case "left":
		//fmt.Println("left")
		s.board.MovePiece(gameboard.MoveLeft)
	case "space":
		//fmt.Println("space")
		s.board.MoveToBottom1()
	case "up":
		//fmt.Println("up")
		s.board.RotatePiece()
	}

	if Time == 0 {
		//fmt.Println("Gravity")
		s.board.Gravity()
	}
	score := s.board.Score()
	s.board.SetScore(0)
	//fmt.Println("Score: ", score)
	//fmt.Println("Nuevo Score: ", s.board.Score())
	//if s.board.Score() != 0 {
	//fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	//}

	return s, float64(score) - 1
}

func (s *State) SetBoard(board gameboard.Board) {
	s.board = board
}
func NewState() State {
	s := State{}
	s.board.FillArray()
	s.v = &sarsa.ValueFunction{}
	s.v.New(3, 1000000, 8, 0.4/8)
	s.hash_table = make(map[string]int, 1000000)
	return s
}

func (s *State) Idx(key string) int {
	idx, ok := s.hash_table[key]
	//if the element is in the hast table the idx is returned
	if ok {
		return idx
	}

	//overflow control
	if len(s.hash_table) >= len(s.v.Weights) {
		return hash(key) % len(s.hash_table)
	}
	//if the elemen is not in the hast table, the element is added.
	s.hash_table[key] = len(s.hash_table)
	return len(s.hash_table) - 1
}

//hash function
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}
