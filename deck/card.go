//go:generate stringer -type=Suit,Rank

package deck

import "fmt"

// Suit is one of the categories into which the cards of a deck are divided.
type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker // special card
)

// Rank is the ranking of cards from low to high (1-13)
type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Card is an individual card which has a Suit and a Rank
type Card struct {
	Suit
	Rank
}

var suits = [...]Suit{Spade, Diamond, Club, Heart}

const (
	minRank = Ace
	maxRank = King
)

// String returns a formatted string of the Card Rank and Suit
func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New creates a new deck of cards
func New() []Card {
	var cards []Card

	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{
				Suit: suit,
				Rank: rank,
			})
		}
	}

	return cards
}
