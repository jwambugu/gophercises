package main

import (
	"fmt"
	"github.com/jwambugu/gophercises/deck"
	"strings"
)

type Hand []deck.Card

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

func main() {
	cards := deck.New(deck.Deck(3), deck.Shuffle)

	var card deck.Card
	var player, dealer Hand

	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&player, &dealer} {
			card, cards = draw(cards)
			*hand = append(*hand, card)
		}
	}

	var input string
	for input != "s" {
		fmt.Println("Player:", player)
		fmt.Println("Dealer:", dealer.DealerString())

		fmt.Println("What will you do? (h)it, (s)tand")
		_, _ = fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			card, cards = draw(cards)
			player = append(player, card)
		}
	}

	playerScore, dealerScore := player.Score(), dealer.Score()

	// If dealer score <= 16, hit
	// If dealer has a sort 17, hit
	// Soft 17 is when an ace as 11 and the score is 17
	for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.MinScore() != 17) {
		card, cards = draw(cards)
		dealer = append(dealer, card)
	}

	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", player, "\nScore: ", playerScore)
	fmt.Println("Dealer:", dealer, "\nScore: ", dealerScore)

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
}
