package agent

import (
	"github.com/TetrisAI/project/gameboard"
	"github.com/sarsa"
	"strconv"
	//"strings"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"hash/fnv"
	"math/rand"
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
	s.board.ResetBoard() //board
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

	cad := ""
	for i := 0; i < len(pos); i++ {
		cad += strconv.Itoa(pos[i])
	}

	return cad
}

func getInfoActivePiece(p gameboard.Piece) string {
	return strconv.Itoa(p.GetId())
}

func getBoardInfo(board gameboard.Board) string {
	Ys := make([]int, 0)

	for i := 0; i < len(board.GetGameBoard()); i++ {
		for j := 0; j < len(board.GetGameBoard()[i]); j++ {
			if board.GetGameBoard()[i][j] != 0 {
				Ys = append(Ys, 1+i+j)
			}
		}
	}

	cad := ""
	for i := 0; i < len(Ys); i++ {
		cad += strconv.Itoa(Ys[i])
	}
	return cad
}

func (s State) GetActiveTiles(action string) [][]int {

	activeTiles := make([][]int, s.v.Features)

	for idx := range activeTiles {
		activeTiles[idx] = make([]int, 0)
	}
	ap := 1
	for tile := 0; tile < s.v.Tilings; tile++ {
		key := strconv.Itoa(tile*ap) + action
		ap *= s.v.Tilings
		activeTiles[0] = append(activeTiles[0], s.Idx(getYs(s.board)+key+getPosActivePiece(s.board.GetActivePiece())))
		activeTiles[1] = append(activeTiles[1], s.Idx(getBoardInfo(s.board)+key))
		//activeTiles[1] = append(activeTiles[1], s.Idx(getPosActivePiece(s.board.GetActivePiece())+key))
		activeTiles[2] = append(activeTiles[2], s.Idx(getInfoActivePiece(s.board.GetActivePiece())+key+getYs(s.board)))

	}
	//}
	//fmt.Println("MIS ACTIVE TILES PENDEJO: ", activeTiles)

	return activeTiles
	/*FEATURES
	-getY: 72357635+key
	getPosActivePiece: [x1,y1,...,x4,y4]+ key


	º*/

}

func getYs(board gameboard.Board) string {
	Ys := make([]int, 0)

	for i := 0; i < len(board.GetGameBoard()[0]); i++ {
		j := 0
		//fmt.Println("muero aqui")
		for j < 20 && board.GetGameBoard()[j][i] != 0 {
			j++
		}
		//fmt.Println("asi es")
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
	//fmt.Println("h")
	/*for {

		//fmt.Println("Velocidad")

		if Time >= 9 {
			Time = 0
			//if Time == 0 {
			//fmt.Println("Gravity")
			//s.board.Gravity()
			//}

		}

		//time.Sleep(25 * time.Millisecond)
		Time++
	}*/
}

func (s State) InGoalState() bool {
	return s.board.Game_over
}

func (s *State) Start(option int) int {
	steps := 0

	switch option {
	case 0:
		steps = sarsa.SemiGradientSarsa(s, valueOf, getAction, s.v)
	case 1:
		steps = sarsa.SemiGradientSarsa(s, valueOf, getEAction, s.v)
	case 2:
		steps = sarsa.SemiGradientSarsa(s, valueOf, GetActionPlayer, s.v)
	}

	return steps

}

//***********

func valueOf(state sarsa.State, action string, vf *sarsa.ValueFunction) float64 {
	if state.InGoalState() {
		return 0.0
	}

	activeTiles := state.GetActiveTiles(action)
	estimations := make([]float64, vf.Features)

	for feature := 0; feature < vf.Features; feature++ {
		for idx := 0; idx < vf.Tilings; idx++ {
			estimations[feature] += vf.Weights[activeTiles[feature][idx]]
		}
	}
	val := 0.0
	for estimation := range estimations {
		val += estimations[estimation]
	}

	return val
}

func getEAction(state sarsa.State, vf *sarsa.ValueFunction) string {
	if rand.Float64() > 0.7 {
		ac := state.GetActions()
		idx := rand.Intn(len(ac))
		return ac[idx]
	}
	return getAction(state, vf)
}

func getAction(state sarsa.State, vf *sarsa.ValueFunction) string {
	values := make([]float64, 0)
	actions := state.GetActions()
	for _, action := range actions {
		values = append(values, valueOf(state, action, vf))
	}
	//	fmt.Println("Actions: ", values)
	ac := actions[getIdxMax(values)]
	return ac
}

func (s State) GetWin() *pixelgl.Window {
	return s.board.GetWin()
}

func GetActionPlayer(state sarsa.State, vf *sarsa.ValueFunction) string {

	values := make([]float64, 0)
	actions := state.GetActions()
	for _, action := range actions {
		values = append(values, valueOf(state, action, vf))
	}
	fmt.Println("Actions: ", values)

	win := state.GetWin()
	if win.Pressed(pixelgl.KeyDown) {
		return "down"
	}
	if win.Pressed(pixelgl.KeyRight) {
		return "right"
	}
	if win.Pressed(pixelgl.KeyLeft) {
		return "left"
	}
	if win.Pressed(pixelgl.KeyUp) {
		return "up"
	}
	if win.Pressed(pixelgl.KeySpace) {
		return "space"
	}
	return "none"
}

func getIdxMax(slice []float64) int {
	idx := 0
	max := slice[idx]
	//get the idx of the biggest element
	for i := 1; i < len(slice); i++ {
		if max < slice[i] {
			idx = i
			max = slice[i]
		}
		//If max and slice are equal, we randomly change so have more exploration in the algorithm
		if max == slice[i] {
			if rand.Float64() <= 0.5 {
				idx = i
				max = slice[i]
			}
		}
	}
	return idx
}

//***********
var last_time = time.Now()
var last_action = "none"

func (s State) TakeAction(action string) (sarsa.State, float64) {
	Time++
	if Time >= 9 {
		//fmt.Println("Gravity")
		s.board.Gravity()
		Time = 0
	}
	if time.Since(last_time).Seconds() < 0.3 && last_action == action {
		return s, 0
	}
	last_time = time.Now()
	last_action = action
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

	//fmt.Println("Time: ", Time)

	score := s.board.Score()
	s.board.SetScore(0)

	return s, float64(score)
}

func (s *State) SetBoard(board gameboard.Board) {
	s.board = board
}
func NewState() State {
	s := State{}
	s.board.FillArray()
	s.v = &sarsa.ValueFunction{}
	s.v.New(3, 10000000, 16, 0.4/8)
	s.hash_table = make(map[string]int, 10000000)
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
