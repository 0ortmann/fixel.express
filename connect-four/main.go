package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
)

type GameStore struct {
	Games map[string]*Game
	mu    sync.RWMutex
}

type Game struct {
	Id    string
	Mode  string
	Board [7][6]string
}

func NewGame() *Game {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error generating UUID ", err)
		return nil
	}
	return &Game{
		Id:   fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]),
		Mode: "random", // minmax, heuristics etc..
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
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/play", playHandler)

	http.ListenAndServe(":5000", nil)
}

func newHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "new")
	game := NewGame()
	gs.Set(game)
	enc := json.NewEncoder(w)
	enc.Encode(game)
}

func playHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	type Post struct {
		GameId string
		Col    int
	}

	var post Post
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&post)

	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}
	if post.GameId == "" || post.Col >= 7 {
		http.Error(w, "Wrong JSON format", 400)
		return
	}
	game := gs.Get(post.GameId)

	if game == nil {
		http.Error(w, "No game with that ID", 404)
		return
	}
	err = apply(game, post.Col, "player")

	if err != nil {
		http.Error(w, "Field already in use", 409)
		return
	}

	col, err := autoPlay(game)
	if err != nil {
		http.Error(w, "Unable to play", 400)
	}

	win := checkWin(game)
	type Resp struct {
		Col int
		Winner string
	}
	var resp Resp
	resp.Col = col
	resp.Winner = win
	enc := json.NewEncoder(w)
	enc.Encode(&resp)
}

// put a token with value "player" into the column "col".
// Returns an error if no game with that id exists of the column is already full
func apply(game *Game, col int, player string) error {
	rc := count(game.Board[col])
	if rc == 6 {
		return errors.New("Column already full")
	}
	game.Board[col][rc] = player
	fmt.Println(game)
	return nil
}

// The computer plays automatically and returns the column number that was picked to insert a chip
func autoPlay(game *Game) (int, error) {
	switch game.Mode {
	case "random": 
		return playRandom(game)
	}
	return 0, nil
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

func checkWin(game *Game) string {
	return ""
}


func count(arr [6]string) (c int) {
	for i := 0; i < len(arr); i++ {
		if arr[i] == "" {
			return
		}
		c++
	}
	return
}
