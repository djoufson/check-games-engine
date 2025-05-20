package deck_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/deck"
)

// TestDeckCardAddition tests adding cards to the deck (top and bottom)
func TestDeckCardAddition(t *testing.T) {
	d := deck.New()
	initialCount := d.Count()

	// Draw a card to add back later
	drawnCard, _ := d.Draw()

	// Add a card to the bottom
	d.AddToBottom(drawnCard)
	if d.Count() != initialCount {
		t.Errorf("Expected deck to have %d cards after adding to bottom, got %d", initialCount, d.Count())
	}

	// The added card should be the last one
	if d.Cards[d.Count()-1].Suit != drawnCard.Suit || d.Cards[d.Count()-1].Rank != drawnCard.Rank {
		t.Error("Card added to bottom is not the last card in the deck")
	}

	// Draw again to test adding to top
	drawnCard, _ = d.Draw()

	// Add a card to the top
	d.AddToTop(drawnCard)

	// The added card should be the first one
	if d.Cards[0].Suit != drawnCard.Suit || d.Cards[0].Rank != drawnCard.Rank {
		t.Error("Card added to top is not the first card in the deck")
	}
}

// TestAddManyCardsToDeck tests adding multiple cards to the deck at once
func TestAddManyCardsToDeck(t *testing.T) {
	d := deck.New()

	// Draw 10 cards
	drawnCards, _ := d.DrawN(10)
	initialCount := d.Count()

	// Add the cards back to the bottom
	d.AddManyToBottom(drawnCards)

	if d.Count() != initialCount+10 {
		t.Errorf("Expected deck to have %d cards after adding many to bottom, got %d", initialCount+10, d.Count())
	}
}

// TestDeckEmptyStatus tests checking if a deck is empty
func TestDeckEmptyStatus(t *testing.T) {
	d := deck.New()
	if d.IsEmpty() {
		t.Error("New deck should not be empty")
	}

	// Draw all cards
	d.DrawN(d.Count())

	if !d.IsEmpty() {
		t.Error("Deck should be empty after drawing all cards")
	}
}
