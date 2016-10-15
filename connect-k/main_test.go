package main

import "testing"

//import "fmt"

func TestNewGame(t *testing.T) {
	cols, rows, k, mode, level := 7, 6, 4, "intelligent", 4
	g := NewGame(cols, rows, k, mode, level)

	if g.Id == "" || g.Mode != mode || g.Winner != "" || len(g.Board) != 7 || g.Cols != cols || g.Rows != rows || g.Win != k || g.Level != level {
		t.Error("Game not properly initialized")
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

func TestApplyFullBoard(t *testing.T) {
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
	board := make([][]bool, cols)
	token := false
	// fillup!
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			_, err := apply(board, i, token, rows)
			if err != nil {
				t.Error("Unexpected error when applying", err)
			}
			token = !token
		}
	}
	// cause overflow in each row
	for i := 0; i < cols; i++ {
		row, err := apply(board, i, token, rows)
		if row != -1 || err == nil {
			t.Error("Expected full row to not be applicable any more. Applied", row, i)
		}
	}

	row, err := apply(board, -1, token, rows)
	if row != -1 || err == nil {
		t.Error("Expected error when applying to col -1, but no error was thrown")
	}
	row, err = apply(board, cols+1, token, rows)
	if row != -1 || err == nil {
		t.Error("Expected error when applying to cols+1, but no error was thrown")
	}

}
