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
	game := NewGame(7, 6, 4, "random")
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
	pRow, err := apply(game, post.Col, "player")
	if err != nil {
		return err
	}
	setWinner(post.Col, pRow, game)
	if game.Winner != "" {
		return playResult(w, game, -1)
	}
	cCol, cRow, err := autoPlay(game)
	if err != nil {
		return err
	}
	setWinner(cCol, cRow, game)
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

// checks if the insert at position (c, r) lead to a win for the player with that token
// sets the Winner-field on the game, if so.
func setWinner(c int, r int, game *Game) {
	name := game.Board[c][r]
	if insertWins(c, r, name, game) {
		game.Winner = name
	}
}

// checks if an insert of the string 'name' at position (c, r) wins the game
func insertWins(c int, r int, name string, game *Game) bool {
	cr := rateAlongAxis(c, r, game, name, "c-r", checkCrossRightUp, checkCrossLeftDown)
	cl := rateAlongAxis(c, r, game, name, "c-l", checkCrossLeftUp, checkCrossRightDown)
	rl := rateAlongAxis(c, r, game, name, "r-l", checkRightOf, checkLeftOf)
	blw := rateAlongAxis(c, r, game, name, "blw", checkBelow)

	rate := 0
	axis := 0
	for r := range merge(cr, cl, rl, blw) {
		axis ++
		if rate < r {
			rate = r
		}
		if rate +1 >= game.Win {
			// rate of surrounding + yourself
			return true
		}
		if axis == 4 {
			return false
		}
	}
	return false
}


// rates the point (c, r) along an axis that is specified by the provided Pointcheckers
// calculates the rate for that axis and pushes it into the result channel. If the axis
// is not usable because one cannot win with that axis, it has value 0.
func rateAlongAxis(c int, r int, game *Game, name string, dir string, pcs ...PointChecker) <-chan int {
	res := make(chan int)
	go func() {
		subRates := make([]int, len(pcs), len(pcs))
		subDists := make([]int, len(pcs), len(pcs))
		for i, pc := range pcs {
			subRates[i], subDists[i] = ratePoint(c, r, game, name, pc)
			fmt.Println("rating for", c, r, dir, subRates[i])
		}
		rate, dist := 0, 0
		for i := 0; i < len(pcs); i++ {
			rate += subRates[i]
			dist += subDists[i]
		}
		if dist < game.Win {
			// not enough fields in a row usable for 'name', whole axis is not valuable
			res <- -1
			return
		}
		res <- rate // axis usable, return rate
	}()
	return res
}

// rates the point (c, r) in the context of 'game' and uses the given pointchecker.
// returns the rate and the max distance of fields that were considered for rating.
// rating function:
// introspect next neighbor, using the Pointchecker:
// found same chip as self: increments the rate by one
// found empty field: no influence to rate
// found field out of bounds or unfriendly chip: terminate checking
func ratePoint(c, r int, game *Game, name string, pc PointChecker) (int, int) {
	rate := 0
	dist := 1
	StraightInARow:
		for ; dist < game.Win; dist++ {
			switch pc(c, r, game.Board, name, dist, game.Rows) {
			case 1: // same chip
				rate++
			case 0: // no chip, not out of bounds
				break
			case -1, -2:
				break StraightInARow
			}
		}
	return rate, dist
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
	case c+d >= len(board) || r+d > avail:
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
	case c-d < 0 || r+d > avail:
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
	case c-d < 0 || r > avail:
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
	switch {
	case c+d >= len(board) || r > avail:
		return -1
	case r >= len(board[c+d]):
		return 0
	case name == board[c+d][r]:
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
