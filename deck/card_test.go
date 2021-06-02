package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Nine, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Two of Spades
	// Nine of Diamonds
	// Jack of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()

	// 13 ranks * 4 suits
	if len(cards) != 52 {
		t.Errorf("expected %d cards in a deck, got %d", 52, len(cards))
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)

	firstCard := Card{
		Suit: Spade,
		Rank: Ace,
	}

	if cards[0] != firstCard {
		t.Errorf("expected '%v' as first card, got '%v'", firstCard, cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))

	firstCard := Card{
		Suit: Spade,
		Rank: Ace,
	}

	if cards[0] != firstCard {
		t.Errorf("expected '%v' as first card, got '%v'", firstCard, cards[0])
	}
}

func TestJokers(t *testing.T) {
	jokers := 3

	cards := New(Jokers(jokers))
	count := 0

	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}

	if count != 3 {
		t.Errorf("expected %d Jokers, got %d", jokers, count)
	}
}
