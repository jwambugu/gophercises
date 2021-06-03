package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/blackjack_ai/blackjack"
)

func main() {
	game := blackjack.New()

	winnings := game.Play(blackjack.HumanAI())

	fmt.Println(winnings)
}
