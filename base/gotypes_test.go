package base

import (
	"testing"
)

func TestOtherPlayer(t *testing.T) {
	expectWhite := Black.otherPlayer()
	if expectWhite != White {
		t.Errorf("Expected White got %v", expectWhite)
	}
	expectBlack := White.otherPlayer()
	if expectBlack != Black {
		t.Errorf("Expected Black got %v", expectBlack)
	}
}

func TestNeighbours(t *testing.T) {
	ns := (&Point{1, 1}).Neighbours()
	if len(ns) != 4 {
		t.Errorf("should have 4 neighbours got %v", len(ns))
	}
}
