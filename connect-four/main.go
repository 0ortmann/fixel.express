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
}

func NewGame() *Game {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Print("Error generating UUID ", err)
		return nil
	}
	return &Game{
		Id:     fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]),
		Mode:   "random", // minmax, heuristics etc..
		Winner: "",
		Board:  make([][]string, 7),
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
	game := NewGame()
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
	if post.GameId == "" || post.Col >= 7 {
		return errors.New("Wrong JSON format")
	}
	game := gs.Get(post.GameId)
	if game == nil {
		return errors.New("No game with that ID")
	}
	if game.Winner != "" {
		return errors.New("Game already ended, winner was " + game.Winner)
	}
	pRow, err := apply(game, post.Col, "player")
	if err != nil {
		return err
	}
	setWinner(post.Col, pRow, game, 4)
	if game.Winner != "" {
		return playResult(w, game, -1)
	}
	cCol, cRow, err := autoPlay(game)
	if err != nil {
		return err
	}
	setWinner(cCol, cRow, game, 4)
	return playResult(w, game, cCol)
}

// fomulate a response with game winner (may be empty) and the column picked by the computer
func playResult(w http.ResponseWriter, game *Game, cCol int) error {
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
func apply(game *Game, col int, player string) (int, error) {
	if len(game.Board[col]) == 6 {
		return -1, errors.New("Column already full")
	}
	game.Board[col] = append(game.Board[col], player)
	return len(game.Board[col]) - 1, nil
}

// The computer plays automatically and returns the column number and row where a chip was inserted
func autoPlay(game *Game) (int, int, error) {
	switch game.Mode {
	case "random":
		return playRandom(game)
	}
	return -1, -1, errors.New("Unknown game mode")
}

func playRandom(game *Game) (int, int, error) {
	// has to terminate if all cols are full
	rCol, _ := rand.Int(rand.Reader, big.NewInt(7))
	c := int(rCol.Int64())
	for i := 0; i < 6; i++ {
		c = (c + i) % 6
		r, err := apply(game, c, "computooor")
		if err == nil {
			return c, r, nil
		}
	}
	return -1, -1, errors.New("No more usable columns, all full!")
}

func setWinner(c int, r int, game *Game, k int) {
	name := game.Board[c][r]
	if insertWins(c, r, game.Board, name, k) {
		game.Winner = name
	}
}

// checks if an insert of the string 'name' at position (c, r)
// has 'k' direct neighbors with the same 'name' on the board
func insertWins(c int, r int, board [][]string, name string, k int) bool {
	cr := checkAxisWith(c, r, board, name, k, checkCrossRightUp, checkCrossLeftDown)
	cl := checkAxisWith(c, r, board, name, k, checkCrossLeftUp, checkCrossRightDown)
	rl := checkAxisWith(c, r, board, name, k, checkRightOf, checkLeftOf)
	blw := checkAxisWith(c, r, board, name, k, checkBelow)
	succ := false
	axis := 0
	for {
		select {
		case res := <-cr:
			succ = succ || res
			axis++
		case res := <-cl:
			succ = succ || res
			axis++
		case res := <-rl:
			succ = succ || res
			axis++
		case res := <-blw:
			succ = succ || res
			axis++
		}
		if succ || axis == 4 {
			return succ
		}
	}
}

// Invokes all the given Pointcheckers with the other arguments.
// Counts the in-order successes of point checking and stops counting on a Pointchecker once a negative success was detected.
// If the total count is greater or eaqual than the given k, true is sent in the channel, false otherwise
func checkAxisWith(c int, r int, board [][]string, name string, k int, pcs ...PointChecker) chan bool {
	res := make(chan bool)
	go func() {
		cntPcs := make([]int, len(pcs), len(pcs))
		for dist := 1; dist < k; dist++ {
			for idx, pc := range pcs {
				// counter on pointchecker guard their next invocation
				if cntPcs[idx] == dist-1 && pc(c, r, board, name, dist) {
					cntPcs[idx]++
				}
			}
		}
		sum := 0
		for _, c := range cntPcs {
			sum += c
		}
		res <- (sum + 1) >= k // count of connected neighbors + self greater k?
	}()
	return res
}

// signature: col, row, board, name, distance
// check if 'name' is found on a point on the board with distance d from point (c, r)
type PointChecker func(int, int, [][]string, string, int) bool

func checkCrossRightUp(c int, r int, board [][]string, name string, d int) bool {
	return c+d < len(board) && r+d < len(board[c+d]) && name == board[c+d][r+d]
}

func checkCrossRightDown(c int, r int, board [][]string, name string, d int) bool {
	return c+d < len(board) && r-d >= 0 && r-d < len(board[c+d]) && name == board[c+d][r-d]
}

func checkCrossLeftDown(c int, r int, board [][]string, name string, d int) bool {
	return c-d >= 0 && r-d >= 0 && r-d < len(board[c-d]) && name == board[c-d][r-d]
}

func checkCrossLeftUp(c int, r int, board [][]string, name string, d int) bool {
	return c-d >= 0 && r+d < len(board[c-d]) && name == board[c-d][r+d]
}

func checkLeftOf(c int, r int, board [][]string, name string, d int) bool {
	return c-d >= 0 && r < len(board[c-d]) && name == board[c-d][r]
}

func checkRightOf(c int, r int, board [][]string, name string, d int) bool {
	return c+d < len(board) && r < len(board[c+d]) && name == board[c+d][r]
}

func checkBelow(c int, r int, board [][]string, name string, d int) bool {
	return r-d >= 0 && r-d < len(board[c]) && name == board[c][r-d]
}
