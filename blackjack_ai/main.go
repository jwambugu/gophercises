package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/blackjack_ai/blackjack"
)

func main() {
	game := blackjack.New(blackjack.Options{})

	winnings := game.Play(blackjack.HumanAI())

	fmt.Println(winnings)
}
