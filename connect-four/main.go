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
		fmt.Println("Error generating UUID ", err)
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
	checkWin(post.Col, pRow, game, 4)
	if game.Winner != "" {
		return playResult(w, game, -1)
	}
	cCol, cRow, err := autoPlay(game)
	if err != nil {
		return err
	}
	checkWin(cCol, cRow, game, 4)
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
	fmt.Println(game)
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

func checkWin(c int, r int, game *Game, k int) {
	name := game.Board[c][r]
	if checkInsertWins(c, r, game.Board, name, k) {
		game.Winner = name
	}
}

func checkInsertWins(c int, r int, board [][]string, name string, k int) bool {
	fmt.Println("check win", c, r, name)
	cr := checkCrossRight(c, r, board, name, k)
	cl := checkCrossLeft(c, r, board, name, k)
	rl := checkBesides(c, r, board, name, k)
	blw := checkBelow(c, r, board, name, k)
	succ := false
	axis := 0
	for {
		select {
		case res := <-cr:
			fmt.Println("res from cr", res)
			succ = succ || res
			axis++
		case res := <-cl:
			fmt.Println("res from cl", res)
			succ = succ || res
			axis++
		case res := <-rl:
			fmt.Println("res from rl", res)
			succ = succ || res
			axis++
		case res := <-blw:
			fmt.Println("res from blw", res)
			succ = succ || res
			axis++
		}
		if succ || axis == 4 {
			fmt.Println("hit", succ, axis)
			return succ
		}
	}
}

func checkCrossRight(c int, r int, board [][]string, name string, k int) chan bool {
	res := make(chan bool)
	cru := 0 // cross right up
	cld := 0 // cross left down
	go func() {
		for i := 1; i < k; i++ {
			if cru == i-1 && c+i < len(board) && r+i < len(board[c+i]) {
				if name == board[c+i][r+i] {
					cru++
				}
			}
			if cld == i-1 && c-i >= 0 && r-i >= 0 && r-i < len(board[c-i]) {
				if name == board[c-i][r-i] {
					cld++
				}
			}
		}
		res <- (cru + cld + 1) >= k
	}()
	return res
}

func checkCrossLeft(c int, r int, board [][]string, name string, k int) chan bool {
	res := make(chan bool)
	clu := 0 // cross left up
	crd := 0 // cross right down
	go func() {
		for i := 1; i < k; i++ {
			if clu == i-1 && c-i >= 0 && r+i < len(board[c-i]) {
				if name == board[c-i][r+i] {
					clu++
				}
			}
			if crd == i-1 && c+i < len(board) && r-i >= 0 && r-i < len(board[c+i]) {
				if name == board[c+i][r-i] {
					crd++
				}
			}
		}
		res <- (clu + crd + 1) >= k
	}()
	return res
}

func checkBesides(c int, r int, board [][]string, name string, k int) chan bool {
	res := make(chan bool)
	rgt := 0 // right
	lft := 0 // left
	go func() {
		for i := 1; i < k; i++ {
			if rgt == i-1 && c+i < len(board) && r < len(board[c+i]) {
				if name == board[c+i][r] {
					fmt.Println("add right", name, "==", board[c+i][r])
					rgt++
				}
			}
			if lft == i-1 && c-i >= 0 && r < len(board[c-i]) {
				if name == board[c-i][r] {
					fmt.Println("add left", name, "==", board[c-i][r])
					lft++
				}
			}
		}
		res <- (rgt + lft + 1) >= k
	}()
	return res
}

func checkBelow(c int, r int, board [][]string, name string, k int) chan bool {
	res := make(chan bool)
	blw := 0 // below
	go func() {
		for i := 1; i < k; i++ {
			if blw == i-1 && r-i >= 0 && r-i < len(board[c]) {
				if name == board[c][r-i] {
					blw++
				}
			}
		}
		res <- (blw + 1) >= k
	}()
	return res
}
