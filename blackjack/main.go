package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/deck"
	"strings"
)

type Hand []deck.Card

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

func (h Hand) String() string {
	strs := make([]string, len(h))

	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func (h Hand) MinScore() int {
	score := 0

	for _, c := range h {
		score += min(int(c.Rank), 10)
	}

	return score
}

func (h Hand) Score() int {
	minScore := h.MinScore()

	if minScore > 11 {
		return minScore
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			// Ace is currently worth 1 point, we are making it 11 points
			return minScore + 10
		}
	}

	return minScore
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func clone(gs GameState) GameState {
	s := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}

	copy(s.Deck, gs.Deck)
	copy(s.Player, gs.Player)
	copy(s.Dealer, gs.Dealer)

	return s
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func Shuffle(gs GameState) GameState {
	s := clone(gs)
	s.Deck = deck.New(deck.Deck(3), deck.Shuffle)

	return s
}

func Deal(gs GameState) GameState {
	s := clone(gs)
	s.Player = make(Hand, 0, 5)
	s.Dealer = make(Hand, 0, 5)

	var card deck.Card

	for i := 0; i < 2; i++ {
		card, s.Deck = draw(s.Deck)
		s.Player = append(s.Player, card)

		card, s.Deck = draw(s.Deck)
		s.Dealer = append(s.Dealer, card)
	}

	s.State = StatePlayerTurn
	return s
}

func Hit(gs GameState) GameState {
	s := clone(gs)
	hand := s.CurrentPlayer()

	var card deck.Card
	card, s.Deck = draw(s.Deck)

	*hand = append(*hand, card)
	return s
}

func Stand(gs GameState) GameState {
	s := clone(gs)
	s.State++
	return s
}

func EndHand(gs GameState) GameState {
	s := clone(gs)

	playerScore, dealerScore := gs.Player.Score(), gs.Dealer.Score()

	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", gs.Player, "\nScore: ", playerScore)
	fmt.Println("Dealer:", gs.Player, "\nScore: ", dealerScore)

	switch {
	case playerScore > 21:
		fmt.Println("You busted!")
	case dealerScore > 21:
		fmt.Println("Dealer busted!")
	case playerScore > dealerScore:
		fmt.Println("You win!")
	case dealerScore > playerScore:
		fmt.Println("You lose!")
	case dealerScore == playerScore:
		fmt.Println("You draw!")
	}

	fmt.Println()

	s.Player = nil
	s.Dealer = nil
	return s
}

func main() {
	var gs GameState
	gs = Shuffle(gs)

	for i := 0; i < 10; i++ {
		gs = Deal(gs)
		var input string
		for gs.State == StatePlayerTurn {
			fmt.Println("Player:", gs.Player)
			fmt.Println("Dealer:", gs.Dealer.DealerString())

			fmt.Println("What will you do? (h)it, (s)tand")
			_, _ = fmt.Scanf("%s\n", &input)

			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid option: ", input)
			}
		}

		for gs.State == StateDealerTurn {
			// If dealer score <= 16, hit || If dealer has a sort 17, hit. Soft 17 is when an ace as 11 and the score is 17
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}
}
