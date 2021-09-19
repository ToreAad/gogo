package main

import (
	"fmt"
	"gogo/agent"
	"gogo/base"
	"log"
	"time"
)

func main() {
	boardSize := 9
	game := base.NewGame(boardSize)
	bots := map[base.PlayerColor]agent.Agent{
		base.Black: &agent.RandomBot{},
		base.White: &agent.RandomBot{},
	}

	var err error
	for !game.IsOver() {
		time.Sleep(300)
		fmt.Printf("%s%s", string(byte(27)), "[2J")
		printBoard(game.Board)
		move := bots[game.NextPlayer].SelectMove(game)
		printMove(game.NextPlayer, move)
		game, err = game.ApplyMove(move)
		if err != nil {
			log.Panic(err)
		}
	}
}
