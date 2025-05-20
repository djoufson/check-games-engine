package deck_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/deck"
)

// TestShouldMaintainCardCount_WhenAddingCardToBottom tests that adding to bottom doesn't change count
func TestShouldMaintainCardCount_WhenAddingCardToBottom(t *testing.T) {
	// Arrange
	d := deck.New()
	initialCount := d.Count()
	drawnCard, _ := d.Draw()

	// Act
	d.AddToBottom(drawnCard)

	// Assert
	if d.Count() != initialCount {
		t.Errorf("Expected deck to have %d cards after adding to bottom, got %d", initialCount, d.Count())
	}
}

// TestShouldPositionCardAtBottom_WhenAddingCardToBottom tests card is placed at the bottom
func TestShouldPositionCardAtBottom_WhenAddingCardToBottom(t *testing.T) {
	// Arrange
	d := deck.New()
	drawnCard, _ := d.Draw()

	// Act
	d.AddToBottom(drawnCard)

	// Assert
	lastIndex := d.Count() - 1
	if d.Cards[lastIndex].Suit != drawnCard.Suit || d.Cards[lastIndex].Rank != drawnCard.Rank {
		t.Error("Card added to bottom is not the last card in the deck")
	}
}

// TestShouldPositionCardAtTop_WhenAddingCardToTop tests card is placed at the top
func TestShouldPositionCardAtTop_WhenAddingCardToTop(t *testing.T) {
	// Arrange
	d := deck.New()
	drawnCard, _ := d.Draw()

	// Act
	d.AddToTop(drawnCard)

	// Assert
	if d.Cards[0].Suit != drawnCard.Suit || d.Cards[0].Rank != drawnCard.Rank {
		t.Error("Card added to top is not the first card in the deck")
	}
}

// TestShouldIncreaseCardCount_WhenAddingMultipleCardsToBottom tests adding multiple cards increases count correctly
func TestShouldIncreaseCardCount_WhenAddingMultipleCardsToBottom(t *testing.T) {
	// Arrange
	d := deck.New()
	drawnCards, _ := d.DrawN(10)
	initialCount := d.Count()

	// Act
	d.AddManyToBottom(drawnCards)

	// Assert
	if d.Count() != initialCount+10 {
		t.Errorf("Expected deck to have %d cards after adding many to bottom, got %d", initialCount+10, d.Count())
	}
}

// TestShouldReportNotEmpty_WhenDeckHasCards tests IsEmpty returns false for non-empty deck
func TestShouldReportNotEmpty_WhenDeckHasCards(t *testing.T) {
	// Arrange
	d := deck.New()

	// Act & Assert
	if d.IsEmpty() {
		t.Error("New deck should not be empty")
	}
}

// TestShouldReportEmpty_WhenDeckHasNoCards tests IsEmpty returns true for empty deck
func TestShouldReportEmpty_WhenDeckHasNoCards(t *testing.T) {
	// Arrange
	d := deck.New()

	// Act
	d.DrawN(d.Count()) // Draw all cards

	// Assert
	if !d.IsEmpty() {
		t.Error("Deck should be empty after drawing all cards")
	}
}
