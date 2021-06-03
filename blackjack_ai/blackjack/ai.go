package blackjack

import (
	"fmt"
	"github.com/jwambugu/gophercises/deck"
)

type Move func(GameState) GameState

type AI interface {
	Results(hand, dealer []deck.Card)
	Play(hand [][]deck.Card, dealer deck.Card) Move
	Bet() int
}

type HumanAI struct {
}

type GameState struct {
}

func Hit(gs GameState) GameState {
	return gs
}

func Stand(gs GameState) GameState {
	return gs
}

func (ai *HumanAI) Bet() int {
	return 1
}

func (ai *HumanAI) Results(hand, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
}

func (ai *HumanAI) Play(hand [][]deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)

		fmt.Println("What will you do? (h)it, (s)tand")

		var input string
		_, _ = fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			return Hit
		case "s":
			return Stand
		default:
			fmt.Println("Invalid option: ", input)
		}
	}
}
