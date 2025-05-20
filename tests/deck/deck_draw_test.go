package deck_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/deck"
)

// TestShouldRemoveOneCard_WhenDrawingSingleCard tests that drawing reduces deck size by one
func TestShouldRemoveOneCard_WhenDrawingSingleCard(t *testing.T) {
	// Arrange
	d := deck.New()
	initialCount := d.Count()

	// Act
	_, ok := d.Draw()

	// Assert
	if !ok {
		t.Error("Failed to draw a card from a non-empty deck")
	}

	if d.Count() != initialCount-1 {
		t.Errorf("Expected deck to have %d cards after drawing, got %d", initialCount-1, d.Count())
	}
}

// TestShouldReturnFailure_WhenDrawingFromEmptyDeck tests that drawing from empty deck fails
func TestShouldReturnFailure_WhenDrawingFromEmptyDeck(t *testing.T) {
	// Arrange
	emptyDeck := deck.New()
	emptyDeck.Cards = nil

	// Act
	_, ok := emptyDeck.Draw()

	// Assert
	if ok {
		t.Error("Expected drawing from empty deck to return false")
	}
}

// TestShouldReturnRequestedNumberOfCards_WhenDrawingMultipleCards tests drawing multiple cards returns correct count
func TestShouldReturnRequestedNumberOfCards_WhenDrawingMultipleCards(t *testing.T) {
	// Arrange
	d := deck.New()

	// Act
	cards, ok := d.DrawN(5)

	// Assert
	if !ok {
		t.Error("Failed to draw 5 cards from a full deck")
	}

	if len(cards) != 5 {
		t.Errorf("Expected to draw 5 cards, got %d", len(cards))
	}
}

// TestShouldReduceDeckSize_WhenDrawingMultipleCards tests that drawing multiple cards reduces deck size accordingly
func TestShouldReduceDeckSize_WhenDrawingMultipleCards(t *testing.T) {
	// Arrange
	d := deck.New()
	initialCount := d.Count()
	drawCount := 5

	// Act
	_, ok := d.DrawN(drawCount)

	// Assert
	if !ok {
		t.Error("Failed to draw cards from a full deck")
	}

	if d.Count() != initialCount-drawCount {
		t.Errorf("Expected deck to have %d cards after drawing %d, got %d",
			initialCount-drawCount, drawCount, d.Count())
	}
}

// TestShouldReturnFailure_WhenDrawingMoreCardsThanAvailable tests over-drawing from deck
func TestShouldReturnFailure_WhenDrawingMoreCardsThanAvailable(t *testing.T) {
	// Arrange
	d := deck.New()

	// Act
	_, ok := d.DrawN(100) // Try to draw more cards than available

	// Assert
	if ok {
		t.Error("Expected drawing more cards than available to return false")
	}
}
