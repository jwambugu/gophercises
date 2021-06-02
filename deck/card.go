//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

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
var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

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
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card

	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{
				Suit: suit,
				Rank: rank,
			})
		}
	}

	for _, opt := range opts {
		cards = opt(cards)
	}

	return cards
}

// absoluteRank ensures the minimum rank a card can get is 13 for other cards except Spade
// which will have a minimum rank of 1
func absoluteRank(c Card) int {
	return int(c.Suit) * int(maxRank+c.Rank)
}

// Less sorts the cards from the min absolute rank to the max absolute rank
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absoluteRank(cards[i]) < absoluteRank(cards[j])
	}
}

// DefaultSort sorts the decks as if they were new
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// Sort takes in a custom sorting function and returns the sorted cards
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// Shuffle shuffles a deck in random order
func Shuffle(cards []Card) []Card {
	shuffledCards := make([]Card, len(cards))

	perm := shuffleRand.Perm(len(cards))

	for i, j := range perm {
		shuffledCards[i] = cards[j]
	}

	return shuffledCards
}

// Jokers adds n number of Joker to the deck
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Suit: Joker,
				Rank: Rank(i),
			})
		}
		return cards
	}
}

// Filter creates a new deck without the cards passed in the filter function
func Filter(filter func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var newCards []Card

		for _, c := range cards {
			if !filter(c) {
				newCards = append(newCards, c)
			}
		}

		return newCards
	}
}

// Deck creates n number of decks
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var deck []Card

		for i := 0; i < n; i++ {
			deck = append(deck, cards...)
		}

		return deck
	}
}
