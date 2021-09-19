package agent

import (
	"gogo/base"
	"math/rand"
	"time"
)

func isPointAnEye(board base.Board, point base.Point, color base.PlayerColor) (bool, error) {
	if !board.IsEmpty(point) {
		return false, nil
	}

	for _, neighbour := range point.Neighbours() {
		if board.IsOnBoard(neighbour) {
			neighbourColor, _ := board.Get(neighbour)
			if neighbourColor != color {
				return false, nil
			}
		}
	}

	friendlyCorners := 0
	offBoardCorners := 0
	corners := []base.Point{
		{Row: point.Row - 1, Col: point.Col - 1},
		{Row: point.Row - 1, Col: point.Col + 1},
		{Row: point.Row + 1, Col: point.Col - 1},
		{Row: point.Row + 1, Col: point.Col + 1}}

	for _, corner := range corners {
		cornerColor, _ := board.Get(corner)
		if cornerColor == color {
			friendlyCorners++
		} else {
			offBoardCorners++
		}
	}
	if offBoardCorners > 0 {
		return offBoardCorners+friendlyCorners == 4, nil
	}
	return friendlyCorners >= 3, nil
}

type Agent interface {
	SelectMove(gs *base.GameState) base.Move
}

type RandomBot struct{}

func (r *RandomBot) SelectMove(gs *base.GameState) base.Move {
	candidates := make([]base.Play, 0)
	for r := 1; r <= gs.Board.Rows; r++ {
		for c := 1; c <= gs.Board.Cols; c++ {
			candidate := base.Play{Point: base.Point{Row: r, Col: c}}
			validMove, _ := gs.IsValidMove(&candidate)
			isAnEye, _ := isPointAnEye(*gs.Board, candidate.Point, gs.NextPlayer)
			if validMove && !isAnEye {
				candidates = append(candidates, candidate)
			}
		}
	}

	if len(candidates) == 0 {
		return &base.Pass{}
	}

	rand.Seed(time.Now().Unix())
	return &candidates[rand.Intn(len(candidates))]
}
