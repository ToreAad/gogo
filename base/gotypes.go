package base

type PlayerColor int

const (
	Empty PlayerColor = iota
	Black
	White
)

func (player PlayerColor) otherPlayer() PlayerColor {
	if player == Black {
		return White
	}
	return Black
}

type Point struct {
	Row int
	Col int
}

func (p *Point) Neighbours() []Point {
	return []Point{
		{p.Row - 1, p.Col},
		{p.Row + 1, p.Col},
		{p.Row, p.Col - 1},
		{p.Row, p.Col + 1}}
}
