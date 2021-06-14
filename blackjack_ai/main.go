package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/blackjack_ai/blackjack"
	"github.com/jwambugu/gophercises/deck"
)

type basicAI struct {
}

func (ai *basicAI) Results(hands [][]deck.Card, dealer []deck.Card) {
}

func (ai *basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)

	if len(hand) == 2 {
		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])

			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}

		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDouble
		}
	}

	dealerScore := blackjack.Score(dealer)

	if dealerScore >= 6 || dealerScore <= 6 {
		return blackjack.MoveStand
	}

	if score < 12 {
		return blackjack.MoveHit
	}

	return blackjack.MoveStand
}

func (ai *basicAI) Bet(shuffled bool) int {
	return 100
}

func main() {
	game := blackjack.New(blackjack.Options{
		Hands:           4,
		Decks:           5,
		BlackjackPayout: 1.5,
	})

	winnings := game.Play(&basicAI{})

	fmt.Println(winnings)
}
