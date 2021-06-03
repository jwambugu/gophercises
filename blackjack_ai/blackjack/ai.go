package blackjack

import (
	"fmt"
	"github.com/jwambugu/gophercises/deck"
)

type AI interface {
	Results(hand [][]deck.Card, dealer []deck.Card)
	Play(hand []deck.Card, dealer deck.Card) Move
	Bet() int
}

type humanAI struct {
}

func HumanAI() AI {
	return humanAI{}
}

type dealerAI struct {
}

func (ai *dealerAI) Bet() int {
	return 1
}

func (ai humanAI) Bet() int {
	return 1
}

func (ai *dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	// DO NOTHING
}

func (ai humanAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
	fmt.Println()
}

func (ai *dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	// If dealer score <= 16, hit || If dealer has a sort 17, hit.
	// Soft 17 is when an ace as 11 and the score is 17
	dealerScore := Score(hand...)

	if dealerScore <= 16 || (dealerScore == 17 && Soft(hand...)) {
		return MoveHit
	}

	return MoveStand
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)

		fmt.Println("What will you do? (h)it, (s)tand")

		var input string
		_, _ = fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("Invalid option: ", input)
		}
	}
}
