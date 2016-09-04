package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
)

type GameStore struct {
	Games map[string]*Game
	mu    sync.RWMutex
}

type Game struct {
	Id     string
	Mode   string
	Winner string
	Board  [][]string
	Cols   int
	Rows   int
	Win    int
}

// create a (cols x rows) game board. 'win' neighbored pieces are needed to win
// and mode determines the computer playmode (eg random)
func NewGame(cols, rows, win int, mode string) *Game {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Print("Error generating UUID", err)
		return nil
	}
	if cols < 1 || rows < 1 {
		log.Print("Please provide sane size", err)
		return nil
	}
	return &Game{
		Id:     fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]),
		Mode:   mode, // minmax, heuristics etc..
		Winner: "",
		Board:  make([][]string, cols),
		Cols:   cols,
		Rows:   rows,
		Win:    win,
	}
}

func NewGameStore() *GameStore {
	return &GameStore{
		Games: make(map[string]*Game),
	}
}

func (gs *GameStore) Get(id string) *Game {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.Games[id]
}

func (gs *GameStore) Set(game *Game) bool {
	if gs.Get(game.Id) != nil {
		return false
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Games[game.Id] = game
	return true
}

var gs = NewGameStore()

func main() {
	http.HandleFunc("/new", allowCors(newHandler))
	http.HandleFunc("/play", allowCors(errorHandler(playHandler)))

	http.ListenAndServe(":5000", nil)
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	// function adapter for errors
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		err := f(w, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error handling request for %q: %v", req.RequestURI, err)
		}
	}

}

func allowCors(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-type")
		if req.Method == "OPTIONS" {
			fmt.Fprintf(w, "OK")
			return
		}
		f(w, req)
	}
}

func newHandler(w http.ResponseWriter, req *http.Request) {
	game := NewGame(7, 6, 4, "intelligent")
	gs.Set(game)
	w.Header().Set("Content-type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(game)
}

func playHandler(w http.ResponseWriter, req *http.Request) error {
	type Post struct {
		GameId string
		Col    int
	}
	var post Post
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&post); err != nil {
		return err
	}
	game := gs.Get(post.GameId)
	if game == nil {
		return errors.New("No game with that ID")
	}
	if post.Col >= game.Cols {
		return errors.New("Wrong JSON format")
	}
	if game.Winner != "" {
		return errors.New("Game already ended, winner was " + game.Winner)
	}
	pRow, err := apply(game.Board, post.Col, "player")
	if err != nil {
		return err
	}
	setWinner(post.Col, pRow, game)
	if game.Winner != "" {
		return sendResult(w, game, -1)
	}
	cCol, cRow, err := autoPlay(game)
	if err != nil {
		return err
	}
	setWinner(cCol, cRow, game)
	return sendResult(w, game, cCol)
}

// fomulate a response with game winner (may be empty) and the column picked by the computer
func sendResult(w http.ResponseWriter, game *Game, cCol int) error {
	type Resp struct {
		Col    int    `json:"col"`
		Winner string `json:"winner"`
	}
	var resp Resp
	resp.Col = cCol
	resp.Winner = game.Winner
	enc := json.NewEncoder(w)
	return enc.Encode(&resp)
}

// put a token with value "player" into the column "col".
// Returns the row it was inserted into and an error if the column is already full
func apply(board [][]string, col int, player string) (int, error) {
	if len(board[col]) == 6 {
		return -1, errors.New("Column already full")
	}
	board[col] = append(board[col], player)
	return len(board[col]) - 1, nil
}

// The computer plays automatically and returns the column number and row where a chip was inserted
func autoPlay(game *Game) (int, int, error) {
	switch game.Mode {
	case "random":
		return playRandom(game)
	case "intelligent":
		return playIntelligent(game, 2)
	}
	return -1, -1, errors.New("Unknown game mode")
}

func playRandom(game *Game) (int, int, error) {
	// has to terminate if all cols are full
	rCol, _ := rand.Int(rand.Reader, big.NewInt(7))
	c := int(rCol.Int64())
	for i := 0; i < 6; i++ {
		c = (c + i) % 6
		r, err := apply(game.Board, c, "computooor")
		if err == nil {
			return c, r, nil
		}
	}
	return -1, -1, errors.New("No more usable columns, all full!")
}

// depth = 2
func playIntelligent(game *Game, depth int) (int, int, error) {
	maximums := make(map[int]int, game.Cols) // col -> val
	minimums := make(map[int]int, game.Cols) // cal -> val

	for c := 0; c < game.Cols; c++ {
		comBc := make([][]string, game.Cols)
		copy(comBc, game.Board)
		_, err := apply(comBc, c, "computooor")
		if err != nil {
			continue
		}
		// minimizer need to do this at least once
		for c2 := 0; c2 < game.Cols; c2++ {
			pBc := make([][]string, game.Cols)
			copy(pBc, comBc)
			r2, err := apply(pBc, c2, "player")
			if err != nil { // else column full
				continue
			}
			minimums[c2] = 0 - rateInsertAt(c2, r2, "player", pBc, game.Win, game.Rows)
		}
		pMin, pCol := min(minimums)
		fmt.Println("player rates for computer move", c, minimums, "would pick col", pCol)

		maximums[c] = pMin
	}
	_, cCol := max(maximums)
	fmt.Println("player will choose smallest of", maximums, "so i pick the highest. thats col", cCol)
	myRow, err := apply(game.Board, cCol, "computooor")
	if err != nil {
		return -1, -1, err
	}
	fmt.Println("")
	fmt.Println("")
	return cCol, myRow, nil
}

// checks if the insert at position (c, r) lead to a win for the player with that token
// sets the Winner-field on the game, if so.
func setWinner(c int, r int, game *Game) {
	name := game.Board[c][r]
	if insertWins(c, r, name, game) {
		game.Winner = name
	}
}

func insertWins(c, r int, name string, game *Game) bool {
	rate := rateInsertAt(c, r, name, game.Board, game.Win, game.Rows) + 1 // count self
	if rate >= game.Win {
		return true
	}
	return false
}

// get max rate of insert at pos (c, r), regarding all possible axis on the board
func rateInsertAt(c int, r int, name string, board [][]string, dist int, rows int) int {
	cr := rateAlongAxis(c, r, name, board, dist, rows, "c-r", checkCrossRightUp, checkCrossLeftDown)
	cl := rateAlongAxis(c, r, name, board, dist, rows, "c-l", checkCrossLeftUp, checkCrossRightDown)
	rl := rateAlongAxis(c, r, name, board, dist, rows, "r-l", checkRightOf, checkLeftOf)
	ab := rateAlongAxis(c, r, name, board, dist, rows, "a-b", checkAbove, checkBelow)

	rate := 0
	axis := 0
	for r := range merge(cr, cl, rl, ab) {
		axis++
		if rate < r {
			rate = r
		}
		if axis == 4 {
			break
		}
	}
	return rate
}

// rates the point (c, r) along an axis that is specified by the provided Pointcheckers.
// calculates the rate for that axis and pushes it into the result channel.
// If the axis is not usable because one cannot win with that axis, it has value -1.
// else the rate of an axis is defined as the max of all rates that are
// found in distance d = game.Win from (c, r)
func rateAlongAxis(c int, r int, name string, board [][]string, dist int, rows int, dir string, pc1 PointChecker, pc2 PointChecker) <-chan int {
	res := make(chan int)
	go func() {

		ratesPc1, maxD1 := ratePoint(c, r, name, board, dist, rows, pc1)
		ratesPc2, maxD2 := ratePoint(c, r, name, board, dist, rows, pc2)
		//fmt.Println("distances rated", c, r, dir, "pc1", ratesPc1, maxD1)
		//fmt.Println("distances rated", c, r, dir, "pc2", ratesPc2, maxD2)

		if maxD1+maxD2-1 < dist {
			// cannot form k adjacent chips for name, axis is useless
			//fmt.Println("maxrate for", c, r, dir, -1)
			res <- -1
			return
		}

		maxRate := 0 // holds the maximal rate that is possible within a game winning dist
		for i := 0; i < dist; i++ {
			// build possible distances
			rate1, ok1 := ratesPc1[i]
			rate2, ok2 := ratesPc2[dist-i-1]
			if !ok1 && ok2 && maxRate < rate2 {
				maxRate = rate2
			}
			if ok1 && !ok2 && maxRate < rate1 {
				maxRate = rate1
			}
			if ok1 && ok2 && maxRate < rate1+rate2 {
				maxRate = rate1 + rate2
			}
		}
		//fmt.Println("maxrate for", c, r, dir, maxRate)
		res <- maxRate // axis usable, return rate
	}()
	return res
}

// rates the point (c, r) in on the 'board' with bounds 'maxDist' and 'maxRows' and uses the given pointchecker.
// returns the rate for each distance of fields that were considered for rating as a map,
// eg. at dist 1: rate 0, at dist 2 : rate 1 etc..
// returns the max dist of considered points as second argument.
// rating function:
// introspect k-th neighbor, using the Pointchecker (guarantees angle):
// found same chip as self: rate at that point is 1 + rates of previos point
// found empty field: rate for that point is rate of previos point
// found field out of bounds or unfriendly chip: rate at that point is undefined, terminate checking
func ratePoint(c, r int, name string, board [][]string, maxDist int, maxRows int, pc PointChecker) (map[int]int, int) {
	rates := make(map[int]int)
	dist := 1
StraightInARow:
	for ; dist < maxDist; dist++ {
		switch pc(c, r, board, name, dist, maxRows) {
		case 1: // same chip
			rates[dist] = rates[dist-1] + 1
		case 0: // no chip, not out of bounds
			rates[dist] = rates[dist-1]
		case -1, -2:
			break StraightInARow
		}
	}
	return rates, dist
}

// signature: col, row, board, name, distance
// checks a point (dc, dr) which is 'd' distance away from the point (col, row)
// for one specific angle (cross/right etc).
// Returns a 1 if the given name equals the name on (dc, dr),
// a 0 if (dc, dr) is empty, a -1 if (dc, dr) out of bounds
// and a -2 if there is a different name found on that found.
// check if 'name' is found on a point on the board with distance d from point (c, r)
type PointChecker func(int, int, [][]string, string, int, int) int

func checkCrossRightUp(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c+d >= len(board) || r+d >= avail:
		return -1 // out of bounds
	case r+d >= len(board[c+d]):
		return 0 // empty
	case name == board[c+d][r+d]:
		return 1
	default:
		return -2
	}
}

func checkCrossRightDown(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c+d >= len(board) || r-d < 0:
		return -1
	case r-d >= len(board[c+d]):
		return 0
	case name == board[c+d][r-d]:
		return 1
	default:
		return -2
	}
}

func checkCrossLeftDown(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c-d < 0 || r-d < 0:
		return -1
	case r-d >= len(board[c-d]):
		return 0
	case name == board[c-d][r-d]:
		return 1
	default:
		return -2
	}
}

func checkCrossLeftUp(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c-d < 0 || r+d >= avail:
		return -1
	case r+d >= len(board[c-d]):
		return 0
	case name == board[c-d][r+d]:
		return 1
	default:
		return -2
	}
}

func checkLeftOf(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c-d < 0 || r >= avail:
		return -1
	case r >= len(board[c-d]):
		return 0
	case name == board[c-d][r]:
		return 1
	default:
		return -2
	}
}

func checkRightOf(c int, r int, board [][]string, name string, d, avail int) int {
	//fmt.Println("check right of", c, r, d)
	switch {
	case c+d >= len(board) || r >= avail:
		return -1
	case r >= len(board[c+d]):
		return 0
	case name == board[c+d][r]:
		return 1
	default:
		return -2
	}
}

func checkAbove(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c >= len(board) || r+d >= avail:
		return -1
	case r+d >= len(board[c]):
		return 0
	case name == board[c][r+d]:
		return 1
	default:
		return -2
	}
}

func checkBelow(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c >= len(board) || r-d < 0:
		return -1
	case name == board[c][r-d]:
		return 1
	default:
		return -2
	}
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// returns the max value of passed map and its key. if multiple vals are highest, the first found
// key is returned for that val. only considers vals that have a key
func max(nums map[int]int) (int, int) {
	if len(nums) == 0 {
		return 0, 0
	}
	max := -1000
	idx := -1
	for i, num := range nums {
		if max < num {
			max = num
			idx = i
		}
	}
	return max, idx
}

// returns the min value of passed map and its key. if multiple vals are lowest, the first found
// key is returned for that val. only considers vals that have a key
func min(nums map[int]int) (int, int) {
	if len(nums) == 0 {
		return 0, 0
	}
	min := 1000
	idx := -1
	for i, num := range nums {
		if min > num {
			min = num
			idx = i
		}
	}
	return min, idx
}
