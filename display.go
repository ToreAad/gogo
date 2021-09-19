package main

import (
	"fmt"
	"gogo/base"
	"strings"
)

var COLS = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"}

var STONE_TO_CHAR = map[base.PlayerColor]string{
	base.Empty: " .",
	base.Black: " x",
	base.White: " 0",
}

func printMove(player base.PlayerColor, move base.Move) {
	switch m := move.(type) {
	case *base.Pass:
		fmt.Printf("%s %s", STONE_TO_CHAR[player], "passes")
	case *base.Resign:
		fmt.Printf("%s %s", STONE_TO_CHAR[player], "resigns")
	case *base.Play:
		fmt.Printf("%s %s%d", STONE_TO_CHAR[player], COLS[m.Point.Col-1], m.Point.Row)
	}
}

func printBoard(b *base.Board) {
	for row := b.Rows; row > 0; row-- {
		line := make([]string, 0)
		if row <= 9 {
			line = append(line, " ")
		}
		line = append(line, fmt.Sprintf("%d ", row))
		for col := 1; col <= b.Cols; col++ {
			stone, _ := b.Get(base.Point{Row: row, Col: col})
			line = append(line, STONE_TO_CHAR[stone])
		}
		for _, s := range line {
			fmt.Print(s)
		}
		fmt.Print("\n")
	}
	fmt.Printf("\n    %s\n", strings.Join(COLS[:b.Cols], " "))
}
