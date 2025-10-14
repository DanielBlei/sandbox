// Knights Game - Challenge Solution focused on game core logic.
// Intentionally minimal and non-modular to reflect a timed coding exercise.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type knightID int

const (
	k1 knightID = iota
	k2
	k3
	k4
	k5
	k6
)

var allKnights = []knightID{k1, k2, k3, k4, k5, k6}

type knight struct {
	idx    knightID
	health int
	alive  bool
}

type game struct {
	gameOver          bool
	currentKnightTurn knightID
	knights           []knight
	seed              *rand.Rand
}

// Get the next knight to either hit or play
func (g *game) getNextKnight() knightID {
	nextKnightIndex := (int(g.currentKnightTurn) + 1) % len(g.knights)
	for !g.knights[nextKnightIndex].alive {
		nextKnightIndex = (nextKnightIndex + 1) % len(g.knights)
	}
	return knightID(nextKnightIndex)
}

func (g *game) knightHitValue() int {
	return g.seed.Intn(6) + 1
}

func main() {
	var seed int64
	flag.Int64Var(&seed, "seed", time.Now().Unix(), "Random seed for reproducible games (default: current time)")
	flag.Parse()

	random := rand.New(rand.NewSource(seed))
	gameState := game{
		knights:           make([]knight, len(allKnights)),
		currentKnightTurn: k1,
		seed:              random,
	}

	for i, k := range allKnights {
		gameState.knights[i] = knight{idx: k, health: 10, alive: true}
	}
	// fmt.Println(gameState)

	countAlive := len(gameState.knights)
	for !gameState.gameOver {
		// Check for winner first
		if countAlive <= 1 {
			fmt.Printf("K%d wins\n", gameState.currentKnightTurn+1)
			gameState.gameOver = true
			break
		}

		// Get the next knight to hit
		nextKnight := gameState.getNextKnight()

		currentKnightHitsWith := gameState.knightHitValue()
		gameState.knights[nextKnight].health = gameState.knights[nextKnight].health - currentKnightHitsWith
		if gameState.knights[nextKnight].health < 1 {
			gameState.knights[nextKnight].alive = false
			countAlive--
		}
		fmt.Printf("K%d hits K%d for %d\n", gameState.currentKnightTurn+1, gameState.knights[nextKnight].idx+1, currentKnightHitsWith)

		// Move to next turn (the knight who just got hit, or next alive knight if they died)
		if gameState.knights[nextKnight].alive {
			gameState.currentKnightTurn = nextKnight
		} else {
			gameState.currentKnightTurn = gameState.getNextKnight()
		}
	}

	// if gameState.gameOver {
	// 	fmt.Print("Game Over")
	// 	// Could add more post-game logic here e.g ask to play again, save, statistics, etc.
	// }

}
