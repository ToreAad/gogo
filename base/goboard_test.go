package base

import (
	"testing"
)

func checkBoard(t *testing.T, board *Board, p Point, expected PlayerColor) {
	if pc, err := board.Get(p); err != nil || expected != pc {
		t.Errorf("Expected %v got %v %v", expected, pc, err)
	}
}

func boardIsNil(t *testing.T, board *Board, p Point) {
	if pc, err := board.Get(p); err == nil || Empty != pc {
		t.Errorf("Expected %v got %v %v", Empty, pc, err)
	}
}

func checkLiberties(t *testing.T, board *Board, p Point, libs []Point) {
	strings, err := board.getGoString(p)
	if err != nil {
		t.Errorf("Got an error %v", err)
	}
	if len(strings.liberties) != len(libs) {
		t.Errorf("Expected %v got %v", len(libs), len(strings.liberties))
	}
	for _, po := range libs {
		if !strings.liberties[po] {
			t.Errorf("Did not fild %v", po)
		}
	}
}

func TestCapture(t *testing.T) {
	board := makeBoard(19, 19)
	board.placeStone(Black, Point{2, 2})
	board.placeStone(White, Point{1, 2})
	checkBoard(t, board, Point{2, 2}, Black)
	board.placeStone(White, Point{2, 1})
	checkBoard(t, board, Point{2, 2}, Black)
	board.placeStone(White, Point{2, 3})
	checkBoard(t, board, Point{2, 2}, Black)
	board.placeStone(White, Point{3, 2})
	boardIsNil(t, board, Point{2, 2})
}

func TestCaptureTwoStones(t *testing.T) {
	board := makeBoard(19, 19)
	board.placeStone(Black, Point{2, 2})
	board.placeStone(Black, Point{2, 3})
	board.placeStone(White, Point{1, 2})
	board.placeStone(White, Point{1, 3})
	checkBoard(t, board, Point{2, 2}, Black)
	checkBoard(t, board, Point{2, 3}, Black)
	board.placeStone(White, Point{3, 2})
	board.placeStone(White, Point{3, 3})
	checkBoard(t, board, Point{2, 2}, Black)
	checkBoard(t, board, Point{2, 3}, Black)
	board.placeStone(White, Point{2, 1})
	board.placeStone(White, Point{2, 4})
	boardIsNil(t, board, Point{2, 2})
	boardIsNil(t, board, Point{2, 3})
}

func TestCaptureIsNotSuicide(t *testing.T) {
	board := makeBoard(19, 19)
	board.placeStone(Black, Point{1, 1})
	board.placeStone(Black, Point{2, 2})
	board.placeStone(Black, Point{1, 3})
	board.placeStone(White, Point{2, 1})
	board.placeStone(White, Point{1, 2})
	boardIsNil(t, board, Point{1, 1})
	checkBoard(t, board, Point{2, 1}, White)
	checkBoard(t, board, Point{1, 2}, White)
}

func TestRemoveLiberties(t *testing.T) {
	board := makeBoard(5, 5)
	board.placeStone(Black, Point{3, 3})
	board.placeStone(White, Point{2, 2})
	checkLiberties(t, board, Point{2, 2}, []Point{{2, 3}, {2, 1}, {1, 2}, {3, 2}})
	board.placeStone(Black, Point{3, 2})
	checkLiberties(t, board, Point{2, 2}, []Point{{2, 3}, {2, 1}, {1, 2}})
}

func TestEmptyTriangle(t *testing.T) {
	board := makeBoard(5, 5)
	board.placeStone(Black, Point{1, 1})
	board.placeStone(Black, Point{1, 2})
	board.placeStone(Black, Point{2, 2})
	board.placeStone(White, Point{2, 1})
	checkLiberties(t, board, Point{1, 1}, []Point{{3, 2}, {2, 3}, {1, 3}})
}

func TestSelfCapture(t *testing.T) {
	board := makeBoard(5, 5)
	board.placeStone(Black, Point{1, 1})
	board.placeStone(Black, Point{1, 3})
	board.placeStone(White, Point{2, 1})
	board.placeStone(White, Point{2, 2})
	board.placeStone(White, Point{2, 3})
	board.placeStone(White, Point{1, 4})
	if !board.isSelfCapture(Black, Point{1, 2}) {
		t.Errorf("Expected self capture")
	}
}

func TestNotSelfCapture(t *testing.T) {
	board := makeBoard(5, 5)
	board.placeStone(Black, Point{1, 1})
	board.placeStone(Black, Point{1, 3})
	board.placeStone(White, Point{2, 1})
	board.placeStone(White, Point{2, 3})
	board.placeStone(White, Point{1, 4})
	if board.isSelfCapture(Black, Point{1, 2}) {
		t.Errorf("Expected not self capture")
	}
}

func TestNotSelfCaptureIsOtherCapture(t *testing.T) {
	board := makeBoard(5, 5)
	board.placeStone(Black, Point{3, 1})
	board.placeStone(Black, Point{3, 2})
	board.placeStone(Black, Point{2, 3})
	board.placeStone(Black, Point{1, 1})
	board.placeStone(White, Point{2, 1})
	board.placeStone(White, Point{2, 2})
	board.placeStone(White, Point{1, 3})
	if board.isSelfCapture(Black, Point{1, 2}) != false {
		t.Errorf("Expected not self capture")
	}
}

func TestNewGame(t *testing.T) {
	start := NewGame(19)
	nextState, err := start.ApplyMove(&Play{Point{16, 16}})

	if err != nil || nextState.previousState != start {
		t.Errorf("Expected %v %v", nextState, err)
	}
	if err != nil || nextState.NextPlayer != White {
		t.Errorf("Expected %v %v", nextState, err)
	}
	checkBoard(t, nextState.Board, Point{16, 16}, Black)
}
