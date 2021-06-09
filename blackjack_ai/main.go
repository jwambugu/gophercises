package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/blackjack_ai/blackjack"
)

func main() {
	game := blackjack.New(blackjack.Options{
		Hands:           2,
		Decks:           3,
		BlackjackPayout: 1.5,
	})

	winnings := game.Play(blackjack.HumanAI())

	fmt.Println(winnings)
}
