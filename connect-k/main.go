package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
)

type GameStore struct {
	Games map[string]*Game
	mu    sync.RWMutex
}

type Game struct {
	Id     string   `json:"id"`
	Mode   string   `json:"mode"`
	Winner string   `json:"winner"`
	Board  [][]bool `json:"board"`
	Cols   int      `json:"cols"`
	Rows   int      `json:"rows"`
	Win    int      `json:"win"`
	Level  int      `json:"level"`
}

// create a (cols x rows) game board. 'win' neighbored pieces are needed to win
// and mode determines the computer playmode (eg random)
func NewGame(cols, rows, win int, mode string, level int) *Game {
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
		Mode:   mode, // random, inteligent
		Winner: "",
		Board:  make([][]bool, cols),
		Cols:   cols,
		Rows:   rows,
		Win:    win,
		Level:  level, // only relevant when playing intelligent
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

const (
	COMPUTER = true
	PLAYER   = false
)

var gs = NewGameStore()

func main() {
	http.HandleFunc("/new", allowCors(newHandler))
	http.HandleFunc("/inspect", allowCors(inspectHandler))
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
	query := req.URL.Query()
	c := toIntOr(query.Get("c"), 7)
	r := toIntOr(query.Get("r"), 6)
	k := toIntOr(query.Get("k"), 4)
	l := toIntOr(query.Get("l"), 4)
	game := NewGame(c, r, k, "intelligent", l)
	gs.Set(game)
	w.Header().Set("Content-type", "application/json")
	enc := json.NewEncoder(w)
	fmt.Println("Created new game", game.Id, "level", game.Level, "board:", game.Cols, "x", game.Rows, "with", k, "wins.")
	enc.Encode(game)
}

func toIntOr(conv string, or int) int {
	v, err := strconv.Atoi(conv)
	if err != nil {
		return or
	}
	return v
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
	pRow, err := apply(game.Board, post.Col, false, game.Rows) // "player" is false, computer is true
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

func inspectHandler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	game := gs.Get(query.Get("id"))
	if game != nil {
		enc := json.NewEncoder(w)
		enc.Encode(&game)
	}
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
func apply(board [][]bool, col int, computer bool, maxRows int) (int, error) {
	//fmt.Println("Applying", col, computer, board)
	if col < 0 || col > len(board) {
		errMsg := fmt.Sprintf("Column out of board-boundaries %d", col)
		return -1, errors.New(errMsg)
	}
	if len(board[col]) >= maxRows {
		return -1, errors.New("Column already full")
	}
	board[col] = append(board[col], computer)
	return len(board[col]) - 1, nil
}

// The computer plays automatically and returns the column number and row where a chip was inserted
func autoPlay(game *Game) (int, int, error) {
	switch game.Mode {
	case "random":
		return playRandom(game)
	case "intelligent":
		return playIntelligent(game)
	}
	return -1, -1, errors.New("Unknown game mode")
}

func playRandom(game *Game) (int, int, error) {
	// has to terminate if all cols are full
	offset := randomNumber(game.Cols)
	for i := 0; i < game.Cols; i++ {
		c := (offset + i) % game.Cols
		r, err := apply(game.Board, c, true, game.Rows)
		if err == nil {
			return c, r, nil
		}
	}
	return -1, -1, errors.New("No more usable columns, all full!")
}

// provides a starting point to the alpha beta algorithm, using the computooor player
func playIntelligent(game *Game) (int, int, error) {

	alpha, beta := -1000, 1000
	_, col := alphaBeta(game.Level-1, alpha, beta, game.Board, true, game.Rows, game.Win)
	row, err := apply(game.Board, col, true, game.Rows)
	if err != nil { // column full
		panic("column full")
	}
	return col, row, nil
}

// Returns the best possible score and one random column with that score for
//'computer (true:comp, false:player)' when thinking 'depth' turns ahead.
// Aserts that 'computer' wants to insert a token at (col, row) into the board.
func scoreInDepth(depth, alpha, beta, col, row int, board [][]bool, computer bool, maxRows, win int) (int, int) {

	s := scoreInsertAt(col, row, computer, board, win, maxRows)
	// if I can win by inserting in c, no further eval of alpha/beta, this option is absolute at this point in time
	if depth == 0 || isWin(s, win, computer) {
		return s, -1 // no further enemy choices

	}
	s, c := alphaBeta(depth-1, alpha, beta, board, !computer, maxRows, win)
	//fmt.Println("When", computer, "inserted in", col, row, "result for AB at", depth-1, "for", !computer, "is", s, "and a b", alpha, beta)
	return s, c
}

// alpha is best score for computer so far in searchtree.
// beta is best score for player so far in searchtree.
//
// implements alpha beta pruning on negamax algorithm.
// returns two arguments: first the score and second a column at which to insert for best result.
func alphaBeta(depth, alpha, beta int, board [][]bool, computer bool, maxRows, win int) (int, int) {
	myCol := -1
	myScore := 1000
	if computer == COMPUTER {
		myScore = -1000
	}

	for _, c := range genShuff(len(board)) {
		bCopy := make([][]bool, len(board))
		copy(bCopy, board)
		r, err := apply(bCopy, c, computer, maxRows)
		if err != nil {
			continue // column full
		}
		s, _ := scoreInDepth(depth, alpha, beta, c, r, bCopy, computer, maxRows, win)

		switch computer {
		case COMPUTER:
			if s > myScore {
				myCol = c
			}
			myScore = max(myScore, s)
			alpha = max(alpha, myScore)
		case PLAYER:
			if s < myScore {
				myCol = c
			}
			myScore = min(myScore, s)
			beta = min(beta, myScore)
		}
		if beta <= alpha {
			break
		}
	}
	//fmt.Println(computer, "with col", myCol, "at depth", depth, "with myScore", myScore, "a b was", alpha, beta)
	return myScore, myCol
}

// checks if the insert at position (c, r) lead to a win for the player with that token
// sets the Winner-field on the game, if so.
func setWinner(c int, r int, game *Game) {
	computer := game.Board[c][r]
	if insertWins(c, r, computer, game) {
		if computer == COMPUTER {
			game.Winner = "computooor"
		} else {
			game.Winner = "player"
		}
	} else if boardFull(game) {
		game.Winner = "tie"
	}
}

func insertWins(c, r int, computer bool, game *Game) bool {
	score := scoreInsertAt(c, r, computer, game.Board, game.Win, game.Rows)
	return isWin(score, game.Win, computer)
}

func boardFull(game *Game) bool {
	if len(game.Board) != game.Cols {
		return false
	}
	for c := range game.Board {
		if len(game.Board[c]) != game.Rows {
			return false
		}
	}
	return true
}

// Determine if "score" wins the game, where k-fields are needed in a row to win.
// Winning is evaluated from the perspective of either "player (false)" or "computer (true)"
func isWin(score, k int, computer bool) bool {
	winScore := k * k * k * k
	if computer == PLAYER {
		winScore = -winScore
	}
	return score == winScore
}

// get max score of insert at pos (c, r), regarding all possible axis on the board
func scoreInsertAt(c int, r int, computer bool, board [][]bool, dist int, rows int) int {
	cr := scoreAxis(c, r, computer, board, dist, rows, "c-r", checkCrossRightUp, checkCrossLeftDown)
	cl := scoreAxis(c, r, computer, board, dist, rows, "c-l", checkCrossLeftUp, checkCrossRightDown)
	rl := scoreAxis(c, r, computer, board, dist, rows, "r-l", checkRightOf, checkLeftOf)
	ab := scoreAxis(c, r, computer, board, dist, rows, "a-b", checkAbove, checkBelow)

	done := make(chan struct{})
	defer close(done)
	res := merge(done, cr, cl, rl, ab)

	score := 0
	for ax := 0; ax <= 4; ax++ {
		r := <-res
		if computer == COMPUTER {
			score = max(score, r)
		} else {
			score = min(score, r)
		}
	}
	//fmt.Println("best score of all 4 axis for insert of", computer, "at", c, r, "is", score)
	return score
}

// Scores the axis, which is defined by the composition of the provided pointcheckers.
// (c, r) is used as center of the axis, scoring takes place for the given 'computer' asserting 'k' as
// necessary amount of similar tokens in a line for winning.
// 'rows' is needed for boundary checking
// Scoring function:
// the rate of on axis is either 0 if it is not usable
// or the maximum o.t. sums of all possible-neighbor scores over all lines of length k or k^4 if k tokens are found in a row.
func scoreAxis(c int, r int, computer bool, board [][]bool, k int, rows int, dir string, pc1 PointChecker, pc2 PointChecker) <-chan int {
	res := make(chan int)
	go func() {

		scoresPc1, afc1 := scoreNeighbors(c, r, computer, board, k, rows, pc1, dir)
		scoresPc2, afc2 := scoreNeighbors(c, r, computer, board, k, rows, pc2, dir)
		//fmt.Println("distances rated", c, r, dir, "pc1", scoresPc1)
		//fmt.Println("distances rated", c, r, dir, "pc2", scoresPc2)
		if afc1+afc2+1 >= k {
			// immediate win cause k in a row
			score := k * k * k * k
			if computer == PLAYER {
				score = -score
			}
			res <- score
			close(res)
			return
		}

		mid := k*k + k*(k-3)
		if computer == PLAYER {
			mid = -mid
		}
		axis := toAxis(scoresPc1, scoresPc2, mid)

		bestScore := 0
		for i := 0; i <= len(axis)-k; i++ {
			if computer == COMPUTER {
				bestScore = max(bestScore, sum(axis, i, i+k))
			} else {
				bestScore = min(bestScore, sum(axis, i, i+k)) // minimizer
			}
		}
		//fmt.Println("axis score", dir, "for", computer, "at", c, r, bestScore)
		res <- bestScore // axis usable, return rate
		close(res)
	}()
	return res
}

// takes the vals of the slice from start to stop and adds them upp.
func sum(arr []int, start, stop int) (sum int) {
	for ; start < stop; start++ {
		sum += arr[start]
	}
	return sum
}

// Scores the neighbors of point (c, r) on the board with bound 'maxRows' and k as necessary amuont to win
// The provided Pointchecker is used for checking the state of the neighbors up to a distance of 'k'.
// The Pointchecker itself guarantees staying in a correct angle.
// Returns the scores for the neighbors as map of dist(neighbor) -> score(neighbor) and the count of adjacent friendly neighbors as second argument.
// scoring function - scores the neighbors based on its property:
// UNR : 0
// EMPTY : k - dist
// FRIEND : k*(k-dist) + k*(k-3)
// OOB / FOE : terminate, do not expand result map any more
// The factor k-3 is needed for scoring arbitary k compositions - 3 is the minimal connect count that makes sense in a connect-k game
func scoreNeighbors(c, r int, computer bool, board [][]bool, k, maxRows int, pc PointChecker, dir string) (map[int]int, int) {
	scores := make(map[int]int)
	afc := 0 // adjacent friend count
	connected := true
OuterLoop:
	for dist := 1; dist < k; dist++ {
		score := 0
		switch pc(c, r, board, computer, dist, maxRows) {
		case OOB, FOE:
			break OuterLoop
		case UNR:
			connected = false
			break
		case EMPTY:
			score = k - dist
			connected = false
		case FRIEND:
			score = k*(k-dist) + k*(k-3)
			if connected {
				afc++
			}
		}
		if computer == PLAYER {
			score = -score
		}
		scores[dist] = score
	}
	//fmt.Println("scores for", dir, c, r, scores)
	return scores, afc
}

// result values of a pointcheck
const (
	OOB    = iota // out of bounds
	FRIEND = iota
	FOE    = iota
	EMPTY  = iota
	UNR    = iota // unreachable
)

// signature: col, row, board, computer, distance, available rows
// checks a point (dc, dr) which is 'd' distance away from the point (col, row)
// for one specific angle (cross/right etc).
// Returns a result integer, that indicates if the checked point is
// out of bounds, un/friendly, empty or unreachable
type PointChecker func(int, int, [][]bool, bool, int, int) int

func checkCrossRightUp(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c+d >= len(board) || r+d >= avail:
		return OOB
	case r+d > len(board[c+d]):
		return UNR
	case r+d == len(board[c+d]):
		return EMPTY
	case computer == board[c+d][r+d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkCrossRightDown(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c+d >= len(board) || r-d < 0:
		return OOB
	case r-d > len(board[c+d]):
		return UNR
	case r-d == len(board[c+d]):
		return EMPTY
	case computer == board[c+d][r-d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkCrossLeftDown(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c-d < 0 || r-d < 0:
		return OOB
	case r-d > len(board[c-d]):
		return UNR
	case r-d == len(board[c-d]):
		return EMPTY
	case computer == board[c-d][r-d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkCrossLeftUp(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c-d < 0 || r+d >= avail:
		return OOB
	case r+d > len(board[c-d]):
		return UNR
	case r+d == len(board[c-d]):
		return EMPTY
	case computer == board[c-d][r+d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkLeftOf(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c-d < 0 || r >= avail:
		return OOB
	case r > len(board[c-d]):
		return UNR
	case r == len(board[c-d]):
		return EMPTY
	case computer == board[c-d][r]:
		return FRIEND
	default:
		return FOE
	}
}

func checkRightOf(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c+d >= len(board) || r >= avail:
		return OOB
	case r > len(board[c+d]):
		return UNR
	case r == len(board[c+d]):
		return EMPTY
	case computer == board[c+d][r]:
		return FRIEND
	default:
		return FOE
	}
}

func checkAbove(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c >= len(board) || r+d >= avail:
		return OOB
	case r+d > len(board[c]):
		return UNR
	case r+d == len(board[c]):
		return EMPTY
	case computer == board[c][r+d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkBelow(c int, r int, board [][]bool, computer bool, d, avail int) int {
	switch {
	case c >= len(board) || r-d < 0:
		return OOB
	case r-d > len(board[c]):
		return UNR
	case r-d == len(board[c]):
		return EMPTY
	case computer == board[c][r-d]:
		return FRIEND
	default:
		return FOE
	}
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return // wg.Done()
			}
		}
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

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func randomNumber(max int) int {
	rBig, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(rBig.Int64())
}

// joins the two maps from one end to the other, uses 'mid' as their middle
func toAxis(map1, map2 map[int]int, mid int) []int {
	var result []int
	for i := len(map1); i >= 0; i-- { // rly want to access the 'len', not 'len-1'
		val, ok := map1[i]
		if ok {
			result = append(result, val)
		}
	}
	result = append(result, mid)
	for i := 0; i <= len(map2); i++ {
		val, ok := map2[i]
		if ok {
			result = append(result, val)
		}
	}
	return result
}

// generates a slice with the numbers from 0-max in it, randomly shuffled
func genShuff(max int) []int {
	a := make([]int, max)
	for i := 0; i < max; i++ {
		a[i] = i
	}
	for i := range a {
		j := randomNumber(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}
