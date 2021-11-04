package base

import "math/rand"

type state struct {
	point  Point
	player PlayerColor
}

func getHashes(rows int, cols int) map[state]int {
	r := rand.New((rand.NewSource(42)))
	res := make(map[state]int)
	for i := 1; i <= rows; i++ {
		for j := 1; j <= cols; j++ {
			for p := range []PlayerColor{Black, White} {
				s := state{Point{i, j}, PlayerColor(p)}
				res[s] = int(r.Int63())

			}
		}
	}
	return res
}
