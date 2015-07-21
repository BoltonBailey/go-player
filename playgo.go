package main

import (
	"gogame"
	"math/rand"
	"time"
)

const REPITITIONS int = 1

func main() {

	rand.Seed(time.Now().Unix())

	for i := 0; i < REPITITIONS; i++ {
		gameToShow := gogame.MakeGame(gogame.RandomAutomatonPlayer(), gogame.RandomAutomatonPlayer())
		gameToShow.PlayGame()
		gameToShow.PrintGame()

	}

}
