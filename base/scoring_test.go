package base

import (
	"testing"
)

func checkValues(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestScoring(t *testing.T) {
	board := makeBoard(5, 5)
	board.placeStone(Black, Point{1, 2})
	board.placeStone(Black, Point{1, 4})
	board.placeStone(Black, Point{2, 2})
	board.placeStone(Black, Point{2, 3})
	board.placeStone(Black, Point{2, 4})
	board.placeStone(Black, Point{2, 5})
	board.placeStone(Black, Point{3, 1})
	board.placeStone(Black, Point{3, 2})
	board.placeStone(Black, Point{3, 3})
	board.placeStone(White, Point{3, 4})
	board.placeStone(White, Point{3, 5})
	board.placeStone(White, Point{4, 1})
	board.placeStone(White, Point{4, 2})
	board.placeStone(White, Point{4, 3})
	board.placeStone(White, Point{4, 4})
	board.placeStone(White, Point{5, 2})
	board.placeStone(White, Point{5, 4})
	board.placeStone(White, Point{5, 5})
	// territory := scoring.EvaluateTerritory(board)
	// checkValues(t, 9, territory.num_black_stones)
	// checkValues(t, 9, territory.num_black_stones)
	// checkValues(t, 4, territory.num_black_territory)
	// checkValues(t, 9, territory.num_white_stones)
	// checkValues(t, 3, territory.num_white_territory)
	// checkValues(t, 0, territory.num_dame)
}
