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
		return errors.New("Game already ended, winner was "+game.Winner)
	}
	if err := apply(game, post.Col, "player"); err != nil {
		return errors.New("Field already in use")
	}

	// fixme: player win -> computer does a last move
	col, err := autoPlay(game)
	if err != nil {
		return errors.New("Unable to play")
	}

	checkWin(game)

	type Resp struct {
		Col    int `json:"col"`
		Winner string `json:"winner"`
	}
	var resp Resp
	resp.Col = col
	resp.Winner = game.Winner
	enc := json.NewEncoder(w)
	enc.Encode(&resp)
	return nil
}

// put a token with value "player" into the column "col".
// Returns an error if the column is already full
func apply(game *Game, col int, player string) error {
	if len(game.Board[col]) == 6 {
		return errors.New("Column already full")
	}
	game.Board[col] = append(game.Board[col], player)
	fmt.Println(game)
	return nil
}

// The computer plays automatically and returns the column number that was picked to insert a chip
func autoPlay(game *Game) (int, error) {
	switch game.Mode {
	case "random":
		return playRandom(game)
	}
	return -1, errors.New("Unknown game mode")
}

func playRandom(game *Game) (int, error) {
	// has to terminate if all cols are full
	rCol, _ := rand.Int(rand.Reader, big.NewInt(7))
	c := int(rCol.Int64())
	for i := 0; i < 6; i++ {
		c = (c + i) % 6
		err := apply(game, c, "computooor")
		if err == nil {
			return c, nil
		}
	}
	return -1, errors.New("No more usable columns, all full!")
}

func checkWin(game *Game) {
	game.Winner = getWinner(game.Board, 4)

}
func getWinner(board [][]string, k int) (w string) {
	for c := 0; c < len(board); c++ {
		for r := 0; r < len(board[c]); r++ {
			w = checkPoint(c, r, board, k)
			if w != "" {
				return
			}
		}
	}
	return
}

// check above, right besides, cross-right-down and cross-right-up
func checkPoint(x int, y int, board [][]string, k int) string {
	name := board[x][y]
	aCnt := 0
	bCnt := 0
	cuCnt := 0
	cdCnt := 0
	fmt.Println("checkpoint ", x, y, name)
	for i := 0; i < k; i++ {
		//above
		if x < len(board) && y+i < len(board[x]) {
			fmt.Println("in above")
			if name == board[x][y+i] {
				aCnt++
			}
		}
		//besides
		if x+i < len(board) && y < len(board[x+i]) {
			fmt.Println("in besides")
			if name == board[x+i][y] {
				bCnt++
			}
		}
		//cross up
		if x+i < len(board) && y+i < len(board[x+i]) {
			fmt.Println("in cu")
			if name == board[x+i][y+i] {
				cuCnt++
			}
		}
		//cross down
		if x+i < len(board) && y-i >= 0 && y-i < len(board[x+i]) {
			fmt.Println("in cd")
			if name == board[x+i][y-i] {
				cdCnt++
			}
		}
	}
	if aCnt == k || bCnt == k || cuCnt == k || cdCnt == k {
		return name
	}
	return ""
}