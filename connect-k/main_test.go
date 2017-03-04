package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"testing"
)

func boardEmpty(board [][]bool, cols int) bool {
	for i := 0; i < cols; i++ {
		if len(board[i]) != 0 {
			return false
		}
	}
	return true
}

func getBoardFromFile(fileName string) [][]bool {
	// not meant to handle mal-formed files.
	file, _ := os.Open("testboards/" + fileName)
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	rows := strings.Split(string(b), "\n")
	//reverse
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
	board := make([][]bool, len(strings.Split(rows[0], "|")))
	for _, row := range rows {
		entries := strings.Split(row, "|")
		for c, e := range entries {
			val := false
			if e == "1" {
				val = true
			}
			if e != " " && (e == "0" || e == "1") {
				board[c] = append(board[c], val)
			}
		}
	}
	return board
}

// first row is filled with true, second false, alternating..
func getDummyBoard(cols, rows int) [][]bool {
	// fillup!
	board := make([][]bool, cols)
	token := false
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			board[i] = append(board[i], token)
			token = !token
		}
	}
	return board
}

func TestNewGame(t *testing.T) {
	cols, rows, k, mode, level := 10, 10, 5, "intelligent", 8
	g := NewGame(cols, rows, k, mode, level)

	if g.Id == "" || g.Mode != mode || g.Winner != "" || len(g.Board) != cols || g.Cols != cols || g.Rows != rows || g.Win != k || g.Level != level {
		t.Error("Game not properly initialized")
	}
	if !boardEmpty(g.Board, g.Cols) {
		t.Error("Expected game board to be empty after init, but have", g.Board)
	}

}

func TestNewGameStore(t *testing.T) {
	gs := NewGameStore()
	if gs == nil || gs.Games == nil {
		t.Error("Gamestore not initialized")
	}
}

func TestGameStoreAccess(t *testing.T) {
	gs := NewGameStore()

	game1 := NewGame(7, 6, 4, "intelligent", 4)

	if gs.Set(game1) == false {
		t.Error("Unable to set game in gamestore")
	}
	if gs.Get(game1.Id) == nil {
		t.Error("Unable to get game from gamestore")
	}
	if gs.Get("Some ID") != nil {
		t.Error("Found unknown game in gamestore")
	}
	if len(gs.Games) != 1 {
		t.Error("Expected 1 game in gamestore, got", len(gs.Games))
	}
}

func TestErrorHandler(t *testing.T) {
	t.Fail()
}

func TestAllowCors(t *testing.T) {
	t.Fail()
}

func TestNewHandler(t *testing.T) {
	t.Fail()
}

func TestPlayHandler(t *testing.T) {
	t.Fail()
}

func TestInspectHandler(t *testing.T) {
	t.Fail()
}

func TestSendResult(t *testing.T) {
	t.Fail()
}

func TestToIntOr(t *testing.T) {

	a := toIntOr("123", 0)
	if a != 123 {
		t.Error("String to int conversion not working properly")
	}
	b := toIntOr("12345678909876543211234567890987654321", 123)
	c := toIntOr("abc", 123)
	if b != 123 || c != 123 {
		t.Error("Expected 'or' fallback to kick in for unparsable strings")
	}
}

func TestApplySimple(t *testing.T) {
	board := make([][]bool, 7)

	apply(board, 0, true, 6)

	if len(board[0]) != 1 || board[0][0] != true {
		t.Error("Error applying token to board. Expected 'true' at board[0][0], got", board[0][0])
	}
}

func TestApplyFillBoard(t *testing.T) {
	// test apply return value
	cols, rows := 10, 8
	board := make([][]bool, cols)
	token := false
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			row, err := apply(board, i, token, rows)
			if err != nil {
				t.Error("Unexpected error when applying", err)
			}
			if row != j {
				t.Error("Applied row differs with expected", row, j)
			}
			if board[i][j] != token {
				t.Error("Applied token differs with present token", token, board[i][j])
			}
			if len(board[i]) != j+1 {
				// applied row + 1 is the length
				t.Error("Expected board length differs with actual", len(board[i]), j)
			}
			token = !token
		}
	}
}

func TestApplyOverflow(t *testing.T) {
	// test apply return value
	cols, rows := 10, 8
	token := false

	board := getDummyBoard(cols, rows)

	// cause overflow in each row
	for i := 0; i < cols; i++ {
		row, err := apply(board, i, token, rows)
		if row != -1 || err == nil {
			t.Error("Expected full row to not be applicable any more. Applied", row, i)
		}
	}
	// cause overflow left/right out of board
	row, err := apply(board, -1, token, rows)
	if row != -1 || err == nil {
		t.Error("Expected error when applying to col -1, but no error was thrown")
	}
	row, err = apply(board, cols+1, token, rows)
	if row != -1 || err == nil {
		t.Error("Expected error when applying to cols+1, but no error was thrown")
	}
}

func TestAutoPlay(t *testing.T) {
	cols, rows, k, mode, level := 7, 6, 4, "intelligent", 4

	g := NewGame(cols, rows, k, mode, level)

	c, r, err := autoPlay(g)
	if boardEmpty(g.Board, cols) {
		t.Error("Autoplay failed to insert token on board, got", g.Mode, g.Board)
	}
	if err != nil {
		t.Error("Unexpected error during autoplay", g.Mode, err)
	}

	mode = "random"
	g = NewGame(cols, rows, k, mode, level)
	c, r, err = autoPlay(g)

	if boardEmpty(g.Board, cols) {
		t.Error("Autoplay failed to insert token on board, got", g.Mode, g.Board)
	}
	if err != nil {
		t.Error("Unexpected error during autoplay", g.Mode, err)
	}

	mode = "UNKNOWN"
	g = NewGame(cols, rows, k, mode, level)
	c, r, err = autoPlay(g)

	if !boardEmpty(g.Board, cols) {
		t.Error("Autoplay inserted token with unknown playmode. Expected error", g.Mode, g.Board)
	}
	if c != -1 || r != -1 || err == nil {
		t.Error("Exptected error on autoplay with unknonw playmode, but got", c, r, g.Board)
	}
}

func TestPlayRandom(t *testing.T) {
	cols, rows, k, mode, level := 7, 6, 4, "random", 4
	game := NewGame(cols, rows, k, mode, level)
	playRandom(game)
	if boardEmpty(game.Board, cols) {
		t.Error("Board empty after playing random mode")
	}
	game.Board = getDummyBoard(cols, rows)

	//delete last elm
	game.Board[cols-1] = game.Board[cols-1][:rows-1]

	// playrandom has to fill that last col without errors
	c, r, err := playRandom(game)

	if c != cols-1 || r != rows-1 || err != nil {
		t.Error("Play random encountered unexpected error", c, r, err)
	}

	for i := 0; i < cols; i++ {
		if len(game.Board[i]) != rows {
			t.Error("Play random did not fill last place on board accordingly", game.Board)
		}
	}
}

func TestPlayIntelligent(t *testing.T) {
	t.Fail()
}

func TestScoreInDepth(t *testing.T) {
	t.Fail()
}

func TestAlphaBeta(t *testing.T) {
	t.Fail()
}

func TestInsertWins(t *testing.T) {
	game := NewGame(7, 6, 4, "intelligent", 4)

	game.Board = getBoardFromFile("insertWins/insertWins_3_0_clu")
	assert.True(t, insertWins(3, 0, true, game))

	game.Board = getBoardFromFile("insertWins/insertWins_3_0_cru")
	assert.True(t, insertWins(3, 0, true, game))

	game.Board = getBoardFromFile("insertWins/insertWins_3_3_below")
	assert.True(t, insertWins(3, 3, true, game))

	game.Board = getBoardFromFile("insertWins/insertWins_3_3_cld")
	assert.True(t, insertWins(3, 3, true, game))

	game.Board = getBoardFromFile("insertWins/insertWins_3_3_crd")
	assert.True(t, insertWins(3, 3, true, game))

	game.Board = getBoardFromFile("insertWins/insertWins_3_3_left")
	assert.True(t, insertWins(3, 3, true, game))

	game.Board = getBoardFromFile("insertWins/insertWins_3_3_right")
	assert.True(t, insertWins(3, 3, true, game))
}

func TestBoardFull(t *testing.T) {
	game := NewGame(7, 6, 4, "intelligent", 4)
	game.Board = getDummyBoard(7, 6)
	assert.True(t, boardFull(game))
	game.Board[6] = game.Board[6][:5]
	assert.False(t, boardFull(game))
	game.Board = game.Board[:6]
	assert.False(t, boardFull(game))
}

func TestScoreInsertAt(t *testing.T) {
	t.Fail()
}

func TestScoreAxis(t *testing.T) {
	t.Fail()
}

func TestSumArray(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Equal(t, sum(arr, 0, 10), 45)
	assert.Equal(t, sum(arr, 10, 0), 0)
	assert.Equal(t, sum(arr, 5, 5), 0)
	assert.Equal(t, sum(arr, 3, 7), 18)
}

func getPointCheckers() (pcs map[string]PointChecker) {
	pcs["a"] = checkAbove
	pcs["b"] = checkBelow
	pcs["r"] = checkRightOf
	pcs["l"] = checkLeftOf
	pcs["cru"] = checkCrossRightUp
	pcs["clu"] = checkCrossLeftUp
	pcs["crd"] = checkCrossRightDown
	pcs["cld"] = checkCrossLeftDown
	return
}

func TestScoreNeighbors(t *testing.T) {
	directions := map[string]string{"0:0": "right", "0:4": "below"}//, "cru"}
	neighborValues := map[string]map[int]int{
		"empty-empty-empty":   map[int]int{1: 3, 2: 2, 3: 1},
		"empty-empty-foe":     map[int]int{1: 3, 2: 2},
		"empty-empty-friend":  map[int]int{1: 3, 2: 2, 3: 8},
		"empty-foe":           map[int]int{1: 3},
		"empty-friend-empty":  map[int]int{1: 3, 2: 12, 3: 1},
		"empty-friend-foe":    map[int]int{1: 3, 2: 12},
		"empty-friend-friend": map[int]int{1: 3, 2: 12, 3: 8},
		"foe":                  map[int]int{},
		"friend-empty-empty":   map[int]int{1: 16, 2: 2, 3: 1},
		"friend-empty-foe":     map[int]int{1: 16, 2: 2},
		"friend-empty-friend":  map[int]int{1: 16, 2: 2, 3: 8},
		"friend-friend-empty":  map[int]int{1: 16, 2: 12, 3: 1},
		"friend-friend-foe":    map[int]int{1: 16, 2: 12},
		"friend-friend-friend": map[int]int{1: 16, 2: 12, 3: 8},
	}

	check := func(expS, s map[int]int, expAfc, afc int, dir, neighbors string) {
		assert.Equal(t, expS, s, "Unexpected neighbor scoring for "+dir+"_"+neighbors)
		assert.Equal(t, expAfc, afc, "Unexpected friendcount for "+dir+"_"+neighbors)
	}
	var board [][]bool
	for pos, dir := range directions {
		adjFriends := 0 // this will be so for a while.
		poss := strings.Split(pos, ":")
		x, _ := strconv.Atoi(poss[0])
		y, _ := strconv.Atoi(poss[1])

		for neighbors, expScores := range neighborValues {
			fmt.Println(dir + neighbors)
			board = getBoardFromFile("scoreNeighbors/" + dir + "_" + neighbors)
			fmt.Println(board)
			if strings.HasPrefix(neighbors, "friend-friend-friend") {
				adjFriends = 3
			} else if strings.HasPrefix(neighbors, "friend-friend") {
				adjFriends = 2
			} else if strings.HasPrefix(neighbors, "friend") {
				adjFriends = 1
			} else {
				adjFriends = 0
			}
			s, afc := scoreNeighbors(x, y, true, board, 4, 1, checkRightOf, "r")
			check(expScores, s, adjFriends, afc, dir, neighbors)
		}
	}
}

func TestPointCheckers(t *testing.T) {
	// board looks like: first row filled with true, second with false, alternating...
	cols, rows := 7, 6
	board := getDummyBoard(cols, rows)

	// self reflexity
	assert.Equal(t, FRIEND, checkCrossRightUp(0, 0, board, false, 0, rows), "CRU 0 0 Checking point with same value should result in 'FRIEND'")
	assert.Equal(t, FRIEND, checkCrossRightDown(0, 0, board, false, 0, rows), "CRD 0 0 Checking point with same value should result in 'FRIEND'")
	assert.Equal(t, FRIEND, checkCrossLeftUp(0, 0, board, false, 0, rows), "CLU 0 0 Checking point with same value should result in 'FRIEND'")
	assert.Equal(t, FRIEND, checkCrossLeftDown(0, 0, board, false, 0, rows), "CLD 0 0 Checking point with same value should result in 'FRIEND'")
	assert.Equal(t, FRIEND, checkRightOf(0, 0, board, false, 0, rows), "CR 0 0 Checking point with same value should result in 'FRIEND'")
	assert.Equal(t, FRIEND, checkLeftOf(0, 0, board, false, 0, rows), "CL 0 0 Checking point with same value should result in 'FRIEND'")
	assert.Equal(t, FRIEND, checkAbove(0, 0, board, false, 0, rows), "CA 0 0 Checking point with same value should result in 'FRIEND'")
	assert.Equal(t, FRIEND, checkBelow(0, 0, board, false, 0, rows), "CB 0 0 Checking point with same value should result in 'FRIEND'")

	// sanity: defined & reachable point
	assert.Equal(t, FOE, checkCrossRightUp(3, 2, board, false, 1, rows), "CRU 3 2 Checking reachable point did not match expected token")
	assert.Equal(t, FOE, checkCrossRightDown(3, 2, board, false, 1, rows), "CRD 3 2 Checking reachable point did not match expected token")
	assert.Equal(t, FOE, checkCrossLeftUp(3, 2, board, false, 1, rows), "CLU 3 2 Checking reachable point did not match expected token")
	assert.Equal(t, FOE, checkCrossLeftDown(3, 2, board, false, 1, rows), "CLD 3 2 Checking reachable point did not match expected token")
	assert.Equal(t, FRIEND, checkRightOf(3, 2, board, false, 1, rows), "CR 3 2 Checking reachable point did not match expected token")
	assert.Equal(t, FRIEND, checkLeftOf(3, 2, board, false, 1, rows), "CL 3 2 Checking reachable point did not match expected token")
	assert.Equal(t, FOE, checkAbove(3, 2, board, false, 1, rows), "CA 3 2 Checking reachable point did not match expected token")
	assert.Equal(t, FOE, checkBelow(3, 2, board, false, 1, rows), "CB 3 2 Checking reachable point did not match expected token")

	//delete last two rows of elms
	for c := 0; c < cols; c++ {
		board[c] = board[c][:rows-2]
	}
	// sanity: empty fields:
	assert.Equal(t, EMPTY, checkCrossRightUp(3, 3, board, false, 1, rows), "CRU 3 3 Checking reachable but empty point did not match expected token")
	assert.Equal(t, EMPTY, checkCrossRightDown(3, 5, board, false, 1, rows), "CRD 3 5 Checking reachable but empty point did not match expected token")
	assert.Equal(t, EMPTY, checkCrossLeftUp(3, 3, board, false, 1, rows), "CLU 3 3 Checking reachable but empty point did not match expected token")
	assert.Equal(t, EMPTY, checkCrossLeftDown(3, 5, board, false, 1, rows), "CLD 3 5 Checking reachable but empty point did not match expected token")
	assert.Equal(t, EMPTY, checkRightOf(3, 4, board, false, 1, rows), "CR 3 4 Checking reachable but empty point did not match expected token")
	assert.Equal(t, EMPTY, checkLeftOf(3, 4, board, false, 1, rows), "CL 3 4 Checking reachable but empty point did not match expected token")
	assert.Equal(t, EMPTY, checkAbove(3, 3, board, false, 1, rows), "CA 3 3 Checking reachable but empty point did not match expected token")
	assert.Equal(t, EMPTY, checkBelow(3, 5, board, false, 1, rows), "CB 3 5 Checking reachable but empty point did not match expected token")

	// sanity: defined (on board) but unreachable point (nothing there to put a token onto)
	board = make([][]bool, 7)
	assert.Equal(t, UNR, checkCrossRightUp(3, 2, board, false, 1, rows), "CRU 3 2 Checking unreachable point did not match expected token")
	assert.Equal(t, UNR, checkCrossRightDown(3, 2, board, false, 1, rows), "CRD 3 2 Checking unreachable point did not match expected token")
	assert.Equal(t, UNR, checkCrossLeftUp(3, 2, board, false, 1, rows), "CLU 3 2 Checking unreachable point did not match expected token")
	assert.Equal(t, UNR, checkCrossLeftDown(3, 2, board, false, 1, rows), "CLD 3 2 Checking unreachable point did not match expected token")
	assert.Equal(t, UNR, checkRightOf(3, 2, board, false, 1, rows), "CR 3 2 Checking unreachable point did not match expected token")
	assert.Equal(t, UNR, checkLeftOf(3, 2, board, false, 1, rows), "CL 3 2 Checking unreachable point did not match expected token")
	assert.Equal(t, UNR, checkAbove(3, 2, board, false, 1, rows), "CA 3 2 Checking unreachable point did not match expected token")
	assert.Equal(t, UNR, checkBelow(3, 2, board, false, 1, rows), "CB 3 2 Checking unreachable point did not match expected token")

	// sanity: columns out of bounds
	assert.Equal(t, OOB, checkCrossRightUp(6, 2, board, false, 1, rows), "CRU 6 2 Checking columns out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkCrossRightDown(6, 2, board, false, 1, rows), "CRD 6 2 Checking columns out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkCrossLeftUp(0, 2, board, false, 1, rows), "CLU 0 2 Checking columns out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkCrossLeftDown(0, 2, board, false, 1, rows), "CLD 0 2 Checking columns out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkRightOf(6, 2, board, false, 1, rows), "CR 6 2 Checking columns out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkLeftOf(0, 2, board, false, 1, rows), "CL 0 2 Checking columns out-of-bound point did not match expected token")

	// sanity: rows out of bounds
	assert.Equal(t, OOB, checkCrossRightUp(3, 5, board, false, 1, rows), "CRU 3 5 Checking rows out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkCrossRightDown(3, 0, board, false, 1, rows), "CRD 3 0 Checking rows out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkCrossLeftUp(3, 5, board, false, 1, rows), "CLU 3 5 Checking rows out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkCrossLeftDown(3, 0, board, false, 1, rows), "CLD 3 0 Checking rows out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkAbove(0, 5, board, false, 1, rows), "CA 0 5 Checking rows out-of-bound point did not match expected token")
	assert.Equal(t, OOB, checkBelow(0, 0, board, false, 1, rows), "CB 0 0 Checking rows out-of-bound point did not match expected token")
}

func TestToAxis(t *testing.T) {
	map1, map2 := make(map[int]int, 5), make(map[int]int, 5)
	for i := 0; i < 5; i++ {
		map1[i] = i + 1
		map2[i] = i + 1
	}
	axis := toAxis(map1, map2, 0)
	expected := []int{5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5}
	assert.Equal(t, axis, expected, "Axis creation failed")
}

func TestGenShuff(t *testing.T) {

	max := 10
	randomSlice := genShuff(max)

	for i := 0; i < max; i++ {
		if !(randomSlice[i] < 10 && randomSlice[i] >= 0) {
			t.Error("Generated Shuffled slice contained unexpected number", randomSlice[i])
		}
	}

	randomSlice = genShuff(0)
	assert.Equal(t, randomSlice, make([]int, 0), "Expected empty slice")
}
