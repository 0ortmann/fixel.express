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
		return playIntelligent(game, 6)
	}
	return -1, -1, errors.New("Unknown game mode")
}

func playRandom(game *Game) (int, int, error) {
	// has to terminate if all cols are full
	c := randomNumber(7)
	for i := 0; i < game.Rows; i++ {
		c = (c + i) % game.Rows
		r, err := apply(game.Board, c, "computooor")
		if err == nil {
			fmt.Println("")
			return c, r, nil
		}
	}
	return -1, -1, errors.New("No more usable columns, all full!")
}

// provides a starting point to the alpha beta algorithm, using the computooor player
func playIntelligent(game *Game, depth int) (int, int, error) {

	alpha, beta := -1000, 1000
	_, choices := goMaximizer(depth-1, alpha, beta, game.Board, "computooor", game.Rows, game.Win)

	col := choices.GetOne()
	row, err := apply(game.Board, col, "computooor")
	if err != nil { // column full
		panic("column full")
	}
	fmt.Println("")
	return col, row, nil
}

// alpha is best score for computer so far in seachtree
// beta is best score for player so far in seachtree
//
// returns the best possible score for 'name' when thinking 'depth' turns ahead
// asserts that 'name' wants to insert a token at (col, row) into the board
func alphaBeta(depth, alpha, beta, col, row int, board [][]string, name string, maxRows, win int) (int, *IntSet) {

	if depth == 0 {
		s := scoreInsertAt(col, row, name, board, win, maxRows)
		return s, NewIntSet()

	}
	// simulate move and aggreate via alpha beta again whats best
	if name == "computooor" { // maximizer
		return goMinimizer(depth-1, alpha, beta, board, "player", maxRows, win)
	}
	return goMaximizer(depth-1, alpha, beta, board, "computooor", maxRows, win)
}

func goMaximizer(depth, alpha, beta int, board [][]string, name string, maxRows, win int) (int, *IntSet) {
	// until all possible columns examined or alpha gte beta
	myScore, enemyScore := -1000, 1000
	myOptions, enemyOptions := make(map[int]*IntSet), make(map[int]*IntSet) // (score -> c)
	for c := 0; c < len(board); c++ {
		bCopy := make([][]string, len(board))
		copy(bCopy, board)
		r, err := apply(bCopy, c, name)
		if err != nil { // column full
			continue
		}
		s, choices := alphaBeta(depth, alpha, beta, c, r, bCopy, name, maxRows, win)
		myOptions[s] = myOptions[s].Add(c)
		enemyOptions[s] = enemyOptions[s].AddAll(choices)
		myScore = max(myScore, s)
		enemyScore = min(enemyScore, s)
		alpha = max(alpha, myScore)
		if beta < alpha {
			break
		}
	}
	fmt.Println(name, "with options", myOptions[myScore], "enemy options", enemyOptions[myScore], "at depth", depth, "with myScore", myScore, "a b was", alpha, beta)
	if myScore == enemyScore && enemyOptions[myScore].Length() > 1 {
		// this avoids traps: if the best I could do is best for my opponent, then I at steal one option from him
		myOptions[myScore] = enemyOptions[myScore]
	}
	fmt.Println("new options for", name, myOptions[myScore])
	return myScore, myOptions[myScore]
}

func goMinimizer(depth, alpha, beta int, board [][]string, name string, maxRows, win int) (int, *IntSet) {
	// until all possible columns examined or alpha gte beta
	myScore, enemyScore := 1000, -1000
	myOptions, enemyOptions := make(map[int]*IntSet), make(map[int]*IntSet) // (score -> c)

	for c := 0; c < len(board); c++ {
		bCopy := make([][]string, len(board))
		copy(bCopy, board)
		r, err := apply(bCopy, c, name)
		if err != nil { // column full
			continue
		}
		s, choices := alphaBeta(depth, alpha, beta, c, r, bCopy, name, maxRows, win)
		myOptions[s] = myOptions[s].Add(c)
		enemyOptions[s] = enemyOptions[s].AddAll(choices)
		myScore = min(myScore, s)
		beta = min(beta, myScore)
		enemyScore = max(enemyScore, s)
		if beta < alpha {
			break
		}

	}
	fmt.Println(name, "with options", myOptions[myScore], "enemy options", enemyOptions[myScore], "at depth", depth, "with myScore", myScore, "a b was", alpha, beta)
	if myScore == enemyScore && enemyOptions[myScore].Length() > 1 {
		// this avoids traps: if the best I could do is best for my opponent, then I at steal one option from him
		myOptions[myScore] = enemyOptions[myScore]
	}
	fmt.Println("new options for", name, myOptions[myScore])
	return myScore, myOptions[myScore]
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
	score := scoreInsertAt(c, r, name, game.Board, game.Win, game.Rows)
	if name == "player" {
		score = 0 - score
	}
	preciseK := float64(game.Win)
	winScore := ((preciseK+1)/2 + 1) * preciseK * preciseK // k*k + k * sum(1..k)
	if float64(score) >= winScore {
		return true
	}
	return false
}

// get max score of insert at pos (c, r), regarding all possible axis on the board
func scoreInsertAt(c int, r int, name string, board [][]string, dist int, rows int) int {
	cr := scoreAxis(c, r, name, board, dist, rows, "c-r", checkCrossRightUp, checkCrossLeftDown)
	cl := scoreAxis(c, r, name, board, dist, rows, "c-l", checkCrossLeftUp, checkCrossRightDown)
	rl := scoreAxis(c, r, name, board, dist, rows, "r-l", checkRightOf, checkLeftOf)
	ab := scoreAxis(c, r, name, board, dist, rows, "a-b", checkAbove, checkBelow)

	score := 0
	axis := 0
	for r := range merge(cr, cl, rl, ab) {
		axis++
		if name == "player" {
			score = min(score, r)
		} else {
			score = max(score, r)
		}
		if axis == 4 {
			break
		}
	}
	fmt.Println("best score of all 4 axis for insert of", name, "at", c, r, "is", score)
	return score
}

// Scores the axis, which is defined by the composition of the provided pointcheckers.
// (c, r) is used as center of the axis, scoring takes place for the given 'name' asserting 'k' as
// necessary amount of similar pieces in a line for winning.
// 'rows' is needed for boundary checking
// Scoring function:
// the rate of on axis is either -1 if it is not usable
// or the maximum o.t. sums of all possible-neighbor scores over all lines of length k
func scoreAxis(c int, r int, name string, board [][]string, k int, rows int, dir string, pc1 PointChecker, pc2 PointChecker) <-chan int {
	res := make(chan int)
	go func() {

		scoresPc1 := scoreNeighbors(c, r, name, board, k, rows, pc1, dir)
		scoresPc2 := scoreNeighbors(c, r, name, board, k, rows, pc2, dir)
		//fmt.Println("distances rated", c, r, dir, "pc1", scoresPc1)
		//fmt.Println("distances rated", c, r, dir, "pc2", scoresPc2)

		var mid int
		if name == "player" {
			mid = -(k*k + k)
		} else {
			mid = k*k + k
		}

		axis := toAxis(scoresPc1, scoresPc2, mid)

		bestScore := 0
		for i := 0; i <= len(axis)-k; i++ {
			if name == "player" {
				bestScore = min(bestScore, straightOrStop(axis, i, i+k, false)) // minimizer
			} else {
				bestScore = max(bestScore, straightOrStop(axis, i, i+k, true))
			}
		}
		//fmt.Println("axis score", dir, "for", name, "at", c, r, bestScore)
		res <- bestScore // axis usable, return rate
	}()
	return res
}

// takes the vals of the slice from start to stop and adds them upp.
// all numbers must / must not be "positive".
// As soon as a different number is found, terminate, return 0, else the sum
func straightOrStop(arr []int, start, stop int, positive bool) int {
	sum := 0
	for ; start < stop; start++ {
		val := arr[start]
		if (positive && val < 0) || (!positive && val > 0) {
			return 0
		}
		sum += val
	}
	return sum
}

// Scores the neighbors of point (c, r) on the board with bound 'maxRows' and k as necessary amuont to win
// The provided Pointchecker is used for checking the state of the neighbors up to a distance of 'k'.
// The Pointchecker itself guarantees staying in a correct angle.
// Returns the scores for the neighbors as map of dist(neighbor) -> score(neighbor)
// scoring function - scores the neighbors based on its property:
// UNR : 0
// EMPTY : 1
// FRIEND : (k-dist) * k
// FOE : (k-dist) * -k
// OOB : terminate, do not expand result map any more
func scoreNeighbors(c, r int, name string, board [][]string, k, maxRows int, pc PointChecker, dir string) map[int]int {
	scores := make(map[int]int)
	for dist := 1; dist < k; dist++ {
		switch pc(c, r, board, name, dist, maxRows) {
		case OOB:
			break
		case UNR:
			scores[dist] = 0
		case EMPTY:
			if name == "player" {
				scores[dist] = -1 * (k - dist)
			} else {
				scores[dist] = 1 * (k - dist)
			}
		case FRIEND:
			if name == "player" {
				scores[dist] = -k*(k-dist) - k
			} else {
				scores[dist] = k*(k-dist) + k
			}
		case FOE:
			if name == "player" {
				scores[dist] = k*(k-dist) + k
			} else {
				scores[dist] = -k*(k-dist) - k
			}
		}
	}
	//fmt.Println("scores for", dir, c, r, scores)
	return scores
}

// result values of a pointcheck
const (
	OOB    = iota // out of bounds
	FRIEND = iota
	FOE    = iota
	EMPTY  = iota
	UNR    = iota // unreachable
)

// signature: col, row, board, name, distance
// checks a point (dc, dr) which is 'd' distance away from the point (col, row)
// for one specific angle (cross/right etc).
// Returns a result integer, that indicates if the checked point is
// out of bounds, un/friendly, empty or unreachable
type PointChecker func(int, int, [][]string, string, int, int) int

func checkCrossRightUp(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c+d >= len(board) || r+d >= avail:
		return OOB
	case r+d > len(board[c+d]):
		return UNR
	case r+d == len(board[c+d]):
		return EMPTY // empty
	case name == board[c+d][r+d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkCrossRightDown(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c+d >= len(board) || r-d < 0:
		return OOB
	case r-d > len(board[c+d]):
		return UNR
	case r-d == len(board[c+d]):
		return EMPTY
	case name == board[c+d][r-d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkCrossLeftDown(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c-d < 0 || r-d < 0:
		return OOB
	case r-d > len(board[c-d]):
		return UNR
	case r-d == len(board[c-d]):
		return EMPTY
	case name == board[c-d][r-d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkCrossLeftUp(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c-d < 0 || r+d >= avail:
		return OOB
	case r+d > len(board[c-d]):
		return UNR
	case r+d == len(board[c-d]):
		return EMPTY
	case name == board[c-d][r+d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkLeftOf(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c-d < 0 || r >= avail:
		return OOB
	case r > len(board[c-d]):
		return UNR
	case r == len(board[c-d]):
		return EMPTY
	case name == board[c-d][r]:
		return FRIEND
	default:
		return FOE
	}
}

func checkRightOf(c int, r int, board [][]string, name string, d, avail int) int {
	//fmt.Println("check right of", c, r, d)
	switch {
	case c+d >= len(board) || r >= avail:
		return OOB
	case r > len(board[c+d]):
		return UNR
	case r == len(board[c+d]):
		return EMPTY
	case name == board[c+d][r]:
		return FRIEND
	default:
		return FOE
	}
}

func checkAbove(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c >= len(board) || r+d >= avail:
		return OOB
	case r+d > len(board[c]):
		return UNR
	case r+d == len(board[c]):
		return EMPTY
	case name == board[c][r+d]:
		return FRIEND
	default:
		return FOE
	}
}

func checkBelow(c int, r int, board [][]string, name string, d, avail int) int {
	switch {
	case c >= len(board) || r-d < 0:
		return OOB
	case r-d > len(board[c]):
		return UNR
	case r-d == len(board[c]):
		return EMPTY
	case name == board[c][r-d]:
		return FRIEND
	default:
		return FOE
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

// sum of the array from start to stop
func sum(arr []int, start, stop int) int {
	sum := 0
	for ; start < stop; start++ {
		sum += arr[start]
	}
	return sum
}

type IntSet struct {
	set map[int]bool
}

func NewIntSet() *IntSet {
	return &IntSet{
		set: make(map[int]bool),
	}
}

func (set *IntSet) Add(i int) *IntSet {
	if set == nil {
		set = NewIntSet()
	}
	set.set[i] = true
	return set
}

func (me *IntSet) AddAll(other *IntSet) *IntSet {
	if me == nil {
		me = NewIntSet()
	}
	for val := range me.set {
		other.Add(val)
	}
	return other
}

// returns a random element of the set. returns -1 if the set is empty.
func (set *IntSet) GetOne() int {
	for elem := range set.set {
		return elem // iteration not in order, so :)
	}
	return -1
}

func (set *IntSet) Length() int {
	return len(set.set)
}
