package deck

import (
	"fmt"
	"math/rand"
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

func TestShuffle(t *testing.T) {
	// [40 35 ...]
	shuffleRand = rand.New(rand.NewSource(0))

	originalDeck := New()
	firstCard := originalDeck[40]
	secondCard := originalDeck[35]

	shuffledCards := New(Shuffle)

	if shuffledCards[0] != firstCard {
		t.Errorf("expected the first card to be %q , got %q", firstCard, shuffledCards[0])
	}

	if shuffledCards[1] != secondCard {
		t.Errorf("expected the first card to be %q , got %q", firstCard, shuffledCards[1])

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

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}

	cards := New(Filter(filter))

	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			t.Error("expected all twos and threes to be filtered out")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))

	// 13 ranks * 4 suits * 3 decks
	expected := 13 * 4 * 3

	if len(cards) != expected {
		t.Errorf("expected %d cards, got %d", expected, len(cards))
	}
}
