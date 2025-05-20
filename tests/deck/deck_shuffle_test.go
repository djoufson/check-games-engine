package deck_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/djoufson/check-games-engine/deck"
)

// TestShouldHaveIdenticalCardOrder_WhenCreatingTwoNewDecks verifies that new decks have identical card orders
func TestShouldHaveIdenticalCardOrder_WhenCreatingTwoNewDecks(t *testing.T) {
	// Arrange
	deck1 := deck.New()
	deck2 := deck.New()

	// Act & Assert
	for i := range deck1.Cards {
		if deck1.Cards[i].Suit != deck2.Cards[i].Suit || deck1.Cards[i].Rank != deck2.Cards[i].Rank {
			t.Errorf("Decks are not identical before shuffling at index %d", i)
		}
	}
}

// TestShouldChangeCardOrder_WhenShufflingDeck verifies that shuffling changes card order
func TestShouldChangeCardOrder_WhenShufflingDeck(t *testing.T) {
	// Arrange
	deck1 := deck.New()
	deck2 := deck.New()

	// Act
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	deck1.Shuffle(r)

	// Assert - Decks should now be different (this could theoretically fail, but probability is extremely low)
	identical := true
	for i := range deck1.Cards {
		if deck1.Cards[i].Suit != deck2.Cards[i].Suit || deck1.Cards[i].Rank != deck2.Cards[i].Rank {
			identical = false
			break
		}
	}

	if identical {
		t.Error("Decks are still identical after shuffling")
	}
}
