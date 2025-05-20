package deck_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/deck"
)

// TestDrawSingleCard tests drawing a single card from the deck
func TestDrawSingleCard(t *testing.T) {
	d := deck.New()
	initialCount := d.Count()

	// Draw a card
	_, ok := d.Draw()
	if !ok {
		t.Error("Failed to draw a card from a non-empty deck")
	}

	// Deck should have one less card
	if d.Count() != initialCount-1 {
		t.Errorf("Expected deck to have %d cards after drawing, got %d", initialCount-1, d.Count())
	}

	// Drawing from an empty deck
	emptyDeck := deck.New()
	emptyDeck.Cards = nil
	_, ok = emptyDeck.Draw()
	if ok {
		t.Error("Expected drawing from empty deck to return false")
	}
}

// TestDrawMultipleCards tests drawing multiple cards from the deck
func TestDrawMultipleCards(t *testing.T) {
	d := deck.New()
	initialCount := d.Count()

	// Draw 5 cards
	cards, ok := d.DrawN(5)
	if !ok {
		t.Error("Failed to draw 5 cards from a full deck")
	}

	if len(cards) != 5 {
		t.Errorf("Expected to draw 5 cards, got %d", len(cards))
	}

	// Deck should have 5 fewer cards
	if d.Count() != initialCount-5 {
		t.Errorf("Expected deck to have %d cards after drawing 5, got %d", initialCount-5, d.Count())
	}

	// Try to draw more cards than available
	_, ok = d.DrawN(100)
	if ok {
		t.Error("Expected drawing more cards than available to return false")
	}
}
