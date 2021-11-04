package base

import (
	"errors"
)

type Move interface {
	isMove()
}

type Play struct {
	Point Point
}
type Pass struct{}
type Resign struct{}

func (*Play) isMove()   {}
func (*Pass) isMove()   {}
func (*Resign) isMove() {}

type GoString struct {
	color     PlayerColor
	stones    map[Point]bool
	liberties map[Point]bool
}

func (gs *GoString) copy() *GoString {
	stones := make(map[Point]bool)
	liberties := make(map[Point]bool)

	for k, v := range gs.liberties {
		liberties[k] = v
	}

	for k, v := range gs.stones {
		stones[k] = v
	}

	return &GoString{gs.color, stones, liberties}
}

func (gs *GoString) withoutLiberty(p Point) *GoString {
	new_string := gs.copy()
	new_string.removeLiberty(p)
	return new_string
}

func (gs *GoString) removeLiberty(p Point) {
	delete(gs.liberties, p)
}

func (gs *GoString) addLiberty(p Point) {
	gs.liberties[p] = true
}

func (gs *GoString) mergedWith(other *GoString) (*GoString, error) {
	if gs.color != other.color {
		return nil, errors.New("string types are not the same")
	}
	combinedStones := make(map[Point]bool)
	for k, v := range gs.stones {
		combinedStones[k] = v
	}
	for k, v := range other.stones {
		combinedStones[k] = v
	}

	newLiberties := make(map[Point]bool)
	for k, v := range gs.liberties {
		if combinedStones[k] {
			continue
		}
		newLiberties[k] = v
	}
	for k, v := range other.liberties {
		if combinedStones[k] {
			continue
		}
		newLiberties[k] = v
	}
	return &GoString{color: gs.color, stones: combinedStones, liberties: newLiberties}, nil
}

func (gs *GoString) equals(o *GoString) bool {
	if gs == o {
		return true
	}
	if gs == nil || o == nil {
		return gs == nil && o == nil
	}
	if gs.color != o.color {
		return false
	}
	if len(gs.liberties) != len(o.liberties) {
		return false
	}
	if len(gs.stones) != len(o.stones) {
		return false
	}
	for k := range gs.stones {
		if o.stones[k] {
			continue
		}
		return false
	}
	for k := range gs.liberties {
		if o.liberties[k] {
			continue
		}
		return false
	}
	return true
}

type Board struct {
	Rows     int
	Cols     int
	grid     map[Point]*GoString
	hash     int
	zoobrist map[state]int
}

func makeBoard(r int, c int) *Board {
	return &Board{r, c, make(map[Point]*GoString), 0, getHashes(r, c)}
}

func (b *Board) copy() *Board {
	grid := make(map[Point]*GoString)
	for p, gs := range b.grid {
		grid[p] = gs.copy()
	}
	return &Board{b.Rows, b.Cols, grid, b.hash, b.zoobrist}
}

func (b *Board) equals(o *Board) bool {
	if b == nil || o == nil {
		return b == nil && o == nil
	}
	if !(b.Rows == o.Rows && b.Cols == o.Cols && len(b.grid) == len(o.grid)) {
		return false
	}
	for k, v := range b.grid {
		ov, e := o.grid[k]
		if e || !ov.equals(v) {
			return false
		}
	}
	return true
}

func (b *Board) IsOnBoard(p Point) bool {
	return 1 <= p.Row && p.Row <= b.Rows && 1 <= p.Col && p.Col <= b.Cols
}

func (b *Board) Get(p Point) (PlayerColor, error) {
	v, ok := b.grid[p]
	if !ok {
		return Empty, errors.New("grid position is empty")
	}
	return v.color, nil
}

func (b *Board) IsEmpty(p Point) bool {
	boardColor, _ := b.Get(p)
	return boardColor == Empty
}

func (b *Board) getGoString(p Point) (*GoString, error) {
	v, ok := b.grid[p]
	if !ok {
		return nil, errors.New("grid position is empty")
	}
	return v, nil
}

func (b *Board) replaceString(newString *GoString) {
	for p := range newString.stones {
		b.grid[p] = newString
	}
}

func (b *Board) removeString(gs *GoString) {
	for p := range gs.stones {
		for _, neighbour := range p.Neighbours() {
			neigbour_string, err := b.getGoString(neighbour)
			if err != nil {
				continue
			}
			if neigbour_string != gs {
				neigbour_string.addLiberty(p)
				b.replaceString(neigbour_string)
			}
		}
		b.hash ^= b.zoobrist[state{p, gs.color}]
		delete(b.grid, p)
	}
}

func contains(arr []*GoString, gs *GoString) bool {
	for _, e := range arr {
		if e.equals(gs) {
			return true
		}
	}
	return false
}

func (b *Board) isSelfCapture(player PlayerColor, point Point) bool {
	nextBoard := b.copy()
	nextBoard.placeStone(player, point)
	newString, _ := nextBoard.getGoString(point)
	return newString != nil && newString.liberties != nil && len(newString.liberties) == 0
}

func (b *Board) placeStone(player PlayerColor, p Point) error {
	if !b.IsOnBoard(p) {
		return errors.New("point not on board")
	}
	if _, err := b.Get(p); err == nil {
		return errors.New("grid position must be empty")
	}

	adjacentSameColor := make([]*GoString, 0)
	adjacentOppositeColor := make([]*GoString, 0)
	liberties := make(map[Point]bool)

	for _, neighbor := range p.Neighbours() {
		if !b.IsOnBoard(neighbor) {
			continue
		}
		neighbor_string, err := b.getGoString(neighbor)
		if err != nil {
			liberties[neighbor] = true
		} else if neighbor_string.color == player {
			if !contains(adjacentSameColor, neighbor_string) {
				adjacentSameColor = append(adjacentSameColor, neighbor_string)
			}
		} else {
			if !contains(adjacentOppositeColor, neighbor_string) {
				adjacentOppositeColor = append(adjacentOppositeColor, neighbor_string)
			}
		}
	}
	stones := make(map[Point]bool)
	stones[p] = true
	new_string := &GoString{player, stones, liberties}
	var err error

	for _, sameColorString := range adjacentSameColor {
		new_string, err = new_string.mergedWith(sameColorString)
		if err != nil {
			return err
		}
	}
	for new_string_point := range new_string.stones {
		b.grid[new_string_point] = new_string
	}
	b.hash ^= b.zoobrist[state{p, Empty}]
	b.hash ^= b.zoobrist[state{p, player}]
	for _, other_color_string := range adjacentOppositeColor {
		replacement_string := other_color_string.withoutLiberty(p)
		if len(replacement_string.liberties) > 0 {
			b.replaceString(replacement_string)
		} else {
			b.removeString(other_color_string)
		}
	}
	return nil
}

type GameState struct {
	Board         *Board
	NextPlayer    PlayerColor
	previousState *GameState
	lastMove      Move
}

func NewGame(n int) *GameState {
	return &GameState{makeBoard(n, n), Black, nil, nil}
}

func (g *GameState) ApplyMove(move Move) (*GameState, error) {
	var nextBoard *Board
	switch m := move.(type) {
	case *Play:
		nextBoard = g.Board.copy()
		err := nextBoard.placeStone(g.NextPlayer, m.Point)
		if err != nil {
			return nil, err
		}
	default:
		nextBoard = g.Board
	}
	return &GameState{nextBoard, g.NextPlayer.otherPlayer(), g, move}, nil
}

func (g *GameState) IsOver() bool {
	if g.lastMove == nil {
		return false
	}
	switch g.lastMove.(type) {
	case *Resign:
		return true
	}
	if g.previousState.lastMove == nil {
		return false
	}
	switch g.lastMove.(type) {
	case *Pass:
		switch g.previousState.lastMove.(type) {
		case *Pass:
			return true
		}
	}
	return false
}

func (gs *GameState) isMoveSelfCapture(player PlayerColor, move Move) (bool, error) {
	switch m := move.(type) {
	case *Pass:
		return false, nil
	case *Resign:
		return false, nil
	case *Play:
		return gs.Board.isSelfCapture(player, m.Point), nil
	default:
		return false, errors.New("unknown move type")
	}
}

func (gs *GameState) situation() *Situation {
	// if gs.previousState == nil {
	// 	return &Situation{gs.NextPlayer, 0}
	// }
	return &Situation{gs.NextPlayer, gs.Board.hash}
}

type Situation struct {
	nextPlayer PlayerColor
	hash       int
}

func (s *Situation) equals(o *Situation) bool {
	return s.nextPlayer == o.nextPlayer && s.hash == o.hash
}

func (gs *GameState) doesMoveViolateKo(player PlayerColor, move Move) (bool, error) {
	switch m := move.(type) {
	case *Pass:
		return false, nil
	case *Resign:
		return false, nil
	case *Play:
		nextBoard := gs.Board.copy()
		nextBoard.placeStone(player, m.Point)
		nextSituation := &Situation{player.otherPlayer(), nextBoard.hash}
		for pastState := gs.previousState; pastState != nil; pastState = pastState.previousState {
			if pastState.situation().equals(nextSituation) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, errors.New("unknown move type")
	}
}

func (gs *GameState) IsValidMove(move Move) (bool, error) {
	if gs.IsOver() {
		return false, nil
	}
	switch m := move.(type) {
	case *Pass:
		return false, nil
	case *Resign:
		return false, nil
	case *Play:
		boardState, _ := gs.Board.Get(m.Point)
		selfCapture, _ := gs.isMoveSelfCapture(gs.NextPlayer, m)
		violatesKo, _ := gs.doesMoveViolateKo(gs.NextPlayer, m)
		return boardState == Empty && !selfCapture && !violatesKo, nil
	default:
		return false, errors.New("unknown move type")
	}
}
